package repository

import (
	"database/sql"
	"tasks-manager/internal/domain"

	"github.com/jmoiron/sqlx"
)

type TaskRepository interface {
	Create(task *domain.Task) error
	GetById(id int64, userID int64, isAdmin bool) (*domain.Task, error)
	GetAll(userID int64) ([]domain.Task, error)
	Update(task *domain.Task, userID int64) error
	Delete(id int64, userID int64, isAdmin bool) error
}

type taskRepository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (t *taskRepository) Create(task *domain.Task) error {
	row := t.db.QueryRowx("INSERT INTO tasks(title, description, user_id) VALUES($1, $2, $3) RETURNING id, created_at", task.Title, task.Description, task.UserId)
	if err := row.Scan(&task.ID, &task.CreatedAt); err != nil {
		return err
	}
	return nil
}

func (t *taskRepository) GetById(id int64, userID int64, isAdmin bool) (*domain.Task, error) {
	var task domain.Task

	if isAdmin {
		err := t.db.Get(&task, "SELECT id, title, description, created_at, completed FROM tasks WHERE id=$1", id)
		if err != nil {
			return nil, err
		}
		return &task, nil
	}

	err := t.db.Get(&task, "SELECT id, title, description, created_at, completed FROM tasks WHERE id=$1 AND user_id=$2", id, userID)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (t *taskRepository) GetAll(userID int64) ([]domain.Task, error) {
	var tasks []domain.Task
	err := t.db.Select(&tasks, "SELECT id, title, description, created_at, completed FROM tasks WHERE user_id=$1", userID)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (t *taskRepository) Update(task *domain.Task, userID int64) error {
	result, err := t.db.Exec("UPDATE tasks set title=$1, description=$2, completed=$3 WHERE id = $4 AND user_id=$5", task.Title, task.Description, task.Completed, task.ID, userID)
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

func (t *taskRepository) Delete(id int64, userID int64, isAdmin bool) error {
	if isAdmin {
		result, err := t.db.Exec("DELETE FROM tasks WHERE id=$1", id)
		if err != nil {
			return err
		}
		row, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if row == 0 {
			return sql.ErrNoRows
		}
		return nil
	}

	result, err := t.db.Exec("DELETE FROM tasks WHERE id=$1 AND user_id=$2", id, userID)
	if err != nil {
		return err
	}
	row, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if row == 0 {
		return sql.ErrNoRows
	}

	return nil
}
