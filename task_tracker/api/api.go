package api

import (
	"controller/api"
	error2 "controller/pkg/ctlerror"
	"controller/task_tracker/param"
	"controller/task_tracker/services"
	"go.opencensus.io/trace"
	"net/http"
)

const (
	TaskTrackerVersion = "v1"
	InvalidParam       = "InvalidParam"
)

type apiHandlers struct {
	service *services.Service
}

func (h *apiHandlers) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(r.Context(), "CreateTask")
	defer span.End()

	request := &api.CreateTaskRequest{Ext: &api.Extend{Ctx: ctx}}
	if err := api.ParseParam(w, r, request); err != nil {
		return
	}
	if request.RequestId == "" {
		api.ProcessFail(w, r, request, InvalidParam, error2.ErrorCodes.ToAPIErr(error2.ErrInvalidArguments))
		return
	}

	rsp, err := h.service.CreateTask(request)
	if err != nil {
		api.ProcessFail(w, r, request, err.Error(), error2.ErrorCodes.ToAPIErr(error2.ErrInternalError))
		return
	}
	//to do register metrics.
	api.WriteSuccessResponseObject(w, rsp)
}

func (h *apiHandlers) UploadFinish(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(r.Context(), "UploadFinish")
	defer span.End()

	request := &param.UploadFinishRequest{Ext: &param.Extend{Ctx: ctx}}
	if err := api.ParseParam(w, r, request); err != nil {
		return
	}
	if request.OrderId == "" || len(request.Tasks) == 0 {
		api.ProcessFail(w, r, request, InvalidParam, error2.ErrorCodes.ToAPIErr(error2.ErrInvalidArguments))
		return
	}

	rsp, err := h.service.UploadFinish(request)
	if err != nil {
		api.ProcessFail(w, r, request, err.Error(), error2.ErrorCodes.ToAPIErr(error2.ErrInternalError))
		return
	}

	api.WriteSuccessResponseObject(w, rsp)
}

func (h *apiHandlers) CallbackUpload(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(r.Context(), "CallbackUpload")
	defer span.End()

	request := &param.CallbackUploadRequest{Ext: &param.Extend{Ctx: ctx}}
	if err := api.ParseParam(w, r, request); err != nil {
		return
	}
	if request.OrderId == "" {
		api.ProcessFail(w, r, request, InvalidParam, error2.ErrorCodes.ToAPIErr(error2.ErrInvalidArguments))
		return
	}

	rsp, err := h.service.CallbackUpload(request)
	if err != nil {
		api.ProcessFail(w, r, request, err.Error(), error2.ErrorCodes.ToAPIErr(error2.ErrInternalError))
		return
	}

	api.WriteSuccessResponseObject(w, rsp)
}

func (h *apiHandlers) CallbackReplicate(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(r.Context(), "CallbackReplicate")
	defer span.End()

	request := &param.CallbackRepRequest{Ext: &param.Extend{Ctx: ctx}}
	if err := api.ParseParam(w, r, request); err != nil {
		return
	}
	if request.OrderId == "" {
		api.ProcessFail(w, r, request, InvalidParam, error2.ErrorCodes.ToAPIErr(error2.ErrInvalidArguments))
		return
	}

	rsp, err := h.service.CallbackReplicate(request)
	if err != nil {
		api.ProcessFail(w, r, request, err.Error(), error2.ErrorCodes.ToAPIErr(error2.ErrInternalError))
		return
	}

	api.WriteSuccessResponseObject(w, rsp)
}

func (h *apiHandlers) CallbackDelete(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(r.Context(), "CallbackDelete")
	defer span.End()

	request := &param.CallbackDeleteRequest{Ext: &param.Extend{Ctx: ctx}}
	if err := api.ParseParam(w, r, request); err != nil {
		return
	}
	if request.OrderId == "" {
		api.ProcessFail(w, r, request, InvalidParam, error2.ErrorCodes.ToAPIErr(error2.ErrInvalidArguments))
		return
	}

	rsp, err := h.service.CallbackDelete(request)
	if err != nil {
		api.ProcessFail(w, r, request, err.Error(), error2.ErrorCodes.ToAPIErr(error2.ErrInternalError))
		return
	}

	api.WriteSuccessResponseObject(w, rsp)
}

