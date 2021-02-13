package users

import "github.com/tnynlabs/wyrm/pkg/utils"

const (
	InvalidInputCode   = utils.ServiceErrCode("INVALID_INPUT")
	DuplicateEmailCode = utils.ServiceErrCode("DUPLICATE_EMAIL")
	DuplicateNameCode  = utils.ServiceErrCode("DUPLICATE_NAME")
)
