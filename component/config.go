package component

import (
	"flag"
	sctx "github.com/viettranx/service-context"
)

type config struct {
	id             string
	urlRPCCategory string
	catgRPCServer  string
	gRPCServerPort int
}

func NewConfigComponent(id string) *config {
	return &config{id: id}
}

func (comp *config) ID() string {
	return comp.id
}

func (comp *config) InitFlags() {
	flag.StringVar(
		&comp.urlRPCCategory,
		"url-rpc-category",
		"http://localhost:3000/v1/category/rpc",
		"URL of category RPC",
	)

	flag.StringVar(
		&comp.catgRPCServer,
		"cat-grpc-server-url",
		"localhost:8080",
		"URL of category gRPC server",
	)

	flag.IntVar(
		&comp.gRPCServerPort,
		"grpc-server-port",
		8080,
		"Port of gRPC server",
	)

}

func (comp *config) Activate(context sctx.ServiceContext) error {
	return nil
}

func (comp *config) Stop() error {
	return nil
}

func (comp *config) GetGRPCServer() string { return comp.catgRPCServer }
func (comp *config) GetGRPCServPort() int  { return comp.gRPCServerPort }
