package repository

import (
	"errors"
	"sync"
	"time"

	"github.com/popvictor123/distributed-calc/internal/orchestrator/models"
	"github.com/google/uuid"
)

var (
	ErrExpressionNotFound = errors.New("expression not found")
	ErrTaskNotFound       = errors.New("task not found")
	ErrNoTasksAvailable   = errors.New("no tasks available")
)

type Repository struct {
	expressions     map[uuid.UUID]*models.Expression
	tasks           map[uuid.UUID]*models.Task
	expressionMutex sync.RWMutex
	taskMutex       sync.RWMutex
}

func NewRepository() *Repository {
	return &Repository{
		expressions: make(map[uuid.UUID]*models.Expression),
		tasks:       make(map[uuid.UUID]*models.Task),
	}
}

func (r *Repository) CreateExpression(expression string) (*models.Expression, error) {
	r.expressionMutex.Lock()
	defer r.expressionMutex.Unlock()

	expr := &models.Expression{
		ID:         uuid.New(),
		Expression: expression,
		Status:     models.StatusPending,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	r.expressions[expr.ID] = expr
	return expr, nil
}

func (r *Repository) GetExpressionByID(id uuid.UUID) (*models.Expression, error) {
	r.expressionMutex.RLock()
	defer r.expressionMutex.RUnlock()

	expr, exists := r.expressions[id]
	if !exists {
		return nil, ErrExpressionNotFound
	}
	return expr, nil
}

func (r *Repository) GetAllExpressions() []*models.Expression {
	r.expressionMutex.RLock()
	defer r.expressionMutex.RUnlock()

	result := make([]*models.Expression, 0, len(r.expressions))
	for _, expr := range r.expressions {
		result = append(result, expr)
	}
	return result
}

func (r *Repository) UpdateExpression(expr *models.Expression) error {
	r.expressionMutex.Lock()
	defer r.expressionMutex.Unlock()

	if _, exists := r.expressions[expr.ID]; !exists {
		return ErrExpressionNotFound
	}

	expr.UpdatedAt = time.Now()
	r.expressions[expr.ID] = expr
	return nil
}

func (r *Repository) SaveTasks(tasks []*models.Task) error {
	r.taskMutex.Lock()
	defer r.taskMutex.Unlock()

	for _, task := range tasks {
		r.tasks[task.ID] = task
	}
	return nil
}

func (r *Repository) GetTaskByID(id uuid.UUID) (*models.Task, error) {
	r.taskMutex.RLock()
	defer r.taskMutex.RUnlock()

	task, exists := r.tasks[id]
	if !exists {
		return nil, ErrTaskNotFound
	}
	return task, nil
}

func (r *Repository) UpdateTask(task *models.Task) error {
	r.taskMutex.Lock()
	defer r.taskMutex.Unlock()

	if _, exists := r.tasks[task.ID]; !exists {
		return ErrTaskNotFound
	}

	r.tasks[task.ID] = task
	return nil
}

func (r *Repository) GetNextPendingTask() (*models.Task, error) {
	r.taskMutex.Lock()
	defer r.taskMutex.Unlock()

	for _, task := range r.tasks {
		if task.Status == models.TaskStatusPending {
			allDependenciesResolved := true
			if task.Dependencies != nil {
				for _, depID := range task.Dependencies {
					if depTask, exists := r.tasks[*depID]; exists {
						if depTask.Status != models.TaskStatusCompleted {
							allDependenciesResolved = false
							break
						}
						
						if depTask.Result != nil {
							if task.Arg1 != nil && task.Arg1.ID == *depID {
								task.Arg1Value = *depTask.Result
							}
							if task.Arg2 != nil && task.Arg2.ID == *depID {
								task.Arg2Value = *depTask.Result
							}
						}
					} else {
						allDependenciesResolved = false
						break
					}
				}
			}

			if allDependenciesResolved {
				now := time.Now()
				task.Status = models.TaskStatusProcessing
				task.StartedAt = &now
				return task, nil
			}
		}
	}

	return nil, ErrNoTasksAvailable
}

func (r *Repository) GetTasksByExpressionID(expressionID uuid.UUID) []*models.Task {
	r.taskMutex.RLock()
	defer r.taskMutex.RUnlock()

	var result []*models.Task
	for _, task := range r.tasks {
		if task.ExpressionID == expressionID {
			result = append(result, task)
		}
	}
	return result
}

func (r *Repository) CheckExpressionCompletion(expressionID uuid.UUID) {
	r.expressionMutex.Lock()
	defer r.expressionMutex.Unlock()
	
	expr, exists := r.expressions[expressionID]
	if !exists {
		return
	}

	tasks := r.GetTasksByExpressionID(expressionID)
	
	allCompleted := true
	var finalTask *models.Task
	
	for _, task := range tasks {
		if task.Status != models.TaskStatusCompleted {
			allCompleted = false
			break
		}
		
		isDependency := false
		for _, t := range tasks {
			if t.Dependencies != nil {
				for _, depID := range t.Dependencies {
					if *depID == task.ID {
						isDependency = true
						break
					}
				}
			}
			if isDependency {
				break
			}
		}
		
		if !isDependency {
			finalTask = task
		}
	}
	
	if allCompleted && finalTask != nil && finalTask.Result != nil {
		expr.Status = models.StatusCompleted
		expr.Result = finalTask.Result
		expr.UpdatedAt = time.Now()
	} else if expr.Status == models.StatusPending {
		expr.Status = models.StatusComputing
		expr.UpdatedAt = time.Now()
	}
}
