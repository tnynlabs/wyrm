package users

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"regexp"
	"time"

	"github.com/tnynlabs/wyrm/pkg/utils"
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
	Update(userID int64, u User) (*User, error)
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
	user, err := s.userRepo.GetByKey(key)
	if err != nil {
		return nil, &utils.ServiceErr{
			Code:    UserNotFoundCode,
			Message: "Invalid Key",
		}
	}

	return user, nil
}

func (s *service) GetByID(userID int64) (*User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, &utils.ServiceErr{
			Code:    UserNotFoundCode,
			Message: "Invalid ID",
		}
	}

	return user, nil
}

func (s *service) CreateWithPwd(u User, pwd string) (*User, error) {
	if u.Name == "" {
		return nil, &utils.ServiceErr{
			Code:    InvalidInputCode,
			Message: "Invalid name",
		}
	}

	if u.DisplayName == "" {
		return nil, &utils.ServiceErr{
			Code:    InvalidInputCode,
			Message: "Invalid display name",
		}
	}

	if !isStrongPwd(pwd) {
		return nil, &utils.ServiceErr{
			Code:    InvalidInputCode,
			Message: "Weak password",
		}
	}

	if !isValidEmail(u.Email) {
		return nil, &utils.ServiceErr{
			Code:    InvalidInputCode,
			Message: "Invalid email",
		}
	}

	if dup, err := s.userRepo.IsDuplicateEmail(u.Email); err != nil || dup {
		return nil, &utils.ServiceErr{
			Code:    DuplicateEmailCode,
			Message: "User with this email already exists",
		}
	}

	if dup, err := s.userRepo.IsDuplicateName(u.Name); err != nil || dup {
		return nil, &utils.ServiceErr{
			Code:    DuplicateNameCode,
			Message: "User with this name already exists",
		}
	}

	u.setPwd(pwd)
	u.AuthKey = genString(64) // Let's hope no collisions lol

	newUser, err := s.userRepo.Create(u)
	if err != nil {
		// TODO: better error handling
		log.Printf("Failed creating new user (error: %v)", err)
		return nil, &utils.ServiceErr{
			Code:    utils.UnexpectedCode,
			Message: "Failed creating new user",
		}
	}

	return newUser, nil
}

func (s *service) Update(userID int64, u User) (*User, error) {
	if u.DisplayName == "" {
		return nil, &utils.ServiceErr{
			Code:    InvalidInputCode,
			Message: "Invalid display name",
		}
	}

	if !isValidEmail(u.Email) {
		return nil, &utils.ServiceErr{
			Code:    InvalidInputCode,
			Message: "Invalid email",
		}
	}

	// Initialize updatable fields only
	updatedData := User{
		Email:       u.Email,
		DisplayName: u.DisplayName,
	}

	user, err := s.userRepo.Update(userID, updatedData)
	if err != nil {
		return nil, &utils.ServiceErr{
			Code:    UserNotFoundCode,
			Message: "Invalid ID",
		}
	}

	return user, nil
}

func (s *service) Delete(userID int64) error {
	err := s.userRepo.Delete(userID)
	if err != nil {
		return &utils.ServiceErr{
			Code:    UserNotFoundCode,
			Message: "Invalid ID",
		}
	}

	return nil
}

func (s *service) AuthWithEmailPwd(email, pwd string) (*User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, &utils.ServiceErr{
			Code:    InvalidInputCode,
			Message: "Invalid email or password",
		}
	}

	if !user.isValidPwd(pwd) {
		return nil, &utils.ServiceErr{
			Code:    InvalidInputCode,
			Message: "Invalid email or password",
		}
	}

	return user, nil
}

func isStrongPwd(pwd string) bool {
	if len(pwd) < minPwdLength {
		return false
	}

	if matched, err := regexp.MatchString("[0-9]", pwd); err != nil || !matched {
		return false
	}

	if matched, err := regexp.MatchString("[a-zA-Z]", pwd); err != nil || !matched {
		return false
	}

	return true
}

func isValidEmail(email string) bool {
	// https://html.spec.whatwg.org/multipage/input.html#valid-e-mail-address
	// https://github.com/badoux/checkmail/blob/master/checkmail.go
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return re.MatchString(email)
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
