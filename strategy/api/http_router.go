package api

import (
	"controller/api"
	xhttp "controller/pkg/http"
	"controller/pkg/logger"
	"controller/strategy/services"
	"github.com/gorilla/mux"
	"github.com/klauspost/compress/gzhttp"
	"github.com/klauspost/compress/gzip"
	"net/http"
)

func RegisterHealthCheckRouter(router *mux.Router) {

	//todo
}

func RegisterMetricsRouter(router *mux.Router) {
	//todo
}

func RegisterAPIRouter(router *mux.Router, ck *services.Service) {
	apiHandlers := apiHandlers{
		service: ck,
	}

	// API Router
	apiRouter := router.PathPrefix("/strategy/" + StrategyVersion).Subrouter()

	gz, err := gzhttp.NewWrapper(gzhttp.MinSize(1000), gzhttp.CompressionLevel(gzip.BestSpeed))
	if err != nil {
		logger.Fatal(err, "Unable to initialize server")
	}
	maxClients := xhttp.MaxClients

	apiRouter.Methods(http.MethodPost).Path("/createStrategy").HandlerFunc(
		maxClients(gz(api.HttpTraceAll(apiHandlers.CreateStrategy))))

	apiRouter.Methods(http.MethodGet).Path("/getReplicateStrategy").HandlerFunc(
		maxClients(gz(api.HttpTraceAll(apiHandlers.GetReplicateStrategy))))

	apiRouter.Methods(http.MethodGet).Path("/getDeleteReplicateStrategy").HandlerFunc( // 上传文件失败，获取删除策略
		maxClients(gz(api.HttpTraceAll(apiHandlers.GetDeleteReplicateStrategy))))
}
