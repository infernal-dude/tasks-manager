package service

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"tasks-manager/internal/domain"
	"tasks-manager/internal/repository"
)

type TaskService interface {
	Create(task *domain.Task, userID int64) error
	GetById(id int64, userID int64, isAdmin bool) (*domain.Task, error)
	GetAll(userID int64) ([]domain.Task, error)
	Update(task *domain.Task, userID int64) error
	Delete(id int64, userID int64, isAdmin bool) error
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

func (s *taskService) GetById(id int64, userID int64, isAdmin bool) (*domain.Task, error) {
	return s.repo.GetById(id, userID, isAdmin)
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

func (s *taskService) Delete(id int64, userID int64, isAdmin bool) error {
	//ЧЕКНИ ПОТОМ ЧЕ БУДЕТ ТО С GETBYID и FALSE ADMIN
	err := s.repo.Delete(id, userID, isAdmin)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("task not found")
		}
		return err
	}

	return nil
}
