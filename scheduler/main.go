package main

import (
	"context"
	"controller/api"
	xhttp "controller/pkg/http"
	"controller/pkg/logger"
	utilruntime "controller/pkg/runtime"
	"controller/pkg/tracing"
	api2 "controller/scheduler/api"
	"controller/scheduler/config"
	"controller/scheduler/database"
	"controller/scheduler/services"
	"github.com/gorilla/mux"
	"net"
	"os"
)

const LocalServiceId = "scheduler"

func main() {
	serve()
}

func serve() error {
	defer utilruntime.HandleCrash()
	c, err := config.LoadServerConfig(LocalServiceId)
	if err != nil {
		logger.Error(err)
		return err
	}
	config.ServerCfg = c

	if err := database.InitDB(c.DB.Url, c.DB.DbName, c.DB.Timeout); err != nil {
		logger.Errorf("init database fail %v", c.DB)
		return nil
	}

	//init logger
	logger.InitLogger(c.Logger.Level)

	//Set system resources to maximum
	if err := api.SetMaxResources(); err != nil {
		logger.Error(err)
	}
	xhttp.ReqConfigInit(c.Request.Max, c.Request.TimeOut)

	//init trace jaeger
	jaeger := tracing.SetupJaegerTracing("mt_node_manager")
	defer func() {
		if jaeger != nil {
			jaeger.Flush()
		}
	}()

	//database.DB.Init(&c.DB)

	//init api server
	scheduler := new(services.Service)
	scheduler.Init(database.Db)

	initHttpRouter(c, scheduler)
	return nil
}

func initHttpRouter(c *config.ServerConfig, scheduler *services.Service) {

	router := mux.NewRouter().SkipClean(true).UseEncodedPath()

	// Add healthcheck router
	api2.RegisterHealthCheckRouter(router)

	//Add server metrics router
	api2.RegisterMetricsRouter(router)
	// Add API router.
	api2.RegisterAPIRouter(router, scheduler)

	// Use all the middlewares
	router.Use(api.GlobalHandlers...)

	ctx := context.Background()
	addr := c.Node.Api
	if addr == "" {
		addr = ":8521"
	}
	httpServer := xhttp.NewServer([]string{addr},
		router, nil)
	httpServer.BaseContext = func(listener net.Listener) context.Context {
		return ctx
	}

	globalHTTPServerErrorCh := make(chan error)

	go func() {
		logger.Infof("starting api Server : %s", addr)
		globalHTTPServerErrorCh <- httpServer.Start()
	}()

	select {
	case <-globalHTTPServerErrorCh:
		//todo: handler signals
		os.Exit(1)
	}
}
