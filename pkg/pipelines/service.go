package pipelines

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/tnynlabs/wyrm/pkg/pipelines/protobuf"
	"github.com/tnynlabs/wyrm/pkg/utils"

	"google.golang.org/grpc"
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
	Create(p Pipeline) (*Pipeline, error)
	Update(pipelineID int64, pipeline Pipeline) (*Pipeline, error)
	Delete(pipelineID int64) error
}

type Service interface {
	GetByID(pipelineID int64) (*Pipeline, error)
	GetByProjectID(projectID int64) ([]Pipeline, error)
	Create(p Pipeline) (*Pipeline, error)
	Update(pipelineID int64, pipeline Pipeline) (*Pipeline, error)
	Delete(pipelineID int64) error
	RunPipeline(pipelineID int64, payload string) error
}

type service struct {
	pipelineRepo Repository
	client       protobuf.PipelineWorkerClient
}

func CreateService(repo Repository, workerAddr string) (Service, error) {
	log.Println((workerAddr))
	conn, err := grpc.Dial(workerAddr, grpc.WithInsecure())
	if err != nil {
		log.Println("GRPC Server Error")
		return nil, err
	}
	workerClient := protobuf.NewPipelineWorkerClient(conn)
	svc := service{
		pipelineRepo: repo,
		client:       workerClient,
	}
	return &svc, nil
}

func (s *service) GetByID(pipelineID int64) (*Pipeline, error) {
	pipeline, err := s.pipelineRepo.GetByID(pipelineID)
	if err != nil {
		return nil, &utils.ServiceErr{
			Code:    PipelineNotFoundCode,
			Message: "Invalid ID",
		}
	}

	return pipeline, nil
}

func (s *service) GetByProjectID(projectID int64) ([]Pipeline, error) {

	pipelines, err := s.pipelineRepo.GetByProjectID(projectID)
	if err != nil {
		return nil, &utils.ServiceErr{
			Code:    ProjectNotFoundCode,
			Message: "Invalid ID",
		}
	}

	return pipelines, nil
}

func (s *service) Create(p Pipeline) (*Pipeline, error) {
	if p.DisplayName == "" {
		return nil, &utils.ServiceErr{
			Code:    InvalidInputCode,
			Message: "Invalid display name",
		}
	}
	if p.Data == "" || p.Data == "{}" {
		return nil, &utils.ServiceErr{
			Code:    InvalidInputCode,
			Message: "Invalid pipeline structure",
		}
	}
	newPipeline, err := s.pipelineRepo.Create(p)
	if err != nil {
		log.Printf("Failed creating new pipeline (error: %v", err)
		return nil, &utils.ServiceErr{
			Code:    utils.UnexpectedCode,
			Message: "Failed creating new pipeline",
		}
	}
	return newPipeline, nil
}

func (s *service) Update(pipelineID int64, p Pipeline) (*Pipeline, error) {
	if p.DisplayName == "" {
		return nil, &utils.ServiceErr{
			Code:    InvalidInputCode,
			Message: "Invalid display name",
		}
	}
	if p.Data == "" || p.Data == "{}" {
		return nil, &utils.ServiceErr{
			Code:    InvalidInputCode,
			Message: "Invalid pipeline data",
		}
	}
	updatedData := Pipeline{
		DisplayName: p.DisplayName,
		Description: p.Description,
		Data:        p.Data,
	}
	pipeline, err := s.pipelineRepo.Update(pipelineID, updatedData)
	if err != nil {
		return nil, &utils.ServiceErr{
			Code:    PipelineNotFoundCode,
			Message: "Invalid ID",
		}
	}

	return pipeline, nil
}

func (s *service) Delete(pipelineID int64) error {
	err := s.pipelineRepo.Delete(pipelineID)
	if err != nil {
		return &utils.ServiceErr{
			Code:    PipelineNotFoundCode,
			Message: "Invalid ID",
		}
	}
	return nil
}

func (s *service) RunPipeline(pipelineID int64, payload string) error {
	pipelineRequest := protobuf.PipelineRequest{
		PipelineId: pipelineID,
		Payload:    payload,
	}

	_, err := s.client.RunPipeline(context.Background(), &pipelineRequest)
	if err != nil {
		errMsg := fmt.Sprintf("Pipeline run failed (%v)", err)
		return &utils.ServiceErr{
			Code:    WorkerConnectionErrorCode,
			Message: errMsg,
		}
	}

	return nil
}
