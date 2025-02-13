package main

import (
	"github.com/crazyfrankie/todolist/app/task/ioc"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	server := ioc.InitTask()

	err = server.Serve()
	if err != nil {
		panic(err)
	}

	err = server.ShutDown()
	if err != nil {
		panic(err)
	}
}
