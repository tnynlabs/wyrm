package postgres

import (
	"database/sql"
	"time"
	"github.com/jmoiron/sqlx"
	"github.com/tnynlabs/wyrm/pkg/pipelines"
	"errors"
)

type PipelineRepository struct {
	db *sqlx.DB
}

func CreatePipelineRepository(db *sqlx.DB) pipelines.Repository {
	return &PipelineRepository{db}
}

func (pR *PipelineRepository) GetByID(pipelineID int64) (*pipelines.Pipeline, error) {
	const getByIDStmt = `
	SELECT
		id, project_id , display_name, data, description,created_at, updated_at, created_by
	FROM pipelines
	WHERE id = $1`
	
	var pipelineData pipelineSQL
	err := pR.db.Get(&pipelineData, getByIDStmt, pipelineID)
	if err != nil {
		return nil, err
	}

	return toPipeline(pipelineData), nil
}

func (pR *PipelineRepository) GetByProjectID(projectID int64) ([]pipelines.Pipeline, error) {
	const getByProjectIDStmt = `
	SELECT
		id, project_id , display_name, data, description,created_at, updated_at, created_by
	FROM pipelines
	WHERE project_id = $1`
	pipelinesSQL := []pipelineSQL{}
	err := pR.db.Select(&pipelinesSQL, getByProjectIDStmt, projectID)
	if err != nil {
		return nil, err
	}
	pipelines := make([]pipelines.Pipeline, len(pipelinesSQL))
	for i := 0; i < len(pipelinesSQL); i++ {
		pipelines[i] = *toPipeline(pipelinesSQL[i])
	}

	return pipelines, nil
}

func (pR *PipelineRepository) Create(p pipelines.Pipeline) (*pipelines.Pipeline, error) {
	p.CreatedAt = time.Now()
	pipelineData := fromPipeline(p)

	const createPipelineStmt = `
	INSERT INTO pipelines (
		project_id, display_name, data, description, created_at, created_by
	) VALUES (
		:project_id, :display_name, :data, :description, :created_at, :created_by
	) RETURNING id`

	query, args, err := sqlx.Named(createPipelineStmt, pipelineData)
	if err != nil {
		return nil, err
	}
	query = sqlx.Rebind(sqlx.DOLLAR, query)
	err = pR.db.Get(&p.ID, query, args...)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (pR *PipelineRepository) Update(pipelineID int64, p pipelines.Pipeline) (*pipelines.Pipeline, error) {
	p.ID = pipelineID
	p.UpdatedAt = time.Now()
	pipelineData := fromPipeline(p)

	const updatePipelineStmt  = `
	UPDATE pipelines
	SET
		display_name 	= COALESCE(:display_name, display_name),
		description 	= COALESCE(:description, description),
		data 			= COALESCE(:data, data),
		updated_at 		= :updated_at
	WHERE id  = :id;`

	result, err := pR.db.NamedExec(updatePipelineStmt, pipelineData)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected != 1 {
		return nil, errors.New("Invalid ID")
	}

	pipeline, err := pR.GetByID(pipelineID)
	if err != nil {
		return nil, err
	}

	return pipeline, nil
}

func (pR *PipelineRepository) Delete(pipelineID int64) error {
	const deletePipelineStmt = `
		DELETE FROM pipelines
		WHERE id = $1
	`
	result, err := pR.db.Exec(deletePipelineStmt, pipelineID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected != 1 {
		return errors.New("Invalid ID")
	}

	return nil
}

type pipelineSQL struct {
	ID          int64          `db:"id"`
	DisplayName sql.NullString `db:"display_name"`
	Data        sql.NullString `db:"data"`
	Description sql.NullString `db:"description"`
	ProjectID   sql.NullInt64  `db:"project_id"`
	CreatedBy   int64          `db:"created_by"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   sql.NullTime   `db:"updated_at"`
}

func toPipeline(pSQL pipelineSQL) *pipelines.Pipeline {
	return &pipelines.Pipeline{
		ID:          pSQL.ID,
		DisplayName: pSQL.DisplayName.String,
		Data:        pSQL.Data.String,
		Description: pSQL.Description.String,
		ProjectID:   pSQL.ProjectID.Int64,
		CreatedBy:   pSQL.CreatedBy,
		CreatedAt:   pSQL.CreatedAt,
		UpdatedAt:   pSQL.UpdatedAt.Time,
	}
}

func fromPipeline(p pipelines.Pipeline) *pipelineSQL {
	var pSQL pipelineSQL
	pSQL.ID = p.ID
	pSQL.CreatedBy = p.CreatedBy
	pSQL.CreatedAt = p.CreatedAt
	pSQL.UpdatedAt = sql.NullTime{
		Time:  p.UpdatedAt,
		Valid: !p.UpdatedAt.IsZero(),
	}
	pSQL.Description = sql.NullString{
		String: p.Description,
		Valid:  p.Description != "",
	}
	pSQL.Data = sql.NullString{
		String: p.Data,
		Valid:  (p.Data != "" || p.Data != "{}"),
	}
	pSQL.ProjectID = sql.NullInt64{
		Int64: p.ProjectID,
		Valid: p.ProjectID != 0,
	}
	pSQL.DisplayName = sql.NullString{
		String: p.DisplayName,
		Valid:  p.DisplayName != "",
	}

	return &pSQL
}
