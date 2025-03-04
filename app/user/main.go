package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/crazyfrankie/todolist/app/user/config"
	"github.com/crazyfrankie/todolist/app/user/ioc"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	config.WatchEnvFile(".env")

	server := ioc.InitUser()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.Serve(); err != nil {
			zap.L().Error("Server serve error", zap.Error(err))
		}
	}()

	<-quit

	if err := server.ShutDown(); err != nil {
		zap.L().Error("Server shutdown error", zap.Error(err))
	}

	zap.L().Info("Server exited gracefully")
}
