package api

import (
	"controller/api"
	"controller/business/param"
	"controller/business/services"
	"controller/business/utils"
	"controller/common/auth"
	error2 "controller/pkg/ctlerror"
	"net/http"
)

const (
	BusinessVersion = "v1"
	InvalidParam    = "InvalidParam"
)

type apiHandlers struct {
	service *services.Service
}

func (h *apiHandlers) Search(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		api.WriteResponse(w, []byte(""))
		return
	}

	request := &param.SearchFileRequest{}
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
	if r.Method == http.MethodOptions {
		api.WriteResponse(w, []byte(""))
		return
	}

	request := &param.UploadTaskRequest{}
	if err := api.ParseParam(w, r, request); err != nil {
		return
	}
	if request.RequestId == "" || request.UserId == "" || len(request.Tasks) == 0 {
		api.ProcessFail(w, r, request, InvalidParam, error2.ErrorCodes.ToAPIErr(error2.ErrInvalidArguments))
		return
	}

	a := r.Header.Get("Authorization")
	secret, err := utils.GetUserSecret(request.UserId)
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
	if r.Method == http.MethodOptions {
		api.WriteResponse(w, []byte(""))
		return
	}

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

func (h *apiHandlers) DownloadTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		api.WriteResponse(w, []byte(""))
		return
	}

	request := &param.DownloadTaskRequest{}
	if err := api.ParseParam(w, r, request); err != nil {
		return
	}
	if request.RequestId == "" || len(request.Tasks) == 0 {
		api.ProcessFail(w, r, request, InvalidParam, error2.ErrorCodes.ToAPIErr(error2.ErrInvalidArguments))
		return
	}

	a := r.Header.Get("Authorization")
	secret, err := utils.GetUserSecret(request.UserId)
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
	if r.Method == http.MethodOptions {
		api.WriteResponse(w, []byte(""))
		return
	}

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
