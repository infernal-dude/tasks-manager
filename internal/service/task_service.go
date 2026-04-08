package service

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"tasks-manager/internal/domain"
	"tasks-manager/internal/repository"
)

// C TaskService и UserService та же картина что и с *Handler.
// Почему они разделены? Судя по названиям методов и наличию интерфейса это больше походе на репозиторий, чем на инкапсуляцию бизнес логики.
// Какую роль здесь играет интерфейс? Таск в будущем можно создать как-то подругому? Кажется что это здесь лишнее
type TaskService interface {
	Create(task *domain.Task, userID int64) error
	GetById(id int64, userID int64) (*domain.Task, error)
	GetAll(userID int64) ([]domain.Task, error)
	Update(task *domain.Task, userID int64) error
	Delete(id int64, userID int64) error
}

type taskService struct {
	repo repository.TaskRepository
}

func NewService(repo repository.TaskRepository) TaskService {
	return &taskService{repo: repo}
}

func (s *taskService) Create(task *domain.Task, userID int64) error {
	if strings.TrimSpace(task.Title) == "" {
		return errors.New("Field \"Title\" can't be empty")
	}
	task.UserId = userID
	if err := s.repo.Create(task); err != nil {
		return err
	}
	return nil
}

func (s *taskService) GetById(id int64, userID int64) (*domain.Task, error) {
	return s.repo.GetById(id, userID)
}

func (s *taskService) GetAll(userID int64) ([]domain.Task, error) {
	return s.repo.GetAll(userID)
}

func (s *taskService) Update(task *domain.Task, userID int64) error {
	if strings.TrimSpace(task.Title) == "" {
		return errors.New("Field \"Title\" can't be empty")
	}

	if err := s.repo.Update(task, userID); err != nil {
		return err
	}
	return nil
}

func (s *taskService) Delete(id int64, userID int64) error {

	task, err := s.repo.GetById(id, userID)
	if err != nil {
		return err
	}

	if task.Completed {
		return errors.New("Can't delete completed task")
	}

	if err := s.repo.Delete(id, userID); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("task not found")
		}
		return err
	}
	return nil
}
