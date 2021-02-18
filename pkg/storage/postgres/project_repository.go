package postgres

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/tnynlabs/wyrm/pkg/projects"
)

// ProjectRepository users.Repository Postgres implementation
type ProjectRepository struct {
	db *sqlx.DB
}

// CreateProjectRepository new instance of postgres.UserRepository
func CreateProjectRepository(db *sqlx.DB) projects.Repository {
	return &ProjectRepository{db}
}

func (pR *ProjectRepository) GetByID(projectID int64) (*projects.Project, error) {
	const getByIDStmt = `
	SELECT
		id, display_name, created_at, updated_at, description, created_by
	FROM projects
	WHERE id = $1
	`
	var projectData projectSQL
	err := pR.db.Get(&projectData, getByIDStmt, projectID)
	if err != nil {
		return nil, err
	}
	return toProject(projectData), nil
}

// func (pR *ProjectRepository) GetByCreatorID(CreatorID int64) (*projects.Project, error){
// 	return nil,nil
// }

func (pR *ProjectRepository) GetAllowedProjects(userID int64) (*[]projects.Project, error) {
	const selectProjectsStmt = `
	SELECT *
	FROM projects
	WHERE created_by = $1
	`

	projectsSQL := []projectSQL{}
	err := pR.db.Select(&projectsSQL, selectProjectsStmt, userID)
	if err != nil {
		return nil, err
	}
	projects := make([]projects.Project, len(projectsSQL))
	for i := 0; i < len(projectsSQL); i++ {
		projects[i] = *toProject(projectsSQL[i])
	}
	return &projects, nil
}

func (pR *ProjectRepository) Create(p projects.Project) (*projects.Project, error) {
	p.CreatedAt = time.Now()
	projectData := fromProject(p)

	const insertProjectStmt = `
		INSERT INTO projects (display_name, created_at, created_by, description)
		VALUES (:display_name, :created_at, :created_by, :description)
		RETURNING id
	`
	query, args, err := sqlx.Named(insertProjectStmt, projectData)
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
func (pR *ProjectRepository) Update(ProjectID int64, p projects.Project) (*projects.Project, error) {
	p.ID = ProjectID
	p.UpdatedAt = time.Now()
	projectData := fromProject(p)

	const updateProjectStmt = `
	UPDATE projects
	SET
		display_name 	= COALESCE(:display_name, display_name),
		description 	= COALESCE(:description, description),
		updated_at 		= :updated_at
	WHERE id  = :id;`
	result, err := pR.db.NamedExec(updateProjectStmt, projectData)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected != 1 {
		return nil, errors.New("Invalid ID")
	}

	project, err := pR.GetByID(ProjectID)
	if err != nil {
		return nil, err
	}

	return project, nil
}
func (pR *ProjectRepository) Delete(ProjectID int64) error {
	const deleteProjectStmt = `
		DELETE FROM projects
		WHERE id = $1`

	result, err := pR.db.Exec(deleteProjectStmt, ProjectID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected != 1 {
		return errors.New("Invalid ID")
	}

	return nil
}

type projectSQL struct {
	ID          int64          `db:"id"`
	CreatedBy   int64          `db:"created_by"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   sql.NullTime   `db:"updated_at"`
	DisplayName sql.NullString `db:"display_name"`
	Description sql.NullString `db:"description"`
}

func toProject(pSQL projectSQL) *projects.Project {
	return &projects.Project{
		ID:          pSQL.ID,
		CreatedAt:   pSQL.CreatedAt,
		UpdatedAt:   pSQL.UpdatedAt.Time,
		DisplayName: pSQL.DisplayName.String,
		Description: pSQL.Description.String,
		CreatedBy:   pSQL.CreatedBy,
	}
}

func fromProject(p projects.Project) *projectSQL {
	var pSQL projectSQL
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
	pSQL.DisplayName = sql.NullString{
		String: p.DisplayName,
		Valid:  p.DisplayName != "",
	}
	return &pSQL
}
