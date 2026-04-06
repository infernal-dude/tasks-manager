package service

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"tasks-manager/internal/domain"
	"tasks-manager/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(user *domain.User) error
	Login(user domain.User) (string, error)
	GetByUsername(username string) (*domain.User, error)
	GetById(id int64) (*domain.User, error)
	GetAll() ([]domain.User, error)
	Update(user *domain.User) error
	Delete(id int64) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (u *userService) Register(user *domain.User) error {
	_, err := u.repo.GetByUsername(strings.TrimSpace(user.Username))
	if err == nil {
		return fmt.Errorf("Username %s already taken!", user.Username)
	}
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	hashedPassword, err := GeneratePassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	err = u.repo.Create(user)
	if err != nil {
		return err
	}
	return nil
}

func (u *userService) Login(user domain.User) (string, error) {
	userComp, err := u.repo.GetByUsername(user.Username)
	if err != nil {
		return "", fmt.Errorf("Incorrect username or password!")
	}
	err = bcrypt.CompareHashAndPassword([]byte(userComp.Password), []byte(user.Password))
	if err != nil {
		return "", fmt.Errorf("Incorrect username or password!")
	}

	token, err := GenerateToken(userComp.ID, userComp.Role)
	if err != nil {
		return "", err
	}

	return token, nil

}

func (u *userService) GetByUsername(username string) (*domain.User, error) {
	user, err := u.repo.GetByUsername(username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("There is no user with such username")
		}
		return nil, err
	}

	return user, nil
}

func (u *userService) GetById(id int64) (*domain.User, error) {
	user, err := u.repo.GetById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("There is no user with such id")
		}
		return nil, err
	}

	return user, nil
}

func (u *userService) GetAll() ([]domain.User, error) {
	users, err := u.repo.GetAll()
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, sql.ErrNoRows
	}

	return users, nil
}

func (u *userService) Update(user *domain.User) error {
	userTest, err := u.repo.GetByUsername(user.Username)
	if err == nil && userTest.ID != user.ID {
		return fmt.Errorf("busy by another user")
	}

	hashedPassword, err := GeneratePassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	err = u.repo.Update(user)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no such user")
		}
		return err
	}

	return nil
}

func (u *userService) Delete(id int64) error {
	err := u.repo.Delete(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("No user with such id")
		}
		return err
	}

	return nil
}

func GenerateToken(userID int64, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	signKey := os.Getenv("SIGNATION")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(signKey))
}

func GeneratePassword(pass string) (string, error) {
	if pass == "" {
		return "", fmt.Errorf("invalid password")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	password := string(hashedPassword)

	return password, nil
}
