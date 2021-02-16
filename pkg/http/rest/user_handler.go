package rest

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/tnynlabs/wyrm/pkg/users"
	"github.com/tnynlabs/wyrm/pkg/utils"

	"github.com/go-chi/chi"
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
		SendInvalidJSONErr(w, r)
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
			SendError(w, r, *serviceErr, http.StatusBadRequest)
		default:
			SendUnexpectedErr(w, r)
		}
		return
	}

	result := &map[string]interface{}{
		"user": fromUser(*user),
	}

	SendResponse(w, r, result)
}

func (h *UserHandler) LoginWithEmailPwd(w http.ResponseWriter, r *http.Request) {
	req := loginEmailPwdRequest{}

	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		SendInvalidJSONErr(w, r)
		return
	}

	user, err := h.userService.AuthWithEmailPwd(req.Email, req.Pwd)
	if err != nil {
		serviceErr := utils.ToServiceErr(err)
		switch serviceErr.Code {
		case users.InvalidInputCode:
			SendError(w, r, *serviceErr, http.StatusUnauthorized)
		default:
			SendUnexpectedErr(w, r)
		}
		return
	}

	// TODO: add SecureOnly
	authKeyCookie := http.Cookie{
		Name:     "auth_key",
		Value:    user.AuthKey,
		HttpOnly: true,
	}

	http.SetCookie(w, &authKeyCookie)

	result := &map[string]interface{}{
		"user": fromUser(*user),
	}

	SendResponse(w, r, result)
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		SendError(w, r, invalidIDErr, http.StatusNotFound)
		return
	}

	user, err := h.userService.GetByID(int64(userID))
	if err != nil {
		serviceErr := utils.ToServiceErr(err)
		switch serviceErr.Code {
		case users.UserNotFoundCode:
			SendError(w, r, *serviceErr, http.StatusNotFound)
		default:
			SendUnexpectedErr(w, r)
		}
		return
	}

	userData := fromUser(*user)
	// omit sensitive values
	userData.CreatedAt = nil
	userData.UpdatedAt = nil
	userData.Email = nil

	result := &map[string]interface{}{
		"user": userData,
	}

	SendResponse(w, r, result)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "welcome\n")
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		SendError(w, r, invalidIDErr, http.StatusNotFound)
		return
	}

	err = h.userService.Delete(int64(userID))
	if err != nil {
		serviceErr := utils.ToServiceErr(err)
		switch serviceErr.Code {
		case users.UserNotFoundCode:
			SendError(w, r, *serviceErr, http.StatusNotFound)
		default:
			SendUnexpectedErr(w, r)
		}
		return
	}

	SendResponse(w, r, nil)
}

type userRest struct {
	ID          *int64     `json:"id,omitempty"`
	Name        *string    `json:"name,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	Email       *string    `json:"email,omitempty"`
	DisplayName *string    `json:"display_name,omitempty"`
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

type loginEmailPwdRequest struct {
	Email string `json:"email"`
	Pwd   string `json:"password"`
}

var invalidIDErr = utils.ServiceErr{
	Code:    users.UserNotFoundCode,
	Message: "Invalid ID (IDs should be integers)",
}
