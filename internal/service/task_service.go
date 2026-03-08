package service

import (
	"errors"
	"strings"
	"tasks-manager/internal/domain"
	"tasks-manager/internal/repository"
)

type TaskService interface {
	Create(task *domain.Task) error
	GetById(id int64) (*domain.Task, error)
	GetAll() ([]domain.Task, error)
	Update(*domain.Task) error
	Delete(id int64) error
}

type taskService struct {
	repo repository.TaskRepository
}

func NewService(repo repository.TaskRepository) TaskService {
	return &taskService{repo: repo}
}

func (s *taskService) Create(task *domain.Task) error {
	if strings.TrimSpace(task.Title) == "" {
		return errors.New("Field \"Title\" can't be empty")
	}

	if err := s.repo.Create(task); err != nil {
		return err
	}
	return nil
}

func (s *taskService) GetById(id int64) (*domain.Task, error) {
	return s.repo.GetById(id)
}

func (s *taskService) GetAll() ([]domain.Task, error) {
	return s.repo.GetAll()
}

func (s *taskService) Update(task *domain.Task) error {
	if strings.TrimSpace(task.Title) == "" {
		return errors.New("Field \"Title\" can't be empty")
	}

	if err := s.repo.Update(task); err != nil {
		return err
	}
	return nil
}

func (s *taskService) Delete(id int64) error {
	task, err := s.repo.GetById(id)
	if err != nil {
		return err
	}

	if task.Completed {
		return errors.New("Can't delete completed task")
	}

	if err := s.repo.Delete(id); err != nil {
		return err
	}
	return nil
}
