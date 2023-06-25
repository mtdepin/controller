package api

import (
	"controller/api"
	"controller/business/param"
	"controller/business/services"
	"controller/business/utils"
	"controller/common/auth"
	error2 "controller/pkg/ctlerror"
	"controller/pkg/logger"
	"go.opencensus.io/trace"
	"net/http"
	"time"
)

const (
	BusinessVersion = "v1"
	InvalidParam    = "InvalidParam"
)

type apiHandlers struct {
	service *services.Service
}

func (h *apiHandlers) Search(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(r.Context(), "Search")
	defer span.End()

	if r.Method == http.MethodOptions {
		api.WriteResponse(w, []byte(""))
		return
	}

	request := &param.SearchFileRequest{Ext: &param.Extend{Ctx: ctx}}
	if err := api.ParseParam(w, r, request); err != nil {
		return
	}

	if request.FileName == "" && request.OrderId == "" {
		api.ProcessFail(w, r, request, InvalidParam, error2.ErrorCodes.ToAPIErr(error2.ErrInvalidArguments))
		return
	}

	rsp, err := h.service.Search(request)
	if err != nil {
		api.ProcessFail(w, r, request, err.Error(), error2.ErrorCodes.ToAPIErr(error2.ErrInternalError))
		return
	}

	api.WriteResponse(w, rsp)
}

func (h *apiHandlers) UploadTask(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(r.Context(), "uploadtask")
	defer span.End()

	if r.Method == http.MethodOptions {
		api.WriteResponse(w, []byte(""))
		return
	}

	request := &api.UploadTaskRequest{Ext: &api.Extend{Ctx: ctx}}
	if err := api.ParseParam(w, r, request); err != nil {
		return
	}

	if request.RequestId == "" || request.UserId == "" || request.PieceNum < 1 {
		api.ProcessFail(w, r, request, InvalidParam, error2.ErrorCodes.ToAPIErr(error2.ErrInvalidArguments))
		return
	}

	a := r.Header.Get("Authorization")
	secret, err := utils.GetUserSecret(request.UserId, &param.Extend{Ctx: request.Ext.Ctx})
	if err != nil || !auth.AuthCheck(a, secret) {
		api.WriteResponse(w, []byte("签名认证未通过"))
		return
	}

	rsp, err := h.service.UploadTask(request)
	if err != nil {
		api.ProcessFail(w, r, request, err.Error(), error2.ErrorCodes.ToAPIErr(error2.ErrInternalError))
		return
	}

	api.WriteSuccessResponseObject(w, rsp)
}

func (h *apiHandlers) UploadFinish(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(r.Context(), "UploadFinish")
	defer span.End()

	if r.Method == http.MethodOptions {
		api.WriteResponse(w, []byte(""))
		return
	}

	t1 := time.Now().UnixMilli()
	request := &param.UploadFinishRequest{Ext: &param.Extend{Ctx: ctx}}
	if err := api.ParseParam(w, r, request); err != nil {
		return
	}
	if request.OrderId == "" || len(request.Tasks) == 0 {
		api.ProcessFail(w, r, request, InvalidParam, error2.ErrorCodes.ToAPIErr(error2.ErrInvalidArguments))
		return
	}

	rsp, err := h.service.UploadFinish(request)
	t2 := time.Now().UnixMilli()
	logger.Infof("UploadFinish, orderId: %v, begin_time: %v,  endtime: %v, total_cost: %v ms", request.OrderId, t1, t2, t2-t1)

	if err != nil {
		api.ProcessFail(w, r, request, err.Error(), error2.ErrorCodes.ToAPIErr(error2.ErrInternalError))
		return
	}

	api.WriteSuccessResponseObject(w, rsp)
}

