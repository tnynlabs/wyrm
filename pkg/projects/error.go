package projects

import "github.com/tnynlabs/wyrm/pkg/utils"

const (
	InvalidInputCode   = utils.ServiceErrCode("INVALID_INPUT")
	ProjectNotFoundCode   = utils.ServiceErrCode("PROJECT_NOT_FOUND")
	UserNotFoundCode = utils.ServiceErrCode("USER_NOT_FOUND")
)
