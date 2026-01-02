package main

import (
	"github.com/google/wire"
	"github.com/jym0818/mywe/interactive/events"
	grpc2 "github.com/jym0818/mywe/interactive/grpc"
	"github.com/jym0818/mywe/interactive/ioc"
	"github.com/jym0818/mywe/interactive/repository"
	"github.com/jym0818/mywe/interactive/repository/cache"
	"github.com/jym0818/mywe/interactive/repository/dao"
	service2 "github.com/jym0818/mywe/interactive/service"
)

var interactive = wire.NewSet(
	service2.NewinteractiveService,
	repository.NewinteractiveRepository,
	dao.NewinteractiveDao,
	cache.NewinteractiveCache,
)

func InitGRPCServer() *App {
	wire.Build(
		interactive,
		ioc.InitDB,
		ioc.InitLogger,
		ioc.InitRedis,
		ioc.InitKafka,
		ioc.InitGRPCxServer,
		grpc2.NewInteractiveServiceServer,
		events.NewInteractiveReadEventConsumer,
		ioc.NewConsumers,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
