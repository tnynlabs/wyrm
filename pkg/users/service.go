package users

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"time"
)

const minPwdLength = 8

// User Contains user core properties
// Note: zero values will not be updated
type User struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time

	// Note: Only these values are updatable
	Email       string
	DisplayName string

	// Note: Never show in output
	AuthKey string
	PwdHash string
	PwdSalt string
}

// setPwd generates a salted hash of the password
// and fills PwdSalt and PwdHash accordingly
func (u *User) setPwd(pwd string) {
	u.PwdSalt = genString(8)
	u.PwdHash = genPwdHash(pwd, u.PwdSalt)
}

func (u *User) isValidPwd(pwd string) bool {
	return genPwdHash(pwd, u.PwdSalt) == u.PwdHash
}

// Repository defines the user.Repository operations
// Storage implementations should follow this interface (e.g. Postgres, In Memory, ...etc)
type Repository interface {
	GetByID(userID int64) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByKey(key string) (*User, error)
	Create(u User) (*User, error)
	Update(u User) (*User, error)
	Delete(userID int64) error

	// Check if email already exists
	IsDuplicateEmail(email string) (bool, error)
	// Check if name already exists
	IsDuplicateName(name string) (bool, error)
}

// Service defines the user.Service operations
type Service interface {
	// CRUD ops
	GetByKey(key string) (*User, error)
	GetByID(userID int64) (*User, error)
	CreateWithPwd(u User, pwd string) (*User, error)
	Update(userID int64, u User) (*User, error)
	Delete(userID int64) error

	AuthWithEmailPwd(email, pwd string) (*User, error)
}

type service struct {
	userRepo Repository
}

// CreateService Create new instance of User Service
func CreateService(repo Repository) Service {
	return &service{repo}
}

func (s *service) GetByKey(key string) (*User, error) {
	return nil, nil
}

func (s *service) GetByID(userID int64) (*User, error) {
	return nil, nil
}

func (s *service) CreateWithPwd(u User, pwd string) (*User, error) {
	return nil, nil
}

func (s *service) Update(userID int64, u User) (*User, error) {
	return nil, nil
}

func (s *service) Delete(userID int64) error {
	return nil
}

func (s *service) AuthWithEmailPwd(email, pwd string) (*User, error) {
	return nil, nil
}

func genString(size int) string {
	randBytes := make([]byte, size)
	rand.Read(randBytes)
	return base64.StdEncoding.EncodeToString(randBytes)
}

func genPwdHash(pwd, salt string) string {
	hashSlice := sha256.Sum256([]byte(pwd + salt))
	hash := base64.StdEncoding.EncodeToString(hashSlice[:])
	return hash
}
