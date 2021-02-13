package rest

import (
	"net/http"

	"github.com/tnynlabs/wyrm/pkg/utils"

	"github.com/go-chi/render"
)

type restErr struct {
	Code    utils.ServiceErrCode `json:"code"`
	Message string               `json:"message"`
}

type response struct {
	Result *map[string]interface{} `json:"result"`
	Err    *restErr                `json:"error"`
}

func sendResponse(w http.ResponseWriter, r *http.Request, result *map[string]interface{}) {
	resp := response{
		Result: result,
		Err:    nil,
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, resp)
}

func sendError(w http.ResponseWriter, r *http.Request, err utils.ServiceErr, status int) {
	resp := response{
		Result: nil,
		Err: &restErr{
			Code:    err.Code,
			Message: err.Message,
		},
	}
	render.Status(r, status)
	render.JSON(w, r, resp)
}

var unexpectedErr = utils.ServiceErr{
	Code:    utils.UnexpectedCode,
	Message: "An unexpected error occurred",
}

var invalidJSONErr = utils.ServiceErr{
	Code:    "INVALID_JSON",
	Message: "Invalid json",
}
