package projects

import "github.com/tnynlabs/wyrm/pkg/utils"

const (
	InvalidInputCode   = utils.ServiceErrCode("INVALID_INPUT")
	DuplicateNameCode  = utils.ServiceErrCode("DUPLICATE_NAME")
	ProjectNotFoundCode   = utils.ServiceErrCode("PROJECT_NOT_FOUND")
)
