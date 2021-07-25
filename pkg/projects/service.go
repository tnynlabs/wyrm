package projects

import (
	"log"
	"time"
	"github.com/tnynlabs/wyrm/pkg/utils"
)

type Project struct {
	ID        int64
	CreatedBy int64
	CreatedAt time.Time
	UpdatedAt time.Time

	DisplayName string
	Description string
}

type Repository interface {
	GetByID(projectID int64) (*Project, error)
	GetAllowed(userID int64) ([]Project, error)
	Create(p Project) (*Project, error)
	Update(projectID int64, p Project) (*Project, error)
	Delete(projectID int64) error
	AddCollaborator(userID int64, projectID int64) error
}

type Service interface {
	GetByID(projectID int64) (*Project, error)
	GetAllowed(userID int64) ([]Project, error)
	Create(p Project) (*Project, error)
	Update(projectID int64, p Project) (*Project, error)
	Delete(projectID int64) error
	AddCollaborator(userID int64, projectID int64) error
}

type service struct {
	projectRepo Repository
}

func CreateService(repo Repository) Service {
	return &service{repo}
}

func (s *service) GetByID(projectID int64) (*Project, error) {
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return nil, &utils.ServiceErr{
			Code:    ProjectNotFoundCode,
			Message: "Invalid ID",
		}
	}

	return project, nil
}

func (s *service) GetAllowed(userID int64) ([]Project, error) {
	projects, err := s.projectRepo.GetAllowed(userID)
	if err != nil {
		return nil, &utils.ServiceErr{
			Code:    UserNotFoundCode,
			Message: "User Not Found",
		}
	}

	return projects, nil
}

func (s *service) Create(p Project) (*Project, error) {
	if p.DisplayName == "" {
		return nil, &utils.ServiceErr{
			Code:    InvalidInputCode,
			Message: "Invalid display name",
		}
	}
	// TODO : Duplicate project name
	newProject, err := s.projectRepo.Create(p)
	if err != nil {
		log.Printf("Failed creating new project (error: %v)", err)
		return nil, &utils.ServiceErr{
			Code:    utils.UnexpectedCode,
			Message: "Failed creating new project",
		}
	}

	return newProject, nil
}
func (s *service) AddCollaborator(userID int64, projectID int64) error{
	err := s.projectRepo.AddCollaborator(userID, projectID)
	if err != nil {
		log.Printf("Failed adding collaborator (error: %v)", err)
		return &utils.ServiceErr{
			Code:    utils.UnexpectedCode,
			Message: "Failed adding collaborator",
		}
	}
	
	return nil
}
func (s *service) Update(projectID int64, p Project) (*Project, error) {
	if p.DisplayName == "" {
		return nil, &utils.ServiceErr{
			Code:    InvalidInputCode,
			Message: "Invalid display name",
		}
	}

	updatedData := Project{
		Description: p.Description,
		DisplayName: p.DisplayName,
	}
	project, err := s.projectRepo.Update(projectID, updatedData)
	if err != nil {
		return nil, &utils.ServiceErr{
			Code:    ProjectNotFoundCode,
			Message: "Invalid ID",
		}
	}

	return project, nil
}

func (s *service) Delete(projectID int64) error {
	err := s.projectRepo.Delete(projectID)
	if err != nil {
		return &utils.ServiceErr{
			Code:    ProjectNotFoundCode,
			Message: "Invalid ID",
		}
	}
	
	return nil
}
