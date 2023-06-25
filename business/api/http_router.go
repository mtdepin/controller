package api

import (
	"controller/api"
	"controller/business/services"
	xhttp "controller/pkg/http"
	"controller/pkg/logger"
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
	apiRouter := router.PathPrefix("/business/" + BusinessVersion).Subrouter()

	gz, err := gzhttp.NewWrapper(gzhttp.MinSize(1000), gzhttp.CompressionLevel(gzip.BestSpeed))
	if err != nil {
		logger.Fatal(err, "Unable to initialize server")
	}
	maxClients := xhttp.MaxClients

	apiRouter.Methods(http.MethodGet).Path("/search").HandlerFunc(
		maxClients(gz(api.HttpTraceAll(apiHandlers.Search))))

	//test
	apiRouter.Methods(http.MethodOptions).Path("/search").HandlerFunc(
		maxClients(gz(api.HttpTraceAll(apiHandlers.Search))))
	//

	apiRouter.Methods(http.MethodPost).Path("/uploadTask").HandlerFunc(
		maxClients(gz(api.HttpTraceAll(apiHandlers.UploadTask))))

	//test
	apiRouter.Methods(http.MethodOptions).Path("/uploadTask").HandlerFunc(
		maxClients(gz(api.HttpTraceAll(apiHandlers.UploadTask))))
	//

	apiRouter.Methods(http.MethodPost).Path("/uploadFinish").HandlerFunc(
		maxClients(gz(api.HttpTraceAll(apiHandlers.UploadFinish))))

	//test
	apiRouter.Methods(http.MethodOptions).Path("/uploadFinish").HandlerFunc(
		maxClients(gz(api.HttpTraceAll(apiHandlers.UploadFinish))))
	//

	apiRouter.Methods(http.MethodPost).Path("/downloadTask").HandlerFunc(
		maxClients(gz(api.HttpTraceAll(apiHandlers.DownloadTask))))

	//test code
	apiRouter.Methods(http.MethodOptions).Path("/downloadTask").HandlerFunc(
		maxClients(gz(api.HttpTraceAll(apiHandlers.DownloadTask))))
	//

	apiRouter.Methods(http.MethodPost).Path("/downloadFinish").HandlerFunc(
		maxClients(gz(api.HttpTraceAll(apiHandlers.DownloadFinish))))

	//test code
	apiRouter.Methods(http.MethodOptions).Path("/downloadFinish").HandlerFunc(
		maxClients(gz(api.HttpTraceAll(apiHandlers.DownloadFinish))))
	//

	apiRouter.Methods(http.MethodPost).Path("/deleteFid").HandlerFunc(
		maxClients(gz(api.HttpTraceAll(apiHandlers.DeleteFid))))

	//test code
	apiRouter.Methods(http.MethodOptions).Path("/deleteFid").HandlerFunc(
		maxClients(gz(api.HttpTraceAll(apiHandlers.DeleteFid))))
	//

	/*apiRouter.Methods(http.MethodGet).Path("/getUploadKNodes").HandlerFunc(
		maxClients(gz(api.HttpTraceAll(apiHandlers.GetUploadKNodes))))

	//test code
	apiRouter.Methods(http.MethodOptions).Path("/getUploadKNodes").HandlerFunc(
		maxClients(gz(api.HttpTraceAll(apiHandlers.GetUploadKNodes))))
	*/
	//

	apiRouter.Methods(http.MethodPost).Path("/uploadPieceFid").HandlerFunc(
		maxClients(gz(api.HttpTraceAll(apiHandlers.UploadPieceFid))))

	//test code
	apiRouter.Methods(http.MethodOptions).Path("/uploadPieceFid").HandlerFunc(
		maxClients(gz(api.HttpTraceAll(apiHandlers.UploadPieceFid))))
	//

}
