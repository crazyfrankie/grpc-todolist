package mws

import (
	"net"
	"net/http"
	"strconv"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const (
	tracerKey = "otel-go-contrib-tracer"
	// ScopeName is the instrumentation scope name.
	ScopeName = "go.opentelemetry.io/contrib/instrumentation/github.com/crazyfrankie"
)

type config struct {
	TracerProvider oteltrace.TracerProvider
	Propagators    propagation.TextMapPropagator
}

func Trace(service string, next http.HandlerFunc) http.HandlerFunc {
	cfg := config{}
	if cfg.TracerProvider == nil {
		cfg.TracerProvider = otel.GetTracerProvider()
	}
	tracer := cfg.TracerProvider.Tracer(
		ScopeName,
		oteltrace.WithInstrumentationVersion(Version()),
	)
	if cfg.Propagators == nil {
		cfg.Propagators = otel.GetTextMapPropagator()
	}
	return func(w http.ResponseWriter, req *http.Request) {
		savedCtx := req.Context()
		defer func() {
			req = req.WithContext(savedCtx)
		}()
		ctx := cfg.Propagators.Extract(savedCtx, propagation.HeaderCarrier(req.Header))
		opts := []oteltrace.SpanStartOption{
			oteltrace.WithAttributes(HTTPServerRequest(service, req)...),
			oteltrace.WithAttributes(semconv.HTTPRoute(req.URL.Path)),
			oteltrace.WithSpanKind(oteltrace.SpanKindServer),
		}
		spanName := req.URL.Path
		ctx, span := tracer.Start(ctx, spanName, opts...)
		defer span.End()

		req = req.WithContext(ctx)

		next(w, req)
	}
}

func Version() string {
	return "1.0.0"
}

func HTTPServerRequest(server string, req *http.Request) []attribute.KeyValue {
	n := 4 // Method, scheme, proto, and host name.
	var host string
	var p int
	if server == "" {
		host, p = splitHostPort(req.Host)
	} else {
		// Prioritize the primary server name.
		host, p = splitHostPort(server)
		if p < 0 {
			_, p = splitHostPort(req.Host)
		}
	}
	hostPort := requiredHTTPPort(req.TLS != nil, p)
	if hostPort > 0 {
		n++
	}
	peer, peerPort := splitHostPort(req.RemoteAddr)
	if peer != "" {
		n++
		if peerPort > 0 {
			n++
		}
	}
	useragent := req.UserAgent()
	if useragent != "" {
		n++
	}

	clientIP := serverClientIP(req.Header.Get("X-Forwarded-For"))
	if clientIP != "" {
		n++
	}

	var target string
	if req.URL != nil {
		target = req.URL.Path
		if target != "" {
			n++
		}
	}
	protoName, protoVersion := netProtocol(req.Proto)
	if protoName != "" && protoName != "http" {
		n++
	}
	if protoVersion != "" {
		n++
	}

	attrs := make([]attribute.KeyValue, 0, n)

	attrs = append(attrs, attribute.String("http.route", req.Method))
	attrs = append(attrs, attribute.Bool("https", req.TLS != nil))
	attrs = append(attrs, attribute.String("net.host.name", host))

	if hostPort > 0 {
		attrs = append(attrs, attribute.Int("net.host.port", hostPort))
	}

	if peer != "" {
		// The Go HTTP server sets RemoteAddr to "IP:port", this will not be a
		// file-path that would be interpreted with a sock family.
		attrs = append(attrs, attribute.String("net.sock.peer.addr", peer))
		if peerPort > 0 {
			attrs = append(attrs, attribute.Int("net.sock.peer.port", peerPort))
		}
	}

	if useragent != "" {
		attrs = append(attrs, attribute.String("user_agent.original", useragent))
	}

	if clientIP != "" {
		attrs = append(attrs, attribute.String("http.client_ip", clientIP))
	}

	if target != "" {
		attrs = append(attrs, attribute.String("http.target", target))
	}

	if protoName != "" && protoName != "http" {
		attrs = append(attrs, attribute.String("net.protocol.name", protoName))
	}
	if protoVersion != "" {
		attrs = append(attrs, attribute.String("net.protocol.version", protoVersion))
	}

	return attrs
}

func splitHostPort(hostport string) (host string, port int) {
	port = -1

	if strings.HasPrefix(hostport, "[") {
		addrEnd := strings.LastIndex(hostport, "]")
		if addrEnd < 0 {
			// Invalid hostport.
			return
		}
		if i := strings.LastIndex(hostport[addrEnd:], ":"); i < 0 {
			host = hostport[1:addrEnd]
			return
		}
	} else {
		if i := strings.LastIndex(hostport, ":"); i < 0 {
			host = hostport
			return
		}
	}

	host, pStr, err := net.SplitHostPort(hostport)
	if err != nil {
		return
	}

	p, err := strconv.ParseUint(pStr, 10, 16)
	if err != nil {
		return
	}
	return host, int(p) // nolint: gosec  // Bitsize checked to be 16 above.
}

func netProtocol(proto string) (name string, version string) {
	name, version, _ = strings.Cut(proto, "/")
	name = strings.ToLower(name)
	return name, version
}

func requiredHTTPPort(https bool, port int) int { // nolint:revive
	if https {
		if port > 0 && port != 443 {
			return port
		}
	} else {
		if port > 0 && port != 80 {
			return port
		}
	}
	return -1
}

func serverClientIP(xForwardedFor string) string {
	if idx := strings.Index(xForwardedFor, ","); idx >= 0 {
		xForwardedFor = xForwardedFor[:idx]
	}
	return xForwardedFor
}
