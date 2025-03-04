package ioc

import (
	"fmt"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/crazyfrankie/todolist/app/user/biz/repository/dao"
	"github.com/crazyfrankie/todolist/app/user/config"
)

var (
	dbOnce     sync.Once
	dbProvider *DBProvider
)

type DBProvider struct {
	db *gorm.DB
	mu sync.RWMutex
}

func NewDBProvider() *DBProvider {
	dbOnce.Do(func() {
		dbProvider = &DBProvider{}
		config.GetConf().AddObserver(dbProvider)
		dbProvider.refreshConnection()
	})
	return dbProvider
}

func (p *DBProvider) OnConfigChange(c *config.Config) {
	p.refreshConnection()
}

func (p *DBProvider) refreshConnection() {
	envCfg := config.GetEnvConfig()
	if envCfg == nil {
		panic("failed to get env config")
	}

	dsn := fmt.Sprintf(config.GetConf().MySQL.DSN,
		envCfg.MySQLUser,
		envCfg.MySQLPassword,
		envCfg.MySQLHost,
		envCfg.MySQLPort,
		envCfg.MySQLDb)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: &schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&dao.User{})

	p.mu.Lock()
	oldDB := p.db
	p.db = db
	p.mu.Unlock()

	if oldDB != nil {
		sqlDB, err := oldDB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
}

// GetDB 获取数据库连接
func (p *DBProvider) GetDB() *gorm.DB {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.db
}

func InitDB() *gorm.DB {
	return NewDBProvider().GetDB()
}
