package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	once sync.Once
	conf *Config
)

type ConfigChangeType int

const (
	ServerChange ConfigChangeType = iota
	DBChange
)

type Observer interface {
	OnConfigChange(*Config, ConfigChangeType)
}

// AddObserver 添加观察者
func (c *Config) AddObserver(o Observer) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.observers = append(c.observers, o)
}

type Config struct {
	Env       string
	Server    Server `yaml:"server"`
	MySQL     MySQL  `yaml:"mysql"`
	JWT       JWT    `yaml:"jwt"`
	ETCD      ETCD   `yaml:"etcd"`
	observers []Observer
	mu        sync.RWMutex
}

type Server struct {
	Addr string `yaml:"addr"`
}

type MySQL struct {
	DSN string `yaml:"dsn"`
}

type ETCD struct {
	Addr string `yaml:"addr"`
}

type JWT struct {
	SecretKey string `yaml:"secretKey"`
}

func GetConf() *Config {
	once.Do(initConf)
	return conf
}

func initConf() {
	prefix := "config"
	path := filepath.Join(prefix, filepath.Join(getGoEnv(), "config.yaml"))
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	conf = new(Config)
	if err := viper.Unmarshal(conf); err != nil {
		panic(err)
	}

	conf.Env = getGoEnv()
	fmt.Printf("%#v", conf)

	viper.OnConfigChange(func(in fsnotify.Event) {
		logger := zap.L()
		logger.Info("Config file changed", zap.String("file", in.Name))

		newConf := new(Config)
		if err := viper.Unmarshal(newConf); err != nil {
			logger.Error("Failed to unmarshal config", zap.Error(err))
			return
		}

		conf.mu.RLock()
		newConf.observers = conf.observers
		conf.mu.RUnlock()

		oldConf := conf
		conf = newConf

		if oldConf.Server.Addr != conf.Server.Addr {
			for _, observer := range conf.observers {
				observer.OnConfigChange(conf, ServerChange)
			}
		}
	})

	viper.WatchConfig()
}

func getGoEnv() string {
	env := os.Getenv("GO_ENV")
	if env == "" {
		return "test"
	}

	return env
}