func (h *apiHandlers) DownloadTask(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(r.Context(), "DownloadTask")
	defer span.End()

	if r.Method == http.MethodOptions {
		api.WriteResponse(w, []byte(""))
		return
	}

	request := &param.DownloadTaskRequest{Ext: &param.Extend{Ctx: ctx}}
	if err := api.ParseParam(w, r, request); err != nil {
		return
	}
	if request.RequestId == "" || len(request.Tasks) == 0 {
		api.ProcessFail(w, r, request, InvalidParam, error2.ErrorCodes.ToAPIErr(error2.ErrInvalidArguments))
		return
	}

	a := r.Header.Get("Authorization")
	secret, err := utils.GetUserSecret(request.UserId, request.Ext)
	if err != nil || !auth.AuthCheck(a, secret) {
		api.WriteResponse(w, []byte("签名认证未通过"))
		return
	}

	rsp, err := h.service.DownloadTask(request)
	if err != nil {
		api.ProcessFail(w, r, request, err.Error(), error2.ErrorCodes.ToAPIErr(error2.ErrInternalError))
		return
	}

	api.WriteSuccessResponseObject(w, rsp)

}

func (h *apiHandlers) DownloadFinish(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(r.Context(), "DownloadFinish")
	defer span.End()

	if r.Method == http.MethodOptions {
		api.WriteResponse(w, []byte(""))
		return
	}

	request := &param.DownloadFinishRequest{Ext: &param.Extend{Ctx: ctx}}
	if err := api.ParseParam(w, r, request); err != nil {
		return
	}
	if request.OrderId == "" || len(request.Tasks) == 0 {
		api.ProcessFail(w, r, request, InvalidParam, error2.ErrorCodes.ToAPIErr(error2.ErrInvalidArguments))
		return
	}

	rsp, err := h.service.DownloadFinish(request)
	if err != nil {
		api.ProcessFail(w, r, request, err.Error(), error2.ErrorCodes.ToAPIErr(error2.ErrInternalError))
		return
	}

	api.WriteSuccessResponseObject(w, rsp)
}

func (h *apiHandlers) DeleteFid(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(r.Context(), "DownloadFinish")
	defer span.End()

	if r.Method == http.MethodOptions {
		api.WriteResponse(w, []byte(""))
		return
	}

	request := &param.DeleteFidRequest{Ext: &param.Extend{Ctx: ctx}}
	if err := api.ParseParam(w, r, request); err != nil {
		return
	}
	if request.OrderId == "" || len(request.Fids) == 0 {
		api.ProcessFail(w, r, request, InvalidParam, error2.ErrorCodes.ToAPIErr(error2.ErrInvalidArguments))
		return
	}

	rsp, err := h.service.DeleteFid(request)
	if err != nil {
		api.ProcessFail(w, r, request, err.Error(), error2.ErrorCodes.ToAPIErr(error2.ErrInternalError))
		return
	}

	api.WriteSuccessResponseObject(w, rsp)
}

/*func (h *apiHandlers) GetUploadKNodes(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(r.Context(), "GetUploadKNodes")
	defer span.End()

	if r.Method == http.MethodOptions {
		api.WriteResponse(w, []byte(""))
		return
	}

	request := &kepler_api.GetKNodesRequest{Ext: &kepler_api.Extend{Ctx: ctx}}
	if err := api.ParseParam(w, r, request); err != nil {
		return
	}

	rsp, err := h.service.GetUploadKNodes(request)
	if err != nil {
		api.ProcessFail(w, r, request, err.Error(), error2.ErrorCodes.ToAPIErr(error2.ErrInternalError))
		return
	}

	api.WriteSuccessResponseObject(w, rsp)
}*/

func (h *apiHandlers) UploadPieceFid(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(r.Context(), "UploadPieceFid")
	defer span.End()

	if r.Method == http.MethodOptions {
		api.WriteResponse(w, []byte(""))
		return
	}

	request := &api.UploadPieceFidRequest{Ext: &api.Extend{Ctx: ctx}}
	if err := api.ParseParam(w, r, request); err != nil {
		return
	}

	rsp, err := h.service.UploadPieceFid(request)
	if err != nil {
		api.ProcessFail(w, r, request, err.Error(), error2.ErrorCodes.ToAPIErr(error2.ErrInternalError))
		return
	}

	api.WriteSuccessResponseObject(w, rsp)
}
