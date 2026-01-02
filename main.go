package main

import (
	"time"

	"github.com/spf13/viper"
)

func main() {
	initViper()
	//initPrometheus()
	app := InitWebServer()

	app.cron.Start()

	app.server.Run(":8080")
	// 这边可以考虑超时强制退出，防止有些任务，执行特别长的时间
	tm := time.NewTimer(time.Minute * 10)
	stop := app.cron.Stop()
	select {
	case <-stop.Done():
	case <-tm.C:

	}
}

func initViper() {
	viper.SetConfigFile("./config/dev.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

//func initPrometheus() {
//	http.Handle("/metrics", promhttp.Handler())
//	http.ListenAndServe(":8081", nil)
//}
