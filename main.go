package main

import (
	"github.com/spf13/viper"
)

func main() {
	initViper()

	app := InitWebServer()
	for _, c := range app.consumers {
		er := c.Start()
		if er != nil {
			panic(er)
		}
	}
	app.server.Run(":8080")
}

func initViper() {
	viper.SetConfigFile("./config/dev.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
