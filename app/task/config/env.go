package config

import (
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// EnvConfig 对 .env 文件的监听
type EnvConfig struct {
	MySQLUser     string
	MySQLPassword string
	MySQLHost     string
	MySQLPort     string
	MySQLDb       string
}

var envConfig *EnvConfig

// LoadEnvConfig 加载环境变量到结构体
func LoadEnvConfig() error {
	envConfig = &EnvConfig{
		MySQLUser:     os.Getenv("MYSQL_USER"),
		MySQLPassword: os.Getenv("MYSQL_PASSWORD"),
		MySQLHost:     os.Getenv("MYSQL_HOST"),
		MySQLPort:     os.Getenv("MYSQL_PORT"),
		MySQLDb:       os.Getenv("MYSQL_DB"),
	}
	return nil
}

// GetEnvConfig 获取环境变量配置
func GetEnvConfig() *EnvConfig {
	if envConfig == nil {
		if err := LoadEnvConfig(); err != nil {
			zap.L().Error("Failed to load env config", zap.Error(err))
			return nil
		}
	}
	return envConfig
}

func WatchEnvFile(envPath string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		zap.L().Error("Failed to create env watcher", zap.Error(err))
		return
	}

	if err := watcher.Add(envPath); err != nil {
		zap.L().Error("Failed to watch .env file", zap.Error(err))
		return
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					zap.L().Info(".env file changed", zap.String("file", event.Name))

					if err := godotenv.Load(); err != nil {
						zap.L().Error("Failed to reload .env", zap.Error(err))
						continue
					}

					if err := LoadEnvConfig(); err != nil {
						zap.L().Error("Failed to reload env config", zap.Error(err))
						continue
					}

					// 通知观察者
					if conf != nil {
						for _, observer := range conf.observers {
							observer.OnConfigChange(conf, DBChange)
						}
					}
				}
			case err := <-watcher.Errors:
				zap.L().Error("Env watcher error", zap.Error(err))
			}
		}
	}()
}
