package devices

import (
	"log"
	"time"

	"github.com/tnynlabs/wyrm/pkg/utils"
)

// Device Contains device core properties
// Note: zero values will not be updated
type Device struct {
	//Basic attributes
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time

	// Note: Only these values are updatable
	Description string
	DisplayName string
	ProjectID   int64

	// Note: Never show in output
	AuthKey string
}

//Defines devices.Repository for Storage Implementation
type Repository interface {
	GetByID(deviceID int64) (*Device, error)
	Create(d Device) (*Device, error)
	Update(deviceID int64, d Device) (*Device, error)
	Delete(deviceID int64) error
	GetByProjectID(projectID int64) ([]Device, error)
}

type Service interface {
	GetByID(deviceID int64) (*Device, error)
	Create(d Device) (*Device, error)
	Update(deviceID int64, d Device) (*Device, error)
	Delete(deviceID int64) error
	GetByProjectID(projectID int64) ([]Device, error)
}

type service struct {
	deviceRepo Repository
}

func CreateDeviceService(deviceRepo Repository) Service {
	return &service{deviceRepo}
}

func (s *service) GetByID(deviceID int64) (*Device, error) {
	device, err := s.deviceRepo.GetByID(deviceID)
	if err != nil {
		return nil, &utils.ServiceErr{
			Code:    DeviceNotFoundCode,
			Message: "Invalid Device id",
		}
	}

	return device, nil
}

func (s *service) Create(d Device) (*Device, error) {
	device, err := s.deviceRepo.Create(d)
	if err != nil {
		return nil, &utils.ServiceErr{
			Code:    InvalidInputCode,
			Message: "Invalid input",
		}
	}

	return device, nil
}

func (s *service) Update(deviceID int64, d Device) (*Device, error) {
	device, err := s.deviceRepo.Update(deviceID, d)
	if err != nil {
		log.Println(d)
		log.Println(err)
		return nil, &utils.ServiceErr{
			Code:    InvalidInputCode,
			Message: "Invalid input",
		}
	}

	return device, nil
}

func (s *service) Delete(deviceID int64) error {
	err := s.deviceRepo.Delete(deviceID)
	if err != nil {
		return &utils.ServiceErr{
			Code:    DeviceNotFoundCode,
			Message: "Invalid Device ID",
		}
	}

	return nil
}

func (s *service) GetByProjectID(projectID int64) ([]Device, error) {
	devices, err := s.deviceRepo.GetByProjectID(projectID)
	if err != nil {
		return nil, &utils.ServiceErr{
			Code:    ProjectNotFoundCode,
			Message: "Invalid Project ID",
		}
	}

	return devices, nil
}
