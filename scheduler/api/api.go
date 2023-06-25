package api

import (
	"controller/api"
	error2 "controller/pkg/ctlerror"
	"controller/scheduler/param"
	"controller/scheduler/services"
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
	request := &param.ReplicationRequest{}
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

func (h *apiHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	request := &param.DeleteOrderRequest{}
	if err := api.ParseParam(w, r, request); err != nil {
		return
	}
	if request.OrderId == "" || len(request.Tasks) == 0 {
		api.ProcessFail(w, r, request, InvalidParam, error2.ErrorCodes.ToAPIErr(error2.ErrInvalidArguments))
		return
	}

	rsp, err := h.service.Delete(request)
	if err != nil {
		api.ProcessFail(w, r, request, err.Error(), error2.ErrorCodes.ToAPIErr(error2.ErrInternalError))
		return
	}

	api.WriteSuccessResponseObject(w, rsp)

}

func (h *apiHandlers) Charge(w http.ResponseWriter, r *http.Request) {
	request := &param.ChargeRequest{}
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
