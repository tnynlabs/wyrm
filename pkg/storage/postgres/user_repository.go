package postgres

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/tnynlabs/wyrm/pkg/users"
)

// UserRepository users.Repository Postgres implementation
type UserRepository struct {
	db *sqlx.DB
}

// CreateUserRepository Create new instance of postgres.UserRepository
func CreateUserRepository(db *sqlx.DB) users.Repository {
	return &UserRepository{db}
}

func (uR *UserRepository) GetByID(userID int64) (*users.User, error) {
	return nil, nil
}

func (uR *UserRepository) GetByEmail(email string) (*users.User, error) {
	return nil, nil
}

func (uR *UserRepository) GetByKey(key string) (*users.User, error) {
	return nil, nil
}

func (uR *UserRepository) Create(u users.User) (*users.User, error) {
	return nil, nil
}

func (uR *UserRepository) Update(u users.User) (*users.User, error) {
	return nil, nil
}

func (uR *UserRepository) Delete(userID int64) error {
	return nil
}

// TODO: cache for fast checks
func (uR *UserRepository) IsDuplicateEmail(email string) (bool, error) {
	return false, nil
}

// TODO: cache for fast checks
func (uR *UserRepository) IsDuplicateName(name string) (bool, error) {
	return false, nil
}

type userSQL struct {
	ID          int64          `db:"id"`
	Name        string         `db:"name"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   sql.NullTime   `db:"updated_at"`
	Email       sql.NullString `db:"email"`
	DisplayName sql.NullString `db:"display_name"`
	AuthKey     sql.NullString `db:"auth_key"`
	PwdHash     sql.NullString `db:"pwd_hash"`
	PwdSalt     sql.NullString `db:"pwd_salt"`
}

func toUser(uSQL userSQL) *users.User {
	return &users.User{
		ID:          uSQL.ID,
		Name:        uSQL.Name,
		CreatedAt:   uSQL.CreatedAt,
		UpdatedAt:   uSQL.UpdatedAt.Time,
		Email:       uSQL.Email.String,
		DisplayName: uSQL.DisplayName.String,
		AuthKey:     uSQL.AuthKey.String,
		PwdHash:     uSQL.PwdHash.String,
		PwdSalt:     uSQL.PwdSalt.String,
	}
}

func fromUser(u users.User) *userSQL {
	var uSQL userSQL
	uSQL.ID = u.ID
	uSQL.Name = u.Name
	uSQL.CreatedAt = u.CreatedAt
	uSQL.UpdatedAt = sql.NullTime{
		Time:  u.UpdatedAt,
		Valid: !u.UpdatedAt.IsZero(),
	}
	uSQL.Email = sql.NullString{
		String: u.Email,
		Valid:  u.Email != "",
	}
	uSQL.DisplayName = sql.NullString{
		String: u.DisplayName,
		Valid:  u.DisplayName != "",
	}
	uSQL.AuthKey = sql.NullString{
		String: u.AuthKey,
		Valid:  u.AuthKey != "",
	}
	uSQL.PwdHash = sql.NullString{
		String: u.PwdHash,
		Valid:  u.PwdHash != "",
	}
	uSQL.PwdSalt = sql.NullString{
		String: u.PwdSalt,
		Valid:  u.PwdSalt != "",
	}
	return &uSQL
}
