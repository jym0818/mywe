package ioc

import (
	grpc2 "github.com/jym0818/mywe/interactive/grpc"
	"github.com/jym0818/mywe/pkg/grpcx"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func InitGRPCxServer(intr *grpc2.InteractiveServiceServer) *grpcx.Server {
	type Config struct {
		Addr string `yaml:"addr"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.server", &cfg)
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	intr.Register(s)

	return &grpcx.Server{
		Server: s,
		Addr:   cfg.Addr,
	}
}
