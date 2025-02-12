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
	google.golang.org/grpc v1.70.0
)

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.26.1 // indirect
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250207221924-e9438ea467c6 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250204164813-702378808489 // indirect
	google.golang.org/protobuf v1.36.5 // indirect
)
