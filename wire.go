//go:build wireinject

package main

import (
	"github.com/google/wire"
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

func InitWebServer() *App {
	wire.Build(
		user,
		code,
		ioc.InitSMS,
		ioc.InitWeb,
		ioc.InitDB,
		ioc.InitRedis,
		ioc.InitHandler,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
