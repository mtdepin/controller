package api

import (
	"controller/api"
	xhttp "controller/pkg/http"
	"controller/pkg/logger"
	"controller/task_tracker/services"
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
	apiRouter := router.PathPrefix("/task_tracker/" + TaskTrackerVersion).Subrouter()

	gz, err := gzhttp.NewWrapper(gzhttp.MinSize(1000), gzhttp.CompressionLevel(gzip.BestSpeed))
	if err != nil {
		logger.Fatal(err, "Unable to initialize server")
	}
	maxClients := xhttp.MaxClients
	//register router handler
	apiRouter.Methods(http.MethodPost).Path("/createTask").HandlerFunc(
		maxClients(gz(api.HttpTraceAll(apiHandlers.CreateTask))))

	apiRouter.Methods(http.MethodPost).Path("/uploadFinish").HandlerFunc(
		maxClients(gz(api.HttpTraceAll(apiHandlers.UploadFinish))))

	apiRouter.Methods(http.MethodPost).Path("/callbackUpload").HandlerFunc(
		maxClients(gz(api.HttpTraceAll(apiHandlers.CallbackUpload))))

	apiRouter.Methods(http.MethodPost).Path("/callbackReplicate").HandlerFunc(
		maxClients(gz(api.HttpTraceAll(apiHandlers.CallbackReplicate))))

	/*apiRouter.Methods(http.MethodPost).Path("/callbackDelete").HandlerFunc(
	maxClients(gz(api.HttpTraceAll(apiHandlers.CallbackDelete))))
	*/

	apiRouter.Methods(http.MethodPost).Path("/callbackCharge").HandlerFunc(
		maxClients(gz(api.HttpTraceAll(apiHandlers.CallbackCharge))))

	apiRouter.Methods(http.MethodPost).Path("/downloadFinish").HandlerFunc(
		maxClients(gz(api.HttpTraceAll(apiHandlers.DownloadFinish))))
}
