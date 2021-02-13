package rest

import (
	"fmt"
	"net/http"
	"time"

	"github.com/tnynlabs/wyrm/pkg/users"
	"github.com/tnynlabs/wyrm/pkg/utils"

	"github.com/go-chi/render"
)

// UserHandler user rest handler
type UserHandler struct {
	userService users.Service
}

func CreateUserHandler(userService users.Service) UserHandler {
	return UserHandler{userService}
}

func (h *UserHandler) RegisterWithPwd(w http.ResponseWriter, r *http.Request) {
	req := registerPwdRequest{}

	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		sendError(w, r, invalidJSONErr, http.StatusBadRequest)
		return
	}

	userData := users.User{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Email:       req.Email,
	}

	user, err := h.userService.CreateWithPwd(userData, req.Pwd)
	if err != nil {
		serviceErr := utils.ToServiceErr(err)
		switch serviceErr.Code {
		case users.InvalidInputCode, users.DuplicateEmailCode, users.DuplicateNameCode:
			sendError(w, r, *serviceErr, http.StatusBadRequest)
		default:
			sendError(w, r, unexpectedErr, http.StatusInternalServerError)
		}
		return
	}

	result := &map[string]interface{}{
		"user": fromUser(*user),
	}

	sendResponse(w, r, result)
}

func (h *UserHandler) LoginWithEmailPwd(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "welcome\n")
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "welcome\n")
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "welcome\n")
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "welcome\n")
}

type userRest struct {
	ID          *int64     `json:"id"`
	Name        *string    `json:"name"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	Email       *string    `json:"email"`
	DisplayName *string    `json:"display_name"`
}

// toUser maps userRest (json request) to users.User type
// Note: Only updatable fields are initialized
func toUser(uRest userRest) *users.User {
	var u users.User

	if uRest.Name != nil {
		u.Name = *uRest.Name
	}
	if uRest.Email != nil {
		u.Email = *uRest.Email
	}
	if uRest.DisplayName != nil {
		u.DisplayName = *uRest.DisplayName
	}

	return &u
}

// fromUser maps users.User to userRest (json response)
// Note: Only fields that might be empty are checked (i.e. UpdatedAt)
func fromUser(u users.User) *userRest {
	var uRest userRest

	uRest.ID = &u.ID
	uRest.Name = &u.Name
	uRest.CreatedAt = &u.CreatedAt
	uRest.Email = &u.Email
	uRest.DisplayName = &u.DisplayName
	if !u.UpdatedAt.IsZero() {
		uRest.UpdatedAt = &u.UpdatedAt
	}

	return &uRest
}

type registerPwdRequest struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Pwd         string `json:"password"`
}
