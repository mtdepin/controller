package api

import (
	"controller/api"
	error2 "controller/pkg/ctlerror"
	"controller/task_tracker/param"
	"controller/task_tracker/services"
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
	request := &param.CreateTaskRequest{}
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

	api.WriteSuccessResponseObject(w, rsp)

}

func (h *apiHandlers) UploadFinish(w http.ResponseWriter, r *http.Request) {
	request := &param.UploadFinishRequest{}
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
	request := &param.CallbackUploadRequest{}
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
	request := &param.CallbackRepRequest{}
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
	request := &param.CallbackDeleteRequest{}
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
	request := &param.CallbackChargeRequest{}
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
	request := &param.DownloadFinishRequest{}
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
