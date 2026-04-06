package repository

import (
	"database/sql"
	"tasks-manager/internal/domain"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type UserRepository interface {
	Create(user *domain.User) error
	GetByUsername(username string) (*domain.User, error)
	GetById(id int64) (*domain.User, error)
	GetAll() ([]domain.User, error)
	Update(user *domain.User) error
	Delete(id int64) error
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (u *userRepository) Create(user *domain.User) error {
	row := u.db.QueryRowx("INSERT INTO users(username, password) VALUES ($1, $2) RETURNING id, created_at, role", user.Username, user.Password)
	if err := row.Scan(&user.ID, &user.CreatedAt, &user.Role); err != nil {
		return err
	}
	return nil
}

func (u *userRepository) GetByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := u.db.Get(&user, "SELECT id, username, password, role, created_at, updated_at FROM users WHERE username = $1", username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepository) GetById(id int64) (*domain.User, error) {
	var user domain.User
	err := u.db.Get(&user, "SELECT id, username, password, role, created_at, updated_at FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userRepository) GetAll() ([]domain.User, error) {
	var users []domain.User
	err := u.db.Select(&users, "SELECT id, username, password, role, created_at, updated_at FROM users")
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *userRepository) Update(user *domain.User) error {
	result, err := u.db.Exec("UPDATE users SET username=$1, password=$2, updated_at=NOW() WHERE id=$3", user.Username, user.Password, user.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (u *userRepository) Delete(id int64) error {
	result, err := u.db.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
