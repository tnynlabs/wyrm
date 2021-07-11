package postgres

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/tnynlabs/wyrm/pkg/devices"
)

type DeviceRepository struct {
	db *sqlx.DB
}

func (dR *DeviceRepository) GetByID(deviceID int64) (*devices.Device, error) {
	const sqlStmt = `
	SELECT id, project_id, display_name, auth_key, description, created_at
	FROM Devices
	WHERE id = $1 `
	var deviceData deviceSQL
	err := dR.db.Get(&deviceData, sqlStmt, deviceID)
	if err != nil {
		return nil, err
	}

	return toDevice(deviceData), nil
}

func (dR *DeviceRepository) GetByKey(authKey string) (*devices.Device, error) {
	const sqlStmt = `
	Select id, project_id, display_name, auth_key, description, created_at
	FROM Devices
	WHERE auth_key = $1 `
	var deviceData deviceSQL
	err := dR.db.Get(&deviceData, sqlStmt, authKey)
	if err != nil {
		return nil, err
	}

	return toDevice(deviceData), nil
}

func (dR *DeviceRepository) Create(d devices.Device) (*devices.Device, error) {
	d.CreatedAt = time.Now()

	deviceData := fromDevice(d)
	const sqlStmt = `
	INSERT INTO devices (
		project_id, display_name, auth_key, description, created_at
	) VALUES (
		:project_id, :display_name, :auth_key, :description, :created_at
	) RETURNING id`

	query, args, err := sqlx.Named(sqlStmt, deviceData)
	if err != nil {
		return nil, err
	}

	// https://pkg.go.dev/github.com/jmoiron/sqlx#Rebind
	// Replace ? with $ for postgres
	query = sqlx.Rebind(sqlx.DOLLAR, query)

	err = dR.db.Get(&d.ID, query, args...)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func (dR DeviceRepository) Update(deviceID int64, d devices.Device) (*devices.Device, error) {
	d.ID = deviceID
	d.UpdatedAt = time.Now()

	deviceData := fromDevice(d)

	const sqlStmt = `
		UPDATE devices
		SET
			project_id = COALESCE(:project_id, project_id),
			display_name = COALESCE(:display_name, display_name),
			description = COALESCE(:description, description),
			updated_at = COALESCE(:updated_at, updated_at)
		WHERE id = :id
	`
	result, err := dR.db.NamedExec(sqlStmt, deviceData)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected != 1 {
		return nil, errors.New("Invalid ID")
	}

	user, err := dR.GetByID(deviceID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (dR *DeviceRepository) Delete(deviceID int64) error {
	const sqlStmt = `
		DELETE FROM devices
		WHERE id = $1
	`
	result, err := dR.db.Exec(sqlStmt, deviceID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected != 1 {
		return errors.New("Invalid ID")
	}

	return nil
}

func (dR *DeviceRepository) GetByProjectID(projectID int64) ([]devices.Device, error) {
	devicesSQL := []deviceSQL{}
	const sqlStmt = `
		SELECT id, project_id, display_name, auth_key, description, created_at 
		FROM devices
		WHERE project_id = $1
	`
	err := dR.db.Select(&devicesSQL, sqlStmt, projectID)
	if err != nil {
		return nil, err
	}

	devices := make([]devices.Device, len(devicesSQL))
	for i := 0; i < len(devicesSQL); i++ {
		devices[i] = *toDevice(devicesSQL[i])
	}

	return devices, nil
}

func CreateDeviceRepository(db *sqlx.DB) devices.Repository {
	return &DeviceRepository{db}
}

//SQL skeleton struct for devices
type deviceSQL struct {
	ID          int64          `db:"id"`
	ProjectID   sql.NullInt64  `db:"project_id"`
	DisplayName sql.NullString `db:"display_name"`
	AuthKey     sql.NullString `db:"auth_key"`
	Description sql.NullString `db:"description"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   sql.NullTime   `db:"updated_at"`
}

//Changing from postgress device implementation to device service implementation
func toDevice(dSQL deviceSQL) *devices.Device {
	return &devices.Device{
		ID:        dSQL.ID,
		CreatedAt: dSQL.CreatedAt,
		UpdatedAt: dSQL.UpdatedAt.Time,

		Description: dSQL.Description.String,
		DisplayName: dSQL.DisplayName.String,
		ProjectID:   dSQL.ProjectID.Int64,

		AuthKey: dSQL.AuthKey.String,
	}
}

//Changing from device service implementation to postgress device implementation
func fromDevice(d devices.Device) *deviceSQL {
	var deviceData deviceSQL
	deviceData.ID = d.ID
	deviceData.CreatedAt = d.CreatedAt
	deviceData.UpdatedAt = sql.NullTime{
		Time:  d.UpdatedAt,
		Valid: !d.UpdatedAt.IsZero(),
	}
	deviceData.Description = sql.NullString{
		String: d.Description,
		Valid:  d.Description != "",
	}
	deviceData.AuthKey = sql.NullString{
		String: d.AuthKey,
		Valid:  d.AuthKey != "",
	}
	deviceData.ProjectID = sql.NullInt64{
		Int64: d.ProjectID,
		Valid: d.ProjectID != 0,
	}
	deviceData.DisplayName = sql.NullString{
		String: d.DisplayName,
		Valid:  d.DisplayName != "",
	}

	return &deviceData
}
