package postgres

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/tnynlabs/wyrm/pkg/endpoints"
)

type EndpointRepository struct {
	db *sqlx.DB
}

func CreateEndpointRepository(db *sqlx.DB) endpoints.Repository {
	return &EndpointRepository{db}
}

func (epR *EndpointRepository) GetByID(endpointID int64) (*endpoints.Endpoint, error) {
	const sqlStmt = `
	Select id, device_id, display_name, description, pattern
	From endpoints
	where id = $1 `
	var endpointData endpointSQL
	err := epR.db.Get(&endpointData, sqlStmt, endpointID)
	if err != nil {
		return nil, err
	}

	return toEndpoint(endpointData), nil
}

func (epR *EndpointRepository) Create(ep endpoints.Endpoint) (*endpoints.Endpoint, error) {
	ep.CreatedAt = time.Now()

	endpointData := fromEndpoint(ep)
	const sqlStmt = `
	INSERT INTO endpoints (
		device_id, display_name, description, pattern, created_at
	) VALUES (
		:device_id, :display_name, :description, :pattern, :created_at
	) RETURNING id`
	query, args, err := sqlx.Named(sqlStmt, endpointData)
	if err != nil {
		return nil, err
	}

	// https://pkg.go.dev/github.com/jmoiron/sqlx#Rebind
	// Replace ? with $ for postgres
	query = sqlx.Rebind(sqlx.DOLLAR, query)

	err = epR.db.Get(&ep.ID, query, args...)
	if err != nil {
		return nil, err
	}

	return &ep, nil
}

func (epR *EndpointRepository) Update(endpointID int64, ep endpoints.Endpoint) (*endpoints.Endpoint, error) {
	ep.ID = endpointID
	ep.UpdatedAt = time.Now()

	endpointData := fromEndpoint(ep)

	const sqlStmt = `
		UPDATE endpoints
		SET
			device_id = COALESCE(:device_id, device_id),
			display_name = COALESCE(:display_name, display_name),
			description = COALESCE(:description, description),
			updated_at = COALESCE(:updated_at, updated_at)
		WHERE id = :id`

	result, err := epR.db.NamedExec(sqlStmt, endpointData)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected != 1 {
		return nil, errors.New("Invalid ID")
	}

	user, err := epR.GetByID(endpointID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (epR *EndpointRepository) Delete(endpointID int64) error {
	const sqlStmt = `
		DELETE FROM endpoint
		WHERE id = $1
	`
	result, err := epR.db.Exec(sqlStmt, endpointID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected != 1 {
		return errors.New("Invalid ID")
	}

	return nil
}

func (epR *EndpointRepository) GetbyDeviceID(deviceID int64) ([]endpoints.Endpoint, error) {
	endpointsSQL := []endpointSQL{}
	const sqlStmt = `
	SELECT id, device_id, display_name, description, pattern
	FROM endpoints
	WHERE device_id = $1
	`
	err := epR.db.Select(endpointsSQL, sqlStmt, deviceID)
	if err != nil {
		return nil, err
	}

	endpoints := make([]endpoints.Endpoint, len(endpointsSQL))
	for i := 0; i < len(endpointsSQL); i++ {
		endpoints[i] = *toEndpoint(endpointsSQL[i])
	}

	return endpoints, nil
}

//stucrt to match the postgres database construction
type endpointSQL struct {
	ID          int64          `db:"id"`
	DeviceID    sql.NullInt64  `db:"device_id"`
	DisplayName sql.NullString `db:"display_name"`
	Description sql.NullString `db:"description"`
	Pattern     sql.NullString `db:"pattern"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   sql.NullTime   `db:"updated_at"`
}

func toEndpoint(epSQL endpointSQL) *endpoints.Endpoint {
	return &endpoints.Endpoint{
		ID:        epSQL.ID,
		CreatedAt: epSQL.CreatedAt,
		UpdatedAt: epSQL.UpdatedAt.Time,

		Description: epSQL.Description.String,
		DisplayName: epSQL.DisplayName.String,
		DeviceID:    epSQL.DeviceID.Int64,
		Pattern:     epSQL.Pattern.String,
	}
}

func fromEndpoint(ep endpoints.Endpoint) *endpointSQL {
	var endpointData endpointSQL
	endpointData.ID = ep.ID
	endpointData.CreatedAt = ep.CreatedAt
	endpointData.UpdatedAt = sql.NullTime{
		Time:  ep.UpdatedAt,
		Valid: !ep.UpdatedAt.IsZero(),
	}
	endpointData.Description = sql.NullString{
		String: ep.Description,
		Valid:  ep.Description != "",
	}
	endpointData.DeviceID = sql.NullInt64{
		Int64: ep.DeviceID,
		Valid: ep.DeviceID != 0,
	}
	endpointData.DisplayName = sql.NullString{
		String: ep.DisplayName,
		Valid:  ep.DisplayName != "",
	}
	endpointData.Pattern = sql.NullString{
		String: ep.Pattern,
		Valid:  ep.Pattern != "",
	}

	return &endpointData
}
