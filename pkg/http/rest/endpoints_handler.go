package rest

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/tnynlabs/wyrm/pkg/endpoints"
	"github.com/tnynlabs/wyrm/pkg/utils"
)

type EndpointHandler struct {
	endpointService endpoints.Service
}

func CreateEndpointHandler(epService endpoints.Service) EndpointHandler {
	return EndpointHandler{epService}
}

func (epHandler *EndpointHandler) Get(w http.ResponseWriter, r *http.Request) {
	endpointID, err := strconv.ParseInt(chi.URLParam(r, "endpointID"), 10, 64)
	if err != nil {
		SendError(w, r, invalidIDErr, http.StatusNotFound)
		return
	}

	endpoint, err := epHandler.endpointService.GetByID(endpointID)
	if err != nil {
		serviceErr := utils.ToServiceErr(err)
		switch serviceErr.Code {
		case endpoints.EndpointNotFoundCode:
			SendError(w, r, *serviceErr, http.StatusNotFound)
		default:
			SendUnexpectedErr(w, r)
		}
		return
	}

	endpointData := fromEndpoint(*endpoint)

	result := &map[string]interface{}{
		"endpoint": endpointData,
	}

	SendResponse(w, r, result)
}

func (epHandler *EndpointHandler) Create(w http.ResponseWriter, r *http.Request) {
	deviceID, err := strconv.ParseInt(chi.URLParam(r, "deviceID"), 10, 64)
	if err != nil {
		SendError(w, r, invalidIDErr, http.StatusNotFound)
		return
	}

	endpointData := endpointRest{}
	err = render.DecodeJSON(r.Body, &endpointData)
	if err != nil {
		SendInvalidJSONErr(w, r)
		return
	}

	endpointData.DeviceID = &deviceID
	endpoint, err := epHandler.endpointService.Create(toEndpoint(endpointData))
	if err != nil {
		serviceErr := utils.ToServiceErr(err)
		switch serviceErr.Code {
		case endpoints.InvalidInputCode:
			SendError(w, r, *serviceErr, http.StatusBadRequest)
		default:
			SendUnexpectedErr(w, r)
		}
		return
	}

	endpointData = fromEndpoint(*endpoint)

	result := &map[string]interface{}{
		"endpoint": endpointData,
	}

	SendResponse(w, r, result)
}

func (ep *EndpointHandler) Update(w http.ResponseWriter, r *http.Request) {
	endpointID, err := strconv.ParseInt(chi.URLParam(r, "endpointID"), 10, 64)
	if err != nil {
		SendError(w, r, invalidIDErr, http.StatusNotFound)
		return
	}

	endpointData := endpointRest{}
	err = render.DecodeJSON(r.Body, &endpointData)
	if err != nil {
		SendInvalidJSONErr(w, r)
		return
	}

	endpoint, err := ep.endpointService.Update(endpointID, toEndpoint(endpointData))
	if err != nil {
		serviceErr := utils.ToServiceErr(err)
		switch serviceErr.Code {
		case endpoints.InvalidInputCode:
			SendError(w, r, *serviceErr, http.StatusBadRequest)
		default:
			SendUnexpectedErr(w, r)
		}
		return
	}

	endpointData = fromEndpoint(*endpoint)

	result := &map[string]interface{}{
		"endpoint": endpointData,
	}

	SendResponse(w, r, result)
}

func (epHandler *EndpointHandler) Delete(w http.ResponseWriter, r *http.Request) {
	endpointID, err := strconv.ParseInt(chi.URLParam(r, "endpointID"), 10, 64)
	if err != nil {
		SendError(w, r, invalidIDErr, http.StatusNotFound)
		return
	}

	err = epHandler.endpointService.Delete(endpointID)
	if err != nil {
		serviceErr := utils.ToServiceErr(err)
		switch serviceErr.Code {
		case endpoints.DeviceNotFoundCode:
			SendError(w, r, *serviceErr, http.StatusNotFound)
		default:
			SendUnexpectedErr(w, r)
		}
		return
	}

	SendResponse(w, r, nil)
}

func (epHandler *EndpointHandler) GetbyDeviceID(w http.ResponseWriter, r *http.Request) {
	deviceID, err := strconv.ParseInt(chi.URLParam(r, "deviceID"), 10, 64)
	if err != nil {
		SendError(w, r, invalidIDErr, http.StatusNotFound)
		return
	}

	deviceEndpoints, err := epHandler.endpointService.GetbyDeviceID(deviceID)
	if err != nil {
		serviceErr := utils.ToServiceErr(err)
		switch serviceErr.Code {
		case endpoints.DeviceNotFoundCode:
			SendError(w, r, *serviceErr, http.StatusNotFound)
		default:
			SendUnexpectedErr(w, r)
		}
	}

	restEndpoints := make([]endpointRest, len(deviceEndpoints))

	for i := 0; i < len(deviceEndpoints); i++ {
		restEndpoints[i] = fromEndpoint(deviceEndpoints[i])
	}

	result := &map[string]interface{}{
		"endpoints": restEndpoints,
	}

	SendResponse(w, r, result)

}

//Endpoint Json Definition
type endpointRest struct {
	ID          *int64     `json:"id,omitempty"`
	DeviceID    *int64     `json:"device_id,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DisplayName *string    `json:"display_name,omitempty"`
	Description *string    `json:"description,omitempty"`
	Pattern     *string    `json:"pattern,omitempty"`
}

func toEndpoint(epRest endpointRest) endpoints.Endpoint {
	var ep endpoints.Endpoint

	if epRest.DisplayName != nil {
		ep.DisplayName = *epRest.DisplayName
	}

	if epRest.Description != nil {
		ep.Description = *epRest.Description
	}

	if epRest.DeviceID != nil {
		ep.DeviceID = *epRest.DeviceID
	}

	if epRest.Pattern != nil {
		ep.Pattern = *epRest.Pattern
	}

	return ep
}

func fromEndpoint(ep endpoints.Endpoint) endpointRest {
	var epRest endpointRest

	epRest.ID = &ep.ID
	epRest.DeviceID = &ep.DeviceID
	epRest.CreatedAt = &ep.CreatedAt
	epRest.DisplayName = &ep.DisplayName
	epRest.Description = &ep.Description
	epRest.Pattern = &ep.Pattern

	if !ep.UpdatedAt.IsZero() {
		epRest.UpdatedAt = &ep.UpdatedAt
	}

	return epRest
}
