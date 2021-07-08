package rest

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/tnynlabs/wyrm/pkg/devices"
	"github.com/tnynlabs/wyrm/pkg/utils"
)

type DeviceHandler struct {
	deviceService devices.Service
}

func CreateDeviceHandler(dService devices.Service) DeviceHandler {
	return DeviceHandler{dService}
}

func (dHandler *DeviceHandler) Get(w http.ResponseWriter, r *http.Request) {
	deviceID, err := strconv.ParseInt(chi.URLParam(r, "deviceID"), 10, 64)
	if err != nil {
		SendError(w, r, invalidIDErr, http.StatusNotFound)
		return
	}

	device, err := dHandler.deviceService.GetByID(deviceID)
	if err != nil {
		serviceErr := utils.ToServiceErr(err)
		switch serviceErr.Code {
		case devices.DeviceNotFoundCode:
			SendError(w, r, *serviceErr, http.StatusNotFound)
		default:
			SendUnexpectedErr(w, r)
		}
		return
	}

	deviceData := fromDevice(*device)

	result := &map[string]interface{}{
		"device": deviceData,
	}

	SendResponse(w, r, result)
}

func (dHandler *DeviceHandler) Create(w http.ResponseWriter, r *http.Request) {
	projectID, err := strconv.ParseInt(chi.URLParam(r, "projectID"), 10, 64)
	if err != nil {
		SendError(w, r, invalidIDErr, http.StatusNotFound)
		return
	}

	deviceData := deviceRest{}
	err = render.DecodeJSON(r.Body, &deviceData)
	if err != nil {
		SendInvalidJSONErr(w, r)
		return
	}

	deviceData.ProjectID = &projectID
	device, err := dHandler.deviceService.Create(toDevice(deviceData))
	if err != nil {
		serviceErr := utils.ToServiceErr(err)
		switch serviceErr.Code {
		case devices.InvalidInputCode:
			SendError(w, r, *serviceErr, http.StatusBadRequest)
		default:
			SendUnexpectedErr(w, r)
		}
		return
	}

	deviceData = fromDevice(*device)

	result := &map[string]interface{}{
		"device": deviceData,
	}

	SendResponse(w, r, result)
}

func (dHandler *DeviceHandler) Update(w http.ResponseWriter, r *http.Request) {
	deviceID, err := strconv.ParseInt(chi.URLParam(r, "deviceID"), 10, 64)
	if err != nil {
		SendError(w, r, invalidIDErr, http.StatusNotFound)
		return
	}

	deviceData := deviceRest{}
	err = render.DecodeJSON(r.Body, &deviceData)
	if err != nil {
		SendInvalidJSONErr(w, r)
		return
	}

	//Only updatable fields are set in device object
	device, err := dHandler.deviceService.Update(deviceID, toDevice(deviceData))
	if err != nil {
		serviceErr := utils.ToServiceErr(err)
		switch serviceErr.Code {
		case devices.InvalidInputCode:
			SendError(w, r, *serviceErr, http.StatusBadRequest)
		default:
			SendUnexpectedErr(w, r)
		}
		return
	}

	deviceData = fromDevice(*device)

	result := map[string]interface{}{
		"device": deviceData,
	}

	SendResponse(w, r, &result)
}

func (dHandler *DeviceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	deviceID, err := strconv.ParseInt(chi.URLParam(r, "deviceID"), 10, 64)
	if err != nil {
		SendError(w, r, invalidIDErr, http.StatusNotFound)
		return
	}

	err = dHandler.deviceService.Delete(deviceID)
	if err != nil {
		serviceErr := utils.ToServiceErr(err)
		switch serviceErr.Code {
		case devices.DeviceNotFoundCode:
			SendError(w, r, *serviceErr, http.StatusNotFound)
		default:
			SendUnexpectedErr(w, r)
		}
		return
	}

	SendResponse(w, r, nil)
}

func (dHandler *DeviceHandler) GetByProjectID(w http.ResponseWriter, r *http.Request) {
	projectID, err := strconv.ParseInt(chi.URLParam(r, "projectID"), 10, 64)
	if err != nil {
		SendError(w, r, invalidIDErr, http.StatusNotFound)
		return
	}

	projectDevices, err := dHandler.deviceService.GetByProjectID(projectID)
	if err != nil {
		serviceErr := utils.ToServiceErr(err)
		switch serviceErr.Code {
		case devices.ProjectNotFoundCode:
			SendError(w, r, *serviceErr, http.StatusNotFound)
		default:
			SendUnexpectedErr(w, r)
		}
	}

	restDevices := make([]deviceRest, len(projectDevices))

	for i := 0; i < len(projectDevices); i++ {
		restDevices[i] = fromDevice(projectDevices[i])
	}

	result := &map[string]interface{}{
		"devices": restDevices,
	}

	SendResponse(w, r, result)
	return
}

//Device Json Definition
type deviceRest struct {
	ID          *int64     `json:"id,omitempty"`
	ProjectID   *int64     `json:"project_id,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DisplayName *string    `json:"display_name,omitempty"`
	AuthKey     *string    `json:"auth_key,omitempty"`
	Description *string    `json:"description,omitempty"`
}

func toDevice(dRest deviceRest) devices.Device {
	var d devices.Device

	if dRest.DisplayName != nil {
		d.DisplayName = *dRest.DisplayName
	}

	if dRest.Description != nil {
		d.Description = *dRest.Description
	}

	if dRest.ProjectID != nil {
		d.ProjectID = *dRest.ProjectID
	}

	return d
}

func fromDevice(d devices.Device) deviceRest {
	var dRest deviceRest

	dRest.ID = &d.ID
	dRest.ProjectID = &d.ProjectID
	dRest.CreatedAt = &d.CreatedAt
	dRest.DisplayName = &d.DisplayName
	dRest.AuthKey = &d.AuthKey
	dRest.Description = &d.Description

	if !d.UpdatedAt.IsZero() {
		dRest.UpdatedAt = &d.UpdatedAt
	}

	return dRest
}
