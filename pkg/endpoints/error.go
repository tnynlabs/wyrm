package endpoints

import "github.com/tnynlabs/wyrm/pkg/utils"

const (
	DeviceNotFoundCode   = utils.ServiceErrCode("DEVICE_NOT_FOUND")
	EndpointNotFoundCode = utils.ServiceErrCode("ENDPOINT_NOT_FOUND")
	InvalidInputCode     = utils.ServiceErrCode("INVALID_INPUT")
)
