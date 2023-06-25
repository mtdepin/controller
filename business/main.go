package main

import (
	"context"
	"controller/api"
	api2 "controller/business/api"
	"controller/business/config"
	"controller/business/database"
	"controller/business/services"
	xhttp "controller/pkg/http"
	"controller/pkg/logger"
	utilruntime "controller/pkg/runtime"
	"controller/pkg/tracing"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
)

const LocalServiceId = "business"

func main() {
	go func() {
		http.ListenAndServe("0.0.0.0:8888", nil)
	}()
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

	logger.Info("config Info: %v", *c)

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
	business := new(services.Service)
	business.Init(database.Db)

	initHttpRouter(c, business)
	return nil
}

func initHttpRouter(c *config.ServerConfig, business *services.Service) {

	router := mux.NewRouter().SkipClean(true).UseEncodedPath()

	// Add healthcheck router
	api2.RegisterHealthCheckRouter(router)

	//Add server metrics router
	api2.RegisterMetricsRouter(router)
	// Add API router.
	api2.RegisterAPIRouter(router, business)

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
