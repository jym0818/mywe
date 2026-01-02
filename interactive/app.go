package main

import (
	"github.com/jym0818/mywe/interactive/events"
	"github.com/jym0818/mywe/pkg/grpcx"
)

type App struct {
	server    *grpcx.Server
	consumers []events.Consumer
}
