//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/jym0818/mywe/ioc"
)

func InitWebServer() *App {
	wire.Build(
		ioc.InitWeb,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
