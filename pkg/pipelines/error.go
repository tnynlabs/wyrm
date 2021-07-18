package pipelines

import "github.com/tnynlabs/wyrm/pkg/utils"

const (
	PipelineNotFoundCode  = utils.ServiceErrCode("DEVICE_NOT_FOUND")
	ProjectNotFoundCode = utils.ServiceErrCode("PROJECT_NOT_FOUND")
	InvalidInputCode    = utils.ServiceErrCode("INVALID_INPUT")
)