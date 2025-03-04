package registry

import (
	"context"
	"fmt"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"go.uber.org/zap"
)

type ServiceRegistry struct {
	Cli        *clientv3.Client
	em         endpoints.Manager
	addr       string
	serviceKey string
	leaseID    clientv3.LeaseID
	stopChan   chan struct{}
	mu         sync.RWMutex
}

func NewServiceRegistry(cli *clientv3.Client, addr string) (*ServiceRegistry, error) {
	em, err := endpoints.NewManager(cli, "service/task")
	if err != nil {
		return nil, err
	}

	return &ServiceRegistry{
		Cli:        cli,
		em:         em,
		addr:       "localhost" + addr,
		serviceKey: "service/task/localhost" + addr,
		stopChan:   make(chan struct{}),
	}, nil
}

func (r *ServiceRegistry) Register() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	leaseResp, err := r.Cli.Grant(ctx, 15)
	if err != nil {
		return fmt.Errorf("failed to grant lease: %v", err)
	}

	r.mu.Lock()
	r.leaseID = leaseResp.ID
	r.mu.Unlock()

	if err = r.em.AddEndpoint(ctx, r.serviceKey, endpoints.Endpoint{Addr: r.addr},
		clientv3.WithLease(leaseResp.ID)); err != nil {
		return err
	}

	// 开始续约
	go r.keepAlive()

	return nil
}

func (r *ServiceRegistry) keepAlive() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, err := r.Cli.KeepAlive(ctx, r.leaseID)
	if err != nil {
		zap.L().Error("KeepAlive failed", zap.Error(err))
	}

	for {
		select {
		case <-r.stopChan:
			return
		case _, ok := <-ch:
			if !ok {
				zap.L().Info("KeepAlive channel closed")
				return
			}
			fmt.Println("Lease renewed")
		case <-ctx.Done():
			return
		}
	}
}

func (r *ServiceRegistry) Unregister() error {
	close(r.stopChan)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := r.em.DeleteEndpoint(ctx, r.serviceKey); err != nil {
		return fmt.Errorf("failed to delete endpoint: %v", err)
	}

	r.mu.RLock()
	leaseID := r.leaseID
	r.mu.RUnlock()

	if _, err := r.Cli.Revoke(ctx, leaseID); err != nil {
		return fmt.Errorf("failed to revoke lease: %v", err)
	}

	return nil
}
