package rest

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/tnynlabs/wyrm/pkg/tunnels"
	"github.com/tnynlabs/wyrm/pkg/utils"
)

type GrpcHandler struct {
	httpGrpcService tunnels.Service
}

func CreateGrpcHandler(tService tunnels.Service) GrpcHandler {
	return GrpcHandler{tService}
}

func (gHandler *GrpcHandler) InvokeDevice(w http.ResponseWriter, r *http.Request) {
	deviceID, err := strconv.ParseInt(chi.URLParam(r, "deviceID"), 10, 64)
	if err != nil {
		SendError(w, r, invalidIDErr, http.StatusNotFound)
		return
	}

	invokeRequest := invokeRequestRest{}
	err = render.DecodeJSON(r.Body, invokeRequest)
	if err != nil {
		SendInvalidJSONErr(w, r)
		return
	}

	invokeResponse, err := gHandler.httpGrpcService.InvokeDevice(deviceID, invokeRequest.Pattern, invokeRequest.Data)
	if err != nil {
		serviceErr := utils.ToServiceErr(err)
		switch serviceErr.Code {
		case tunnels.ConnectionErrorCode:
			SendError(w, r, *serviceErr, http.StatusBadGateway)
		default:
			SendUnexpectedErr(w, r)
		}
	}

	result := &map[string]interface{}{
		"response": invokeResponse.Data,
	}

	SendResponse(w, r, result)
}

type invokeRequestRest struct {
	Pattern string `json:"pattern,omitempty"`
	Data    string `json:"data,omitempty"`
}
