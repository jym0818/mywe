//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/jym0818/mywe/internal/repository"
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
)

func InitWebServer() *App {
	wire.Build(
		user,
		ioc.InitWeb,
		ioc.InitDB,
		ioc.InitRedis,
		ioc.InitHandler,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
