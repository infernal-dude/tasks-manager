package repository

import (
	"database/sql"
	"tasks-manager/internal/domain"

	"github.com/jmoiron/sqlx"
)

type TaskRepository interface {
	Create(task *domain.Task) error
	GetById(id int64) (*domain.Task, error)
	GetAll() ([]domain.Task, error)
	Update(*domain.Task) error
	Delete(id int64) error
}

type taskRepository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (t *taskRepository) Create(task *domain.Task) error {
	row := t.db.QueryRowx("INSERT INTO tasks(title, description) VALUES($1, $2) RETURNING id, created_at", task.Title, task.Description)
	if err := row.Scan(&task.ID, &task.CreatedAt); err != nil {
		return err
	}
	return nil
}

func (t *taskRepository) GetById(id int64) (*domain.Task, error) {
	var task domain.Task
	err := t.db.Get(&task, "SELECT id, title, description, created_at, completed FROM tasks WHERE id=$1", id)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (t *taskRepository) GetAll() ([]domain.Task, error) {
	var tasks []domain.Task
	err := t.db.Select(&tasks, "SELECT id, title, description, created_at, completed FROM tasks")
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (t *taskRepository) Update(task *domain.Task) error {
	result, err := t.db.Exec("UPDATE tasks set title=$1, description=$2, completed=$3 WHERE id = $4", task.Title, task.Description, task.Completed, task.ID)
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

func (t *taskRepository) Delete(id int64) error {
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
