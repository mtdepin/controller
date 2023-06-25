package api

import (
	"controller/api"
	error2 "controller/pkg/ctlerror"
	"controller/scheduler/param"
	"controller/scheduler/services"
	"go.opencensus.io/trace"
	"net/http"
)

const (
	SchedulerVersion = "v1"
	InvalidParam     = "InvalidParam"
)

type apiHandlers struct {
	service *services.Service
}

func (h *apiHandlers) Replicate(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(r.Context(), "Replicate")
	defer span.End()

	request := &param.ReplicationRequest{Ext: &param.Extend{Ctx: ctx}}
	if err := api.ParseParam(w, r, request); err != nil {
		return
	}
	if request.OrderId == "" || len(request.Tasks) < 1 {
		api.ProcessFail(w, r, request, InvalidParam, error2.ErrorCodes.ToAPIErr(error2.ErrInvalidArguments))
		return
	}

	rsp, err := h.service.Replicate(request)
	if err != nil {
		api.ProcessFail(w, r, request, err.Error(), error2.ErrorCodes.ToAPIErr(error2.ErrInternalError))
		return
	}

	api.WriteSuccessResponseObject(w, rsp)
}

func (h *apiHandlers) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(r.Context(), "DeleteOrder")
	defer span.End()

	request := &param.DeleteOrderRequest{Ext: &param.Extend{Ctx: ctx}}
	if err := api.ParseParam(w, r, request); err != nil {
		return
	}
	if request.OrderId == "" || len(request.Tasks) == 0 {
		api.ProcessFail(w, r, request, InvalidParam, error2.ErrorCodes.ToAPIErr(error2.ErrInvalidArguments))
		return
	}

	rsp, err := h.service.DeleteOrder(request)
	if err != nil {
		api.ProcessFail(w, r, request, err.Error(), error2.ErrorCodes.ToAPIErr(error2.ErrInternalError))
		return
	}

	api.WriteSuccessResponseObject(w, rsp)
}

func (h *apiHandlers) Charge(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(r.Context(), "Charge")
	defer span.End()

	request := &param.ChargeRequest{Ext: &param.Extend{Ctx: ctx}}
	if err := api.ParseParam(w, r, request); err != nil {
		return
	}
	if request.OrderId == "" || len(request.Tasks) == 0 {
		api.ProcessFail(w, r, request, InvalidParam, error2.ErrorCodes.ToAPIErr(error2.ErrInvalidArguments))
		return
	}

	rsp, err := h.service.Charge(request)
	if err != nil {
		api.ProcessFail(w, r, request, err.Error(), error2.ErrorCodes.ToAPIErr(error2.ErrInternalError))
		return
	}

	api.WriteSuccessResponseObject(w, rsp)
}

func (h *apiHandlers) SearchRep(w http.ResponseWriter, r *http.Request) {
	request := &param.UploadFinishOrder{}
	if err := api.ParseParam(w, r, request); err != nil {
		return
	}
	if request.OrderId == "" || len(request.Tasks) == 0 {
		api.ProcessFail(w, r, request, InvalidParam, error2.ErrorCodes.ToAPIErr(error2.ErrInvalidArguments))
		return
	}

	rsp, err := h.service.SearchRep(request)
	if err != nil {
		api.ProcessFail(w, r, request, err.Error(), error2.ErrorCodes.ToAPIErr(error2.ErrInternalError))
		return
	}

	api.WriteSuccessResponseObject(w, rsp)
}

func (h *apiHandlers) DeleteFid(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(r.Context(), "DeleteFid")
	defer span.End()

	request := &param.DeleteOrderFidRequest{Ext: &param.Extend{Ctx: ctx}}
	if err := api.ParseParam(w, r, request); err != nil {
		return
	}
	if request.OrderId == "" || len(request.Tasks) == 0 {
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
