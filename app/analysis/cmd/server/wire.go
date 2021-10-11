// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"poa-service/app/analysis/internal/biz"
	"poa-service/app/analysis/internal/conf"
	"poa-service/app/analysis/internal/data"
	"poa-service/app/analysis/internal/server"
	"poa-service/app/analysis/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// initApp init kratos application.
func initApp(*conf.Server, *conf.Data, *conf.Ext, *conf.ZeekLog, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
