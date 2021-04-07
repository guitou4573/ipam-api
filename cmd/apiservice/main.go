package apiservice

import (
  "github.com/guitou4573/ipam/pkg/router"
  "github.com/guitou4573/ipam/pkg/stores"
)

var (
	httpRouter *router.Router
)

type ApiServiceOptions struct {
	SubnetStore *stores.SubnetStore
}

func Execute(options ApiServiceOptions) {
  httpRouter = router.NewRouter(router.RouterOptions{
		SubnetStore: options.SubnetStore,
	})
	httpRouter.Start()
}
