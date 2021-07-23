package pipelines

import "github.com/tnynlabs/wyrm/pkg/utils"

const (
	WorkerConnectionErrorCode = utils.ServiceErrCode("WORKER_CONN_ERROR")
	PipelineNotFoundCode      = utils.ServiceErrCode("PIPELINE_NOT_FOUND")
	ProjectNotFoundCode       = utils.ServiceErrCode("PROJECT_NOT_FOUND")
	InvalidInputCode          = utils.ServiceErrCode("INVALID_INPUT")
)
