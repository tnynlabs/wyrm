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
	const getByEmailStmt = `
		SELECT
			id, email, name, display_name, auth_key,
			pwd_hash, pwd_salt, created_at, updated_at
		FROM users
		WHERE id = $1`

	var userData userSQL
	err := uR.db.Get(&userData, getByEmailStmt, userID)
	if err != nil {
		return nil, err
	}

	return toUser(userData), nil
}

func (uR *UserRepository) GetByEmail(email string) (*users.User, error) {
	const getByEmailStmt = `
		SELECT
			id, email, name, display_name, auth_key,
			pwd_hash, pwd_salt, created_at, updated_at
		FROM users
		WHERE email = $1`

	var userData userSQL
	err := uR.db.Get(&userData, getByEmailStmt, email)
	if err != nil {
		return nil, err
	}

	return toUser(userData), nil
}

func (uR *UserRepository) GetByKey(key string) (*users.User, error) {
	const getByEmailStmt = `
		SELECT
			id, email, name, display_name, auth_key,
			pwd_hash, pwd_salt, created_at, updated_at
		FROM users
		WHERE auth_key = $1`

	var userData userSQL
	err := uR.db.Get(&userData, getByEmailStmt, key)
	if err != nil {
		return nil, err
	}

	return toUser(userData), nil
}

func (uR *UserRepository) Create(u users.User) (*users.User, error) {
	u.CreatedAt = time.Now()

	userData := fromUser(u)

	const insertUserStmt = `
		INSERT INTO users (
			name, email, display_name, auth_key,
			pwd_hash, pwd_salt, created_at
		) VALUES (
			:name, :email, :display_name, :auth_key,
			:pwd_hash, :pwd_salt, :created_at
		) RETURNING id`

	query, args, err := sqlx.Named(insertUserStmt, userData)
	if err != nil {
		return nil, err
	}

	// https://pkg.go.dev/github.com/jmoiron/sqlx#Rebind
	// Replace ? with $ for postgres
	query = sqlx.Rebind(sqlx.DOLLAR, query)

	err = uR.db.Get(&u.ID, query, args...)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (uR *UserRepository) Update(u users.User) (*users.User, error) {
	return nil, nil
}

func (uR *UserRepository) Delete(userID int64) error {
	return nil
}

// TODO: cache for fast checks
func (uR *UserRepository) IsDuplicateEmail(email string) (bool, error) {
	const lookupEmailStmt = `SELECT COUNT(id) FROM users WHERE email = $1`

	var cnt int
	err := uR.db.Get(&cnt, lookupEmailStmt, email)
	if err != nil {
		return true, err
	}

	return (cnt != 0), nil
}

// TODO: cache for fast checks
func (uR *UserRepository) IsDuplicateName(name string) (bool, error) {
	const lookupNameStmt = `SELECT COUNT(id) FROM users WHERE name = $1`

	var cnt int
	err := uR.db.Get(&cnt, lookupNameStmt, name)
	if err != nil {
		return true, err
	}

	return (cnt != 0), nil
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
