package postgres

import (
	"database/sql"
	"time"
	"github.com/jmoiron/sqlx"
	"github.com/tnynlabs/wyrm/pkg/pipelines"
)

type PipelineRepository struct {
	db *sqlx.DB
}

func CreatePipelineRepository(db *sqlx.DB) pipelines.Repository {
	return &PipelineRepository{db}
}
func (pR *PipelineRepository) GetByID(pipelineID int64) (*pipelines.Pipeline, error) {
	return nil, nil
}
func (pR *PipelineRepository) GetByProjectID(projectID int64) ([]pipelines.Pipeline, error) {
	return nil, nil
}
func (pR *PipelineRepository) Create(pipeline pipelines.Pipeline) (*pipelines.Pipeline, error) {
	pipeline.CreatedAt = time.Now()
	pipelineData := fromPipeline(pipeline)

	const sqlStmt = `
	INSERT INTO pipelines (
		project_id, display_name, data, description, created_at, created_by
	) VALUES (
		:project_id, :display_name, :data, :description, :created_at, :created_by
	) RETURNING id`
	
	query, args, err := sqlx.Named(sqlStmt, pipelineData)
	if err != nil {
		return nil, err
	}
	query = sqlx.Rebind(sqlx.DOLLAR, query)
	err = pR.db.Get(&pipeline.ID, query, args...)
	if err != nil {
		return nil, err
	}

	return &pipeline, nil
}
func (pR *PipelineRepository) Update(pipelineID int64, pipeline pipelines.Pipeline) (*pipelines.Pipeline, error) {
	return nil, nil
}
func (pR *PipelineRepository) Delete(pipelineID int64) error {
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
		Valid:  p.Data != "",
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
