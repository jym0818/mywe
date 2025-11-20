//go:build wireinject

package main

import (
	"github.com/google/wire"
	article2 "github.com/jym0818/mywe/internal/events/article"
	"github.com/jym0818/mywe/internal/repository"
	"github.com/jym0818/mywe/internal/repository/cache"
	"github.com/jym0818/mywe/internal/repository/dao"
	"github.com/jym0818/mywe/internal/service"
	"github.com/jym0818/mywe/internal/web"
	"github.com/jym0818/mywe/ioc"
)

var user = wire.NewSet(
	web.NewUserHandler,
	service.NewuserService,
	repository.NewuserRepository,
	dao.NewuserDAO,
	cache.NewuserCache,
)

var code = wire.NewSet(
	service.NewCodeService,
	repository.NewcodeRepository,
	cache.NewcodeCache)

var article = wire.NewSet(
	web.NewArticleHandler,
	service.NewarticleService,
	repository.NewarticleRepository,
	cache.NewarticleCache,
	dao.NewarticleDAO)

var interactive = wire.NewSet(
	service.NewinteractiveService,
	repository.NewinteractiveRepository,
	cache.NewinteractiveCache,
	dao.NewinteractiveDao,
)

func InitWebServer() *App {
	wire.Build(
		user,
		code,
		article,
		interactive,

		ioc.NewConsumers,
		ioc.InitKafka,
		ioc.NewSyncProducer,
		article2.NewKafkaProducer,
		article2.NewInteractiveReadEventConsumer,
		ioc.InitLogger,
		ioc.InitWechat,
		web.NewWechatHandler,
		ioc.InitSMS,
		ioc.InitWeb,
		ioc.InitDB,
		ioc.InitRedis,
		ioc.InitHandler,
		ioc.InitRatelimit,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
