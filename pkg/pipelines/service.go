package pipelines

import (
	"log"
	"time"
	"github.com/tnynlabs/wyrm/pkg/utils"
)

type Pipeline struct {
	ID          int64
	DisplayName string
	Data        string
	Description string
	ProjectID   int64
	CreatedBy   int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Repository interface {
	GetByID(pipelineID int64) (*Pipeline, error)
	GetByProjectID(projectID int64) ([]Pipeline, error)
	Create(pipeline Pipeline) (*Pipeline, error)
	Update(pipelineID int64, pipeline Pipeline) (*Pipeline, error)
	Delete(pipelineID int64) error
}
type Service interface {
	GetByID(pipelineID int64) (*Pipeline, error)
	GetByProjectID(projectID int64) ([]Pipeline, error)
	Create(pipeline Pipeline) (*Pipeline, error)
	Update(pipelineID int64, pipeline Pipeline) (*Pipeline, error)
	Delete(pipelineID int64) error
}
type service struct {
	pipelineRepo Repository
}

func CreateService(repo Repository) Service {
	return &service{repo}
}
func (s *service) GetByID(pipelineID int64) (*Pipeline, error) {
	return nil, nil
}
func (s *service) GetByProjectID(projectID int64) ([]Pipeline, error) {
	return nil, nil
}
func (s *service) Create(pipeline Pipeline) (*Pipeline, error) {
	if pipeline.DisplayName == "" {
		return nil, &utils.ServiceErr{
			Code:    InvalidInputCode,
			Message: "Invalid display name",
		}
	}
	if pipeline.Data == "" {
		return nil, &utils.ServiceErr{
			Code:    InvalidInputCode,
			Message: "Invalid pipeline structure",
		}
	}
	newPipeline, err := s.pipelineRepo.Create(pipeline)
	if err != nil {
		log.Printf("Failed creating new pipeline (error: %v", err)
		return nil, &utils.ServiceErr{
			Code:    utils.UnexpectedCode,
			Message: "Failed creating new pipeline",
		}
	}
	return newPipeline, nil
}
func (s *service) Update(pipelineID int64, pipeline Pipeline) (*Pipeline, error) {
	return nil, nil
}
func (s *service) Delete(pipelineID int64) error {
	return nil
}