func (h *apiHandlers) CallbackCharge(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(r.Context(), "CallbackCharge")
	defer span.End()

	request := &param.CallbackChargeRequest{Ext: &param.Extend{Ctx: ctx}}
	if err := api.ParseParam(w, r, request); err != nil {
		return
	}
	if request.OrderId == "" {
		api.ProcessFail(w, r, request, InvalidParam, error2.ErrorCodes.ToAPIErr(error2.ErrInvalidArguments))
		return
	}

	rsp, err := h.service.CallbackCharge(request)
	if err != nil {
		api.ProcessFail(w, r, request, err.Error(), error2.ErrorCodes.ToAPIErr(error2.ErrInternalError))
		return
	}

	api.WriteSuccessResponseObject(w, rsp)
}

func (h *apiHandlers) DownloadFinish(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(r.Context(), "DownloadFinish")
	defer span.End()

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
	ctx, span := trace.StartSpan(r.Context(), "DeleteFid")
	defer span.End()

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

func (h *apiHandlers) CallbackDownload(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(r.Context(), "CallbackDownload")
	defer span.End()

	request := &param.CallbackDownloadRequest{Extend: &param.Extend{Ctx: ctx}}
	if err := api.ParseParam(w, r, request); err != nil {
		return
	}

	rsp, err := h.service.CallbackDownload(request)
	if err != nil {
		api.ProcessFail(w, r, request, err.Error(), error2.ErrorCodes.ToAPIErr(error2.ErrInternalError))
		return
	}

	api.WriteSuccessResponseObject(w, rsp)
}

func (h *apiHandlers) GetPageOrders(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(r.Context(), "GetOrders")
	defer span.End()

	request := &param.OrderPageQueryRequest{Ext: &param.Extend{Ctx: ctx}}
	if err := api.ParseParam(w, r, request); err != nil {
		return
	}

	rsp, err := h.service.GetPageOrders(request)
	if err != nil {
		api.ProcessFail(w, r, request, err.Error(), error2.ErrorCodes.ToAPIErr(error2.ErrInternalError))
		return
	}

	api.WriteSuccessResponseObject(w, rsp)
}

func (h *apiHandlers) GetOrderDetail(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(r.Context(), "GetOrders")
	defer span.End()

	request := &param.SearchOrderRequest{Ext: &param.Extend{Ctx: ctx}}
	if err := api.ParseParam(w, r, request); err != nil {
		return
	}

	rsp, err := h.service.GetOrderDetail(request)
	if err != nil {
		api.ProcessFail(w, r, request, err.Error(), error2.ErrorCodes.ToAPIErr(error2.ErrInternalError))
		return
	}

	api.WriteSuccessResponseObject(w, rsp)
}

func (h *apiHandlers) GetFidDetail(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(r.Context(), "GetFids")
	defer span.End()

	request := &param.SearchFidRequest{Ext: &param.Extend{Ctx: ctx}}
	if err := api.ParseParam(w, r, request); err != nil {
		return
	}

	rsp, err := h.service.GetFidDetail(request)
	if err != nil {
		api.ProcessFail(w, r, request, err.Error(), error2.ErrorCodes.ToAPIErr(error2.ErrInternalError))
		return
	}

	api.WriteSuccessResponseObject(w, rsp)
}

func (h *apiHandlers) GetPageFids(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(r.Context(), "GetFids")
	defer span.End()

	request := &param.FidPageQueryRequest{Ext: &param.Extend{Ctx: ctx}}
	if err := api.ParseParam(w, r, request); err != nil {
		return
	}

	rsp, err := h.service.GetPageFids(request)
	if err != nil {
		api.ProcessFail(w, r, request, err.Error(), error2.ErrorCodes.ToAPIErr(error2.ErrInternalError))
		return
	}

	api.WriteSuccessResponseObject(w, rsp)
}

func (h *apiHandlers) DebugStack(w http.ResponseWriter, r *http.Request) {
	api.WriteAllGoroutineStacks(w)
}

func (h *apiHandlers) UploadPieceFid(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(r.Context(), "GetFids")
	defer span.End()

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
