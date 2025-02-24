module github.com/carzyfrankie/app/gateway

go 1.23.6

replace (
	github.com/crazyfrankie/todolist/app/task => ../task
	github.com/crazyfrankie/todolist/app/user => ../user

)

require (
	github.com/crazyfrankie/todolist/app/task v0.0.0-00010101000000-000000000000
	github.com/crazyfrankie/todolist/app/user v0.0.0-00010101000000-000000000000
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.3.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.16.0
	go.etcd.io/etcd/client/v3 v3.5.18
	go.uber.org/zap v1.21.0
	google.golang.org/grpc v1.70.0
)

require (
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	go.etcd.io/etcd/api/v3 v3.5.18 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.18 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20241202173237-19429a94021a // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250204164813-702378808489 // indirect
	google.golang.org/protobuf v1.36.5 // indirect
)
