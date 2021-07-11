package rest

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
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
	pattern := chi.URLParam(r, "pattern")
	deviceID, err := strconv.ParseInt(chi.URLParam(r, "deviceID"), 10, 64)
	if err != nil {
		SendError(w, r, invalidIDErr, http.StatusNotFound)
		return
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		SendUnexpectedErr(w, r)
	}

	invokeRequest := string(body[:])

	invokeResponse, err := gHandler.httpGrpcService.InvokeDevice(deviceID, pattern, invokeRequest)
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
