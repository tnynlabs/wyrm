package devices

import "time"

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
}

type Service interface {
	GetByID(deviceID int64) (*Device, error)
}

type service struct {
	deviceRepo Repository
}

func CreateDeviceService(deviceRepo Repository) Service {
	return &service{deviceRepo}
}

func (s *service) GetByID(deviceID int64) (*Device, error) {

}
