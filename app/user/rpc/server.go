package rpc

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/crazyfrankie/framework-plugin/grpcx/interceptor/circuitbreaker"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/crazyfrankie/todolist/app/user/config"
	"github.com/crazyfrankie/todolist/app/user/pkg/registry"
	"github.com/crazyfrankie/todolist/app/user/rpc/server"
)

type Server struct {
	*grpc.Server
	Addr     string
	cli      *clientv3.Client
	register *registry.ServiceRegistry
	mu       sync.RWMutex
	listener net.Listener
	svc      *server.UserServer
}

func NewServer(t *server.UserServer, client *clientv3.Client) *Server {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	logTraceID := func(ctx context.Context) logging.Fields {
		if span := oteltrace.SpanContextFromContext(ctx); span.IsSampled() {
			return logging.Fields{"traceID", span.SpanID().String()}
		}
		return nil
	}

	tp := initTracerProvider("todolist/server/user")
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	s := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(initInterceptor(logger), logging.WithFieldsFromContext(logTraceID)),
			circuitbreaker.NewInterceptorBuilder().Build()),
	)
	t.RegisterServer(s)

	rpcServer := &Server{
		Server: s,
		Addr:   config.GetConf().Server.Addr,
		cli:    client,
		svc:    t,
	}
	register, err := registry.NewServiceRegistry(rpcServer.cli, rpcServer.Addr)
	if err != nil {
		panic(err)
	}
	config.GetConf().AddObserver(rpcServer)
	rpcServer.register = register

	return rpcServer
}

func (s *Server) OnConfigChange(c *config.Config, changeType config.ConfigChangeType) {
	if changeType == config.ServerChange && c.Server.Addr != s.Addr {
		if err := s.register.Unregister(); err != nil {
			zap.L().Error("Failed to unregister service", zap.Error(err))
		}

		s.mu.Lock()
		oldListener := s.listener
		oldServer := s.Server
		s.Addr = c.Server.Addr
		s.mu.Unlock()

		if oldServer != nil {
			oldServer.GracefulStop()
		}

		listener, err := net.Listen("tcp", s.Addr)
		if err != nil {
			zap.L().Error("Failed to create listener", zap.Error(err))
			return
		}

		register, err := registry.NewServiceRegistry(s.cli, s.Addr)
		if err != nil {
			zap.L().Error("Failed to create registry", zap.Error(err))
			listener.Close()
			return
		}

		if err := register.Register(); err != nil {
			zap.L().Error("Failed to register", zap.Error(err))
			listener.Close()
			return
		}

		s.mu.Lock()
		s.listener = listener
		s.register = register
		s.mu.Unlock()

		if oldListener != nil {
			oldListener.Close()
		}

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			if err := s.Server.Serve(listener); err != nil {
				zap.L().Error("Server serve error", zap.Error(err))
			}
		}()

		<-quit

		if err := s.ShutDown(); err != nil {
			zap.L().Error("Server shutdown error", zap.Error(err))
		}

		zap.L().Info("Server exited gracefully")
	}
}

func (s *Server) Serve() error {
	conn, err := net.Listen("tcp", s.Addr)
	if err != nil {
		panic(err)
	}

	s.listener = conn
	err = s.register.Register()
	if err != nil {
		return err
	}

	return s.Server.Serve(conn)
}

func (s *Server) ShutDown() error {
	err := s.cli.Close()
	if err != nil {
		return err
	}

	s.Server.GracefulStop()

	return nil
}

func initInterceptor(l *zap.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, level logging.Level, msg string, fields ...any) {
		f := make([]zap.Field, 0, len(fields)/2)

		for i := 0; i < len(fields); i += 2 {
			key := fields[i]
			value := fields[i+1]

			switch v := value.(type) {
			case string:
				f = append(f, zap.String(key.(string), v))
			case bool:
				f = append(f, zap.Bool(key.(string), v))
			case int:
				f = append(f, zap.Int(key.(string), v))
			default:
				f = append(f, zap.Any(key.(string), v))
			}
		}

		logger := l.WithOptions(zap.AddCallerSkip(1)).With(f...)

		switch level {
		case logging.LevelInfo:
			logger.Info(msg)
		case logging.LevelDebug:
			logger.Debug(msg)
		case logging.LevelError:
			logger.Error(msg)
		case logging.LevelWarn:
			logger.Warn(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", level))
		}
	})
}

func initTracerProvider(servicename string) *trace.TracerProvider {
	res, err := newResource(servicename, "v0.0.1")
	if err != nil {
		fmt.Printf("failed create resource, %s", err)
	}

	tp, err := newTraceProvider(res)
	if err != nil {
		panic(err)
	}

	return tp
}

func newResource(servicename, serviceVersion string) (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceNameKey.String(servicename),
			semconv.ServiceVersionKey.String(serviceVersion)))
}

func newTraceProvider(res *resource.Resource) (*trace.TracerProvider, error) {
	exporter, err := zipkin.New("http://localhost:9411/api/v2/spans")
	if err != nil {
		return nil, err
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter, trace.WithBatchTimeout(time.Second)), trace.WithResource(res))

	return traceProvider, nil
}
