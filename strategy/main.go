package main

import (
	"context"
	"controller/api"
	xhttp "controller/pkg/http"
	"controller/pkg/logger"
	utilruntime "controller/pkg/runtime"
	"controller/pkg/tracing"
	api2 "controller/strategy/api"
	"controller/strategy/config"
	"controller/strategy/database"
	"controller/strategy/services"
	"controller/strategy/version"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"net"
	"net/http"
	"os"
	"runtime"
)

const LocalServiceId = "strategy"

var mainCmd = &cobra.Command{Use: LocalServiceId}

func main() {
	runtime.GOMAXPROCS(8 * runtime.NumCPU())
	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		http.ListenAndServe("0.0.0.0:8982", nil)
	}()

	mainCmd.AddCommand(version.Cmd())
	if mainCmd.Execute() != nil {
		os.Exit(1)
	}

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

	if err := database.InitDB(c.DB.Url, c.DB.DbUser, c.DB.DbPassword, c.DB.DbName, c.DB.Timeout); err != nil {
		logger.Errorf("init database: %v fail: %v", c.DB, err.Error())
		return nil
	}

	//init logger
	logger.InitLogger(c.Logger.Level)
	logger.Info("server config:", *c)

	//Set system resources to maximum
	if err := api.SetMaxResources(); err != nil {
		logger.Error(err)
	}
	xhttp.ReqConfigInit(c.Request.Max, c.Request.TimeOut)

	//init trace jaeger
	jaeger := tracing.SetupJaegerTracing("controller"+LocalServiceId, c.Jaeger.Url)
	if jaeger == nil {
		logger.Errorf("init jaeger fail ")
	} else {
		logger.Errorf("init jaeger success ")
	}
	defer func() {
		if jaeger != nil {
			jaeger.Flush()
		}
	}()

	//database.DB.Init(&c.DB)

	//init api server
	strategy := new(services.Service)
	strategy.Init(database.Db)
	initHttpRouter(c, strategy)

	return nil
}

func initHttpRouter(c *config.ServerConfig, strategy *services.Service) {

	router := mux.NewRouter().SkipClean(true).UseEncodedPath()

	// Add healthcheck router
	api2.RegisterHealthCheckRouter(router)

	//Add server metrics router
	api2.RegisterMetricsRouter(router)
	// Add API router.
	api2.RegisterAPIRouter(router, strategy)

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
