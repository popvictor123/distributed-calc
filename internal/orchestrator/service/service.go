package service

import (

	"github.com/popvictor123/distributed-calc/internal/orchestrator/calculator"
	"github.com/popvictor123/distributed-calc/internal/orchestrator/models"
	"github.com/popvictor123/distributed-calc/internal/orchestrator/repository"
	"github.com/google/uuid"
)

type Service struct {
	repo      *repository.Repository
	calculator *calculator.Calculator
}

func NewService(repo *repository.Repository, calc *calculator.Calculator) *Service {
	return &Service{
		repo:      repo,
		calculator: calc,
	}
}

func (s *Service) CalculateExpression(expression string) (*models.Expression, error) {
	expr, err := s.repo.CreateExpression(expression)
	if err != nil {
		return nil, err
	}

	tasks, err := s.calculator.ProcessExpression(expression, expr.ID)
	if err != nil {
		expr.Status = models.StatusError
		expr.Error = err.Error()
		s.repo.UpdateExpression(expr)
		return nil, err
	}

	err = s.repo.SaveTasks(tasks)
	if err != nil {
		return nil, err
	}

	s.repo.CheckExpressionCompletion(expr.ID)

	return expr, nil
}

func (s *Service) GetExpressionByID(id uuid.UUID) (*models.Expression, error) {
	return s.repo.GetExpressionByID(id)
}

func (s *Service) GetAllExpressions() []*models.Expression {
	return s.repo.GetAllExpressions()
}

func (s *Service) GetNextTask() (*models.Task, error) {
	return s.repo.GetNextPendingTask()
}

func (s *Service) UpdateTaskResult(id uuid.UUID, result float64) error {
	task, err := s.repo.GetTaskByID(id)
	if err != nil {
		return err
	}

	s.calculator.UpdateTaskResult(task, result)
	
	err = s.repo.UpdateTask(task)
	if err != nil {
		return err
	}

	s.repo.CheckExpressionCompletion(task.ExpressionID)
	
	return nil
}
