package main

import (
	"context"
	"controller/api"
	xhttp "controller/pkg/http"
	"controller/pkg/logger"
	"controller/pkg/montior"
	utilruntime "controller/pkg/runtime"
	"controller/pkg/tracing"
	api2 "controller/task_tracker/api"
	"controller/task_tracker/config"
	"controller/task_tracker/database"
	"controller/task_tracker/services"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
)

const LocalServiceId = "task_tracker"

func main() {
	go func() {
		http.ListenAndServe("0.0.0.0:8889", nil)
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

	logger.Infof("config: %v,  schedulerUrl: %v", *config.ServerCfg, config.ServerCfg.Scheduler.Url)

	ding := new(montior.DingTalk)
	ding.Init("task_tracker", config.ServerCfg.Montior.AccessToken, config.ServerCfg.Montior.Secret)
	//init logger
	logger.InitLoggerWithDingTalk(c.Logger.Level, ding)

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

	if err := database.InitDB(c.DB.Url, c.DB.DbName, c.DB.Timeout); err != nil {
		logger.Errorf("init database fail %v", c.DB)
		return nil
	}

	//init api server
	taskTracker := new(services.Service)
	taskTracker.Init(database.Db)

	initHttpRouter(c, taskTracker)
	logger.Info("init task_tracker success")
	return nil
}

func initHttpRouter(c *config.ServerConfig, taskTracker *services.Service) {

	router := mux.NewRouter().SkipClean(true).UseEncodedPath()

	// Add healthcheck router
	api2.RegisterHealthCheckRouter(router)

	//Add server metrics router
	api2.RegisterMetricsRouter(router)
	// Add API router.
	api2.RegisterAPIRouter(router, taskTracker)

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
