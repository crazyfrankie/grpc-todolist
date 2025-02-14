.PHONY: gen-task
gen-task:
	@protoc --go_out=./app/task/rpc_gen --go-grpc_out=./app/task/rpc_gen --grpc-gateway_out=./app/task/rpc_gen .\idl\todolist\task.proto

.PHONY: gen-user
gen-user:
	@protoc --go_out=./app/user/rpc_gen --go-grpc_out=./app/user/rpc_gen --grpc-gateway_out=./app/user/rpc_gen .\idl\todolist\user.proto

.PHONY: gen-user-main
gen-user-main:
	@cd ./app/user && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

.PHONY: gen-task-main
gen-task-main:
	@cd ./app/task && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

.PHONY: gen-gateway-main
gen-gateway-main:
	@cd ./app/gateway/api && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build