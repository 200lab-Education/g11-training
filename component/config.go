package component

import (
	"flag"
	sctx "github.com/viettranx/service-context"
)

type config struct {
	id             string
	urlRPCCategory string
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
}

func (comp *config) Activate(context sctx.ServiceContext) error {
	return nil
}

func (comp *config) Stop() error {
	return nil
}

func (comp *config) GetURLRPCCategory() string { return comp.urlRPCCategory }
