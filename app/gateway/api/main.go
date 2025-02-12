package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/carzyfrankie/app/gateway/mws"
	"github.com/crazyfrankie/todolist/app/task/rpc_gen/task"
	"github.com/crazyfrankie/todolist/app/user/rpc_gen/user"
)

func main() {
	mux := runtime.NewServeMux(runtime.WithMetadata(func(ctx context.Context, request *http.Request) metadata.MD {
		md := metadata.MD{}

		if userID, ok := request.Context().Value("user_id").(string); ok {
			md.Set("user_id", userID)
		}

		return md
	}))

	u := initUserClient()
	t := initTaskClient()

	err := user.RegisterUserServiceHandlerClient(context.Background(), mux, u)
	if err != nil {
		panic(err)
	}
	err = task.RegisterTaskServiceHandlerClient(context.Background(), mux, t)
	if err != nil {
		panic(err)
	}
	handler := mws.NewAuthBuilder().
		IgnorePath("/api/user/login").
		IgnorePath("/api/user/register").
		Auth(mux)

	log.Printf("HTTP server listening on %s", "localhost:9091")
	if err := http.ListenAndServe(fmt.Sprintf("%s", "localhost:9091"), handler); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func initUserClient() user.UserServiceClient {
	conn, err := grpc.NewClient("localhost:8081",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	return user.NewUserServiceClient(conn)
}

func initTaskClient() task.TaskServiceClient {
	conn, err := grpc.NewClient("localhost:8082", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	return task.NewTaskServiceClient(conn)
}
