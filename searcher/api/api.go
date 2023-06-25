package api

import (
	"controller/api"
	error2 "controller/pkg/ctlerror"
	"controller/searcher/param"
	"controller/searcher/services"
	"net/http"
)

const (
	SearcherVersion = "v1"
	InvalidParam    = "InvalidParam"
)

type apiHandlers struct {
	service *services.Service
}

func (h *apiHandlers) Search(w http.ResponseWriter, r *http.Request) {
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

	api.WriteSuccessResponseObject(w, rsp)
}
