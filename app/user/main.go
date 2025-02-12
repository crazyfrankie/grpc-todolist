package main

import (
	"github.com/joho/godotenv"

	"github.com/crazyfrankie/todolist/app/user/ioc"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	server := ioc.InitUser()

	err = server.Serve()
	if err != nil {
		panic(err)
	}
}
