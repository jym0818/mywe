package ioc

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	type Config struct {
		dsn string `yaml:"dsn"`
	}
	var cfg Config
	err := viper.UnmarshalKey("mysql", &cfg)
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(mysql.Open(cfg.dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
