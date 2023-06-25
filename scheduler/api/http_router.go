package api

import (
	"controller/api"
	xhttp "controller/pkg/http"
	"controller/pkg/logger"
	"controller/scheduler/services"
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
	apiRouter := router.PathPrefix("/scheduler/" + SchedulerVersion).Subrouter()

	gz, err := gzhttp.NewWrapper(gzhttp.MinSize(1000), gzhttp.CompressionLevel(gzip.BestSpeed))
	if err != nil {
		logger.Fatal(err, "Unable to initialize server")
	}
	maxClients := xhttp.MaxClients

	apiRouter.Methods(http.MethodPost).Path("/replicate").HandlerFunc(
		maxClients(gz(api.HttpTraceAll(apiHandlers.Replicate))))

	apiRouter.Methods(http.MethodPost).Path("/charge").HandlerFunc(
		maxClients(gz(api.HttpTraceAll(apiHandlers.Charge))))

	apiRouter.Methods(http.MethodPost).Path("/delete").HandlerFunc(
		maxClients(gz(api.HttpTraceAll(apiHandlers.Delete))))

	apiRouter.Methods(http.MethodGet).Path("/searchRep").HandlerFunc(
		maxClients(gz(api.HttpTraceAll(apiHandlers.SearchRep))))
}
