package api

import (
	"controller/api"
	error2 "controller/pkg/ctlerror"
	"controller/strategy/param"
	"controller/strategy/services"
	"go.opencensus.io/trace"
	"net/http"
)

const (
	StrategyVersion = "v1"
	InvalidParam    = "InvalidParam"
)

type apiHandlers struct {
	service *services.Service
}

func (h *apiHandlers) CreateStrategy(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(r.Context(), "CreateStrategy")
	defer span.End()

	request := &param.CreateStrategyRequest{Ext: &param.Extend{Ctx: ctx}}
	if err := api.ParseParam(w, r, request); err != nil {
		return
	}
	/*if request.RequestId == "" || request.OrderId == "" {
		api.ProcessFail(w, r, request, InvalidParam, error2.ErrorCodes.ToAPIErr(error2.ErrInvalidArguments))
		return
	}*/

	rsp, err := h.service.CreateStrategy(request)

	if err != nil {
		api.ProcessFail(w, r, request, err.Error(), error2.ErrorCodes.ToAPIErr(error2.ErrInternalError))
		return
	}

	api.WriteSuccessResponseObject(w, rsp)

}

func (h *apiHandlers) GetReplicateStrategy(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(r.Context(), "GetReplicateStrategy")
	defer span.End()

	request := &param.GetStrategyRequset{Ext: &param.Extend{Ctx: ctx}}
	if err := api.ParseParam(w, r, request); err != nil {
		return
	}
	if request.OrderId == "" {
		api.ProcessFail(w, r, request, InvalidParam, error2.ErrorCodes.ToAPIErr(error2.ErrInvalidArguments))
		return
	}

	rsp, err := h.service.GetReplicateStrategy(request)
	if err != nil {
		api.ProcessFail(w, r, request, err.Error(), error2.ErrorCodes.ToAPIErr(error2.ErrInternalError))
		return
	}

	api.WriteSuccessResponseObject(w, rsp)
}

func (h *apiHandlers) GetOrderDeleteStrategy(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(r.Context(), "GetOrderDeleteStrategy")
	defer span.End()

	request := &param.GetStrategyRequset{Ext: &param.Extend{Ctx: ctx}}
	if err := api.ParseParam(w, r, request); err != nil {
		return
	}
	if request.OrderId == "" {
		api.ProcessFail(w, r, request, InvalidParam, error2.ErrorCodes.ToAPIErr(error2.ErrInvalidArguments))
		return
	}

	rsp, err := h.service.GetOrderDeleteStrategy(request)
	if err != nil {
		api.ProcessFail(w, r, request, err.Error(), error2.ErrorCodes.ToAPIErr(error2.ErrInternalError))
		return
	}

	api.WriteSuccessResponseObject(w, rsp)
}

func (h *apiHandlers) GetFidDeleteStrategy(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(r.Context(), "GetFidDeleteStrategy")
	defer span.End()

	request := &param.GetFidDeleteStrategyRequest{Ext: &param.Extend{Ctx: ctx}}
	if err := api.ParseParam(w, r, request); err != nil {
		return
	}
	if request.OrderId == "" {
		api.ProcessFail(w, r, request, InvalidParam, error2.ErrorCodes.ToAPIErr(error2.ErrInvalidArguments))
		return
	}

	rsp, err := h.service.GetFidDeleteStrategy(request)
	if err != nil {
		api.ProcessFail(w, r, request, err.Error(), error2.ErrorCodes.ToAPIErr(error2.ErrInternalError))
		return
	}

	api.WriteSuccessResponseObject(w, rsp)
}
