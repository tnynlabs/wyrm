package rest

import (
	"net/http"
	"strconv"
	"time"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/tnynlabs/wyrm/pkg/pipelines"
	"github.com/tnynlabs/wyrm/pkg/projects"
	"github.com/tnynlabs/wyrm/pkg/utils"
)

type PipelineHandler struct {
	pipelineService pipelines.Service
	projectService  projects.Service
}

func CreatePipelineHandler(pipelineService pipelines.Service, projectService projects.Service) PipelineHandler {
	return PipelineHandler{pipelineService, projectService}
}

func (h *PipelineHandler) Get(w http.ResponseWriter, r *http.Request) {
	pipelineID, err := strconv.ParseInt(chi.URLParam(r, "pipelineID"), 10, 64)
	if err != nil {
		SendError(w, r, invalidIDErr, http.StatusNotFound)
		return
	}
	fmt.Println(pipelineID)
	pipeline, err := h.pipelineService.GetByID(pipelineID)
	if err != nil {
		serviceErr := utils.ToServiceErr(err)
		switch serviceErr.Code {
		case pipelines.PipelineNotFoundCode:
			SendError(w, r, *serviceErr, http.StatusNotFound)
		default:
			SendUnexpectedErr(w, r)
		}
		return
	}
	result := &map[string]interface{}{
		"pipeline": fromPipeline(*pipeline),
	}
	SendResponse(w, r, result)
}

func (h *PipelineHandler) Create(w http.ResponseWriter, r *http.Request) {
	projectID, err := strconv.ParseInt(chi.URLParam(r, "projectID"), 10, 64)
	if err != nil {
		SendError(w, r, invalidIDErr, http.StatusNotFound)
		return
	}

	project, err := h.projectService.GetByID(projectID)
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

	pipelineData := pipelineRest{}
	err = render.DecodeJSON(r.Body, &pipelineData)
	if err != nil {
		SendInvalidJSONErr(w, r)
		return
	}
	pipelineData.ProjectID = &project.ID
	pipeline, err := h.pipelineService.Create(*toPipeline(pipelineData))
	if err != nil {
		serviceErr := utils.ToServiceErr(err)
		switch serviceErr.Code {
		case pipelines.InvalidInputCode:
			SendError(w, r, *serviceErr, http.StatusBadRequest)
		default:
			SendUnexpectedErr(w, r)
		}
		return
	}
	result := &map[string]interface{}{
		"pipeline": fromPipeline(*pipeline),
	}

	SendResponse(w, r, result)
}

func (h *PipelineHandler) Update(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handler yoo")
	pipelineID, err := strconv.ParseInt(chi.URLParam(r, "pipelineID"), 10, 64)
	if err != nil {
		SendError(w, r, invalidIDErr, http.StatusNotFound)
		return
	}
	fmt.Println("Handler yoo 1")
	pipelineData := pipelineRest{}
	err = render.DecodeJSON(r.Body, &pipelineData)
	if err != nil {
		SendInvalidJSONErr(w, r)
		return
	}
	fmt.Println("Handler yoo 2")
	pipeline, err := h.pipelineService.Update(pipelineID, *toPipeline(pipelineData))
	if err != nil {
		serviceErr := utils.ToServiceErr(err)
		switch serviceErr.Code {
		case pipelines.InvalidInputCode:
			SendError(w, r, *serviceErr, http.StatusBadRequest)
		case pipelines.PipelineNotFoundCode:
			SendError(w, r, *serviceErr, http.StatusNotFound)
		default:
			SendUnexpectedErr(w, r)
		}
		return
	}
	
	result := &map[string]interface{}{
		"pipeline": fromPipeline(*pipeline),
	}

	SendResponse(w, r, result)
}

func (h *PipelineHandler) Delete(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handler yo")
	pipelineID,err := strconv.ParseInt(chi.URLParam(r, "pipelineID"), 10, 64)
	if err != nil {
		SendError(w, r, invalidIDErr, http.StatusNotFound)
		return
	}
	err = h.pipelineService.Delete(int64(pipelineID))
	if err != nil {
		serviceErr := utils.ToServiceErr(err)
		switch serviceErr.Code {
		case pipelines.PipelineNotFoundCode:
			SendError(w, r, *serviceErr, http.StatusNotFound)
		default:
			SendUnexpectedErr(w, r)
		}
		return
	}

	SendResponse(w, r, nil)
}

func (h *PipelineHandler) GetByProjectID(w http.ResponseWriter, r *http.Request) {
	projectID,err := strconv.ParseInt(chi.URLParam(r, "projectID"), 10, 64)
	if err != nil {
		SendError(w, r, invalidIDErr, http.StatusNotFound)
		return
	}
	projectPipelines, err := h.pipelineService.GetByProjectID(projectID)
	if err != nil {
		serviceErr := utils.ToServiceErr(err)
		switch serviceErr.Code {
		case pipelines.PipelineNotFoundCode:
			SendError(w, r, *serviceErr, http.StatusNotFound)
		default:
			SendUnexpectedErr(w, r)
		}
		return 
	}

	restPipelines := make([]*pipelineRest, len(projectPipelines))
	for i := 0; i < len(projectPipelines); i++ {
		restPipelines[i] = fromPipeline(projectPipelines[i])
	}

	result := &map[string]interface{}{
		"pipelines": restPipelines,
	}

	SendResponse(w, r, result)
}

type pipelineRest struct {
	ID          *int64     `json:"id,omitempty"`
	DisplayName *string    `json:"display_name,omitempty"`
	Data        *string    `json:"data,omitempty"`
	Description *string    `json:"description,omitempty"`
	ProjectID   *int64     `json:"project_id,omitempty"`
	CreatedBy   *int64     `json:"created_by,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

func toPipeline(pRest pipelineRest) *pipelines.Pipeline {
	var p pipelines.Pipeline

	if pRest.DisplayName != nil {
		p.DisplayName = *pRest.DisplayName
	}
	if pRest.Description != nil {
		p.Description = *pRest.Description
	}
	if pRest.ProjectID != nil {
		p.ProjectID = *pRest.ProjectID
	}
	if pRest.Data != nil {
		p.Data = *pRest.Data
	}
	if pRest.CreatedBy != nil {
		p.CreatedBy = *pRest.CreatedBy
	}

	return &p
}

func fromPipeline(p pipelines.Pipeline) *pipelineRest {
	pRest := pipelineRest{
		ID:          &p.ID,
		DisplayName: &p.DisplayName,
		Data:        &p.Data,
		Description: &p.Description,
		CreatedAt:   &p.CreatedAt,
		CreatedBy:   &p.CreatedBy,
	}
	if !p.UpdatedAt.IsZero() {
		pRest.UpdatedAt = &p.UpdatedAt
	} else {
		pRest.UpdatedAt = nil
	}

	return &pRest
}
