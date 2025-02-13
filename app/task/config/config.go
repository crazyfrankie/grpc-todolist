package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	once sync.Once
	conf *Config
)

type Config struct {
	Env    string
	Server Server `yaml:"server"`
	MySQL  MySQL  `yaml:"mysql"`
	ETCD   ETCD   `yaml:"etcd"`
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
		fmt.Println("Config file changed:", in.Name)
		if err := viper.ReadInConfig(); err != nil {
			fmt.Printf("Error reading config after change: %v\n", err)
			return
		}
		if err := viper.Unmarshal(conf); err != nil {
			fmt.Printf("Error unmarshalling config: %v\n", err)
			return
		}

		fmt.Printf("Config reloaded: %#v\n", conf)
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
