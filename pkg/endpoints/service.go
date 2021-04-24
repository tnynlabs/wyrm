package endpoints

import (
	"time"

	"github.com/tnynlabs/wyrm/pkg/utils"
)

//Struct contains endpoint main attributes

type Endpoint struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time

	// Note: Only these values are updatable
	DeviceID    int64
	Description string
	DisplayName string
	Pattern     string
}

type Repository interface {
	Create(ep Endpoint) (*Endpoint, error)
	GetByID(endpointID int64) (*Endpoint, error)
	Delete(endpointID int64) error
	Update(endpointID int64, ep Endpoint) (*Endpoint, error)
	GetbyDeviceID(deviceID int64) ([]Endpoint, error)
}

type Service interface {
	Create(ep Endpoint) (*Endpoint, error)
	GetByID(endpointID int64) (*Endpoint, error)
	Delete(endpointID int64) error
	Update(endpointID int64, ep Endpoint) (*Endpoint, error)
	GetbyDeviceID(deviceID int64) ([]Endpoint, error)
}

type service struct {
	endpointRepo Repository
}

func CreateEndpointService(endpointRepo Repository) Service {
	return &service{endpointRepo}
}

func (s *service) Create(ep Endpoint) (*Endpoint, error) {
	endpoint, err := s.endpointRepo.Create(ep)
	if err != nil {
		return nil, &utils.ServiceErr{
			Code:    InvalidInputCode,
			Message: "Invalid input",
		}
	}

	return endpoint, nil
}

func (s *service) GetByID(endpointID int64) (*Endpoint, error) {
	endpoint, err := s.endpointRepo.GetByID(endpointID)
	if err != nil {
		return nil, &utils.ServiceErr{
			Code:    EndpointNotFoundCode,
			Message: "Invalid ID",
		}
	}

	return endpoint, nil
}

func (s *service) Delete(endpointID int64) error {
	err := s.endpointRepo.Delete(endpointID)
	if err != nil {
		return &utils.ServiceErr{
			Code:    EndpointNotFoundCode,
			Message: "Invalid ID",
		}
	}

	return nil
}

func (s *service) Update(endpointID int64, ep Endpoint) (*Endpoint, error) {
	endpoint, err := s.endpointRepo.Update(endpointID, ep)
	if err != nil {
		return nil, &utils.ServiceErr{
			Code:    EndpointNotFoundCode,
			Message: "Invalid ID",
		}
	}

	return endpoint, nil
}

func (s *service) GetbyDeviceID(deviceID int64) ([]Endpoint, error) {
	endpoints, err := s.endpointRepo.GetbyDeviceID(deviceID)
	if err != nil {
		return nil, &utils.ServiceErr{
			Code:    EndpointNotFoundCode,
			Message: "Invalid ID",
		}
	}

	return endpoints, nil
}
