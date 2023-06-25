package api

import (
	"controller/api"
	error2 "controller/pkg/ctlerror"
	"controller/strategy/param"
	"controller/strategy/services"
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
	request := &param.CreateStrategyRequest{}
	if err := api.ParseParam(w, r, request); err != nil {
		return
	}
	if request.RequestId == "" || request.OrderId == "" {
		api.ProcessFail(w, r, request, InvalidParam, error2.ErrorCodes.ToAPIErr(error2.ErrInvalidArguments))
		return
	}

	rsp, err := h.service.CreateStrategy(request)

	if err != nil {
		api.ProcessFail(w, r, request, err.Error(), error2.ErrorCodes.ToAPIErr(error2.ErrInternalError))
		return
	}

	api.WriteSuccessResponseObject(w, rsp)

}

func (h *apiHandlers) GetReplicateStrategy(w http.ResponseWriter, r *http.Request) {
	request := &param.GetStrategyRequset{}
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

func (h *apiHandlers) GetDeleteReplicateStrategy(w http.ResponseWriter, r *http.Request) {
	request := &param.GetDeleteStrategyRequset{}
	if err := api.ParseParam(w, r, request); err != nil {
		return
	}
	if request.OrderId == "" {
		api.ProcessFail(w, r, request, InvalidParam, error2.ErrorCodes.ToAPIErr(error2.ErrInvalidArguments))
		return
	}

	rsp, err := h.service.GetDeleteReplicateStrategy(request)
	if err != nil {
		api.ProcessFail(w, r, request, err.Error(), error2.ErrorCodes.ToAPIErr(error2.ErrInternalError))
		return
	}

	api.WriteSuccessResponseObject(w, rsp)
}
