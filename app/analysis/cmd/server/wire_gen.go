// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"poa-service/app/analysis/internal/biz"
	"poa-service/app/analysis/internal/conf"
	"poa-service/app/analysis/internal/data"
	"poa-service/app/analysis/internal/server"
	"poa-service/app/analysis/internal/service"
)

// Injectors from wire.go:

// initApp init kratos application.
func initApp(confServer *conf.Server, confData *conf.Data, ext *conf.Ext, zeekLog *conf.ZeekLog, logger log.Logger) (*kratos.App, func(), error) {
	dataData, cleanup, err := data.NewData(confData, logger)
	if err != nil {
		return nil, nil, err
	}
	analysisRepo := data.NewAnalysisRepo(dataData, logger)
	externalAnalysisRepo := data.NewExtAnalysis(ext)
	analysisUseCase := biz.NewAnalysisUseCase(analysisRepo, externalAnalysisRepo, logger)
	analysisInterfaceService := service.NewAnalysisInterfaceService(analysisUseCase, logger)
	opinionRepo := data.NewOpinionRepo(dataData, logger)
	zeekLogRepo := data.NewZeekLogRepo(dataData, logger)
	refactorUseCase := biz.NewRefactorUseCase(zeekLog, opinionRepo, zeekLogRepo, logger)
	refactorInterfaceService := service.NewRefactorInterfaceService(refactorUseCase, logger)
	httpServer := server.NewHTTPServer(confServer, analysisInterfaceService, refactorInterfaceService, logger)
	grpcServer := server.NewGRPCServer(confServer, analysisInterfaceService, logger)
	app := newApp(logger, httpServer, grpcServer)
	return app, func() {
		cleanup()
	}, nil
}
