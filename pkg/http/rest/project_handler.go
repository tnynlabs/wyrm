package rest

import (
	"net/http"
	"strconv"
	"time"

	"github.com/tnynlabs/wyrm/pkg/projects"
	"github.com/tnynlabs/wyrm/pkg/users"
	"github.com/tnynlabs/wyrm/pkg/utils"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type ProjectHandler struct {
	projectService projects.Service
}

func CreateProjectHandler(projectService projects.Service) ProjectHandler {
	return ProjectHandler{projectService}
}
func (h *ProjectHandler) GetAllowed(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		SendError(w, r, invalidIDErr, http.StatusNotFound)
		return
	}
	projects, err := h.projectService.GetAllowedProjects(userID)
	if err != nil {
		SendUnexpectedErr(w, r)
		return
	}
	result := &map[string]interface{}{
		"projects": projects,
	}
	SendResponse(w, r, result)
}
func (h *ProjectHandler) Get(w http.ResponseWriter, r *http.Request) {
	projectID, err := strconv.ParseInt(chi.URLParam(r, "projectID"), 10, 64)
	if err != nil {
		SendError(w, r, invalidIDErr, http.StatusNotFound)
		return
	}

	project, err := h.projectService.GetByID(projectID)
	if err != nil {
		SendUnexpectedErr(w, r)
		return
	}
	result := &map[string]interface{}{
		"project": fromProject(*project),
	}
	SendResponse(w, r, result)
}
func (h *ProjectHandler) Update(w http.ResponseWriter, r *http.Request) {
	projectID, err := strconv.ParseInt(chi.URLParam(r, "projectID"), 10, 64)
	if err != nil {
		SendError(w, r, invalidIDErr, http.StatusNotFound)
		return
	}

	projectData := projectRest{}
	err = render.DecodeJSON(r.Body, &projectData)
	if err != nil {
		SendInvalidJSONErr(w, r)
		return
	}
	project, err := h.projectService.Update(projectID, *toProject(projectData))
	if err != nil {
		serviceErr := utils.ToServiceErr(err)
		switch serviceErr.Code {
		case projects.InvalidInputCode:
			SendError(w, r, *serviceErr, http.StatusBadRequest)
		case projects.ProjectNotFoundCode:
			SendError(w, r, *serviceErr, http.StatusNotFound)
		default:
			SendUnexpectedErr(w, r)
		}
		return
	}
	result := &map[string]interface{}{
		"project": fromProject(*project),
	}

	SendResponse(w, r, result)
}
func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	projectID, err := strconv.ParseInt(chi.URLParam(r, "projectID"), 10, 64)
	if err != nil {
		SendError(w, r, invalidIDErr, http.StatusNotFound)
		return
	}
	err = h.projectService.Delete(int64(projectID))
	if err != nil {
		serviceErr := utils.ToServiceErr(err)
		switch serviceErr.Code {
		case projects.ProjectNotFoundCode:
			SendError(w, r, *serviceErr, http.StatusNotFound)
		default:
			SendUnexpectedErr(w, r)
		}
		return
	}

	SendResponse(w, r, nil)
}

type createProjectRequest struct {
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
}

func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	req := createProjectRequest{}
	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		SendInvalidJSONErr(w, r)
		return
	}

	var user = r.Context().Value(UserCtxKey{}).(*users.User)
	projectData := projects.Project{
		DisplayName: req.DisplayName,
		Description: req.Description,
		CreatedBy:   user.ID,
	}

	project, err := h.projectService.Create(projectData)
	if err != nil {
		serviceErr := utils.ToServiceErr(err)
		switch serviceErr.Code {
		// TODO add duplicate
		case projects.InvalidInputCode:
			SendError(w, r, *serviceErr, http.StatusBadRequest)
		default:
			SendUnexpectedErr(w, r)
		}
		return
	}
	result := &map[string]interface{}{
		"project": fromProject(*project),
	}
	SendResponse(w, r, result)
}

type projectRest struct {
	ID          *int64     `json:"id,omitempty"`
	CreatedBy   *int64     `json:"created_by,omitempty"`
	Description *string    `json:"description,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DisplayName *string    `json:"display_name,omitempty"`
}

func toProject(pRest projectRest) *projects.Project {
	var p projects.Project

	if pRest.Description != nil {
		p.Description = *pRest.Description
	}
	if pRest.DisplayName != nil {
		p.DisplayName = *pRest.DisplayName
	}
	return &p
}

func fromProject(p projects.Project) *projectRest {
	pRest := projectRest{
		ID:          &p.ID,
		CreatedBy:   &p.CreatedBy,
		CreatedAt:   &p.CreatedAt,
		DisplayName: &p.DisplayName,
		Description: &p.Description,
	}
	if !p.UpdatedAt.IsZero() {
		pRest.UpdatedAt = &p.UpdatedAt
	} else {
		pRest.UpdatedAt = nil
	}
	return &pRest
}
