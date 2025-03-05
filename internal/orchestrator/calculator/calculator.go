package calculator

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/popvictor123/distributed-calc/internal/orchestrator/models"
	"github.com/google/uuid"
)

type Calculator struct {
	AdditionTime       int
	SubtractionTime    int
	MultiplicationTime int
	DivisionTime       int
}

func NewCalculator() *Calculator {
	addTime, _ := strconv.Atoi(getEnvOrDefault("TIME_ADDITION_MS", "1000"))
	subTime, _ := strconv.Atoi(getEnvOrDefault("TIME_SUBTRACTION_MS", "1000"))
	mulTime, _ := strconv.Atoi(getEnvOrDefault("TIME_MULTIPLICATIONS_MS", "2000"))
	divTime, _ := strconv.Atoi(getEnvOrDefault("TIME_DIVISIONS_MS", "2000"))

	return &Calculator{
		AdditionTime:       addTime,
		SubtractionTime:    subTime,
		MultiplicationTime: mulTime,
		DivisionTime:       divTime,
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func (c *Calculator) ProcessExpression(expression string, expressionID uuid.UUID) ([]*models.Task, error) {
	parser := NewParser(expression)
	ast, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	tasks, err := c.convertASTToTasks(ast, expressionID)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (c *Calculator) convertASTToTasks(node ASTNode, expressionID uuid.UUID) ([]*models.Task, error) {
	var tasks []*models.Task
	switch n := node.(type) {
	case *NumberNode:
		task := &models.Task{
			ID:           uuid.New(),
			ExpressionID: expressionID,
			Operation:    models.OperationValue,
			Status:       models.TaskStatusCompleted, // Numbers are already calculated
			CreatedAt:    time.Now(),
		}
		result := n.Value
		task.Result = &result
		tasks = append(tasks, task)
		return tasks, nil

	case *BinaryOpNode:
		leftTasks, err := c.convertASTToTasks(n.Left, expressionID)
		if err != nil {
			return nil, err
		}

		rightTasks, err := c.convertASTToTasks(n.Right, expressionID)
		if err != nil {
			return nil, err
		}

		leftTask := leftTasks[len(leftTasks)-1]
		rightTask := rightTasks[len(rightTasks)-1]

		operationType, opTime := c.getOperationTypeAndTime(n.Op)
		
		task := &models.Task{
			ID:            uuid.New(),
			ExpressionID:  expressionID,
			Arg1:          leftTask,
			Arg2:          rightTask,
			Operation:     operationType,
			OperationTime: opTime,
			Status:        models.TaskStatusPending,
			Dependencies:  []*uuid.UUID{&leftTask.ID, &rightTask.ID},
			CreatedAt:     time.Now(),
		}

		if leftTask.Result != nil {
			task.Arg1Value = *leftTask.Result
		}
		
		if rightTask.Result != nil {
			task.Arg2Value = *rightTask.Result
		}

		tasks = append(tasks, leftTasks...)
		tasks = append(tasks, rightTasks...)
		tasks = append(tasks, task)
		
		return tasks, nil
	}

	return nil, fmt.Errorf("unknown node type")
}

func (c *Calculator) getOperationTypeAndTime(op string) (models.OperationType, int) {
	switch op {
	case "+":
		return models.OperationAddition, c.AdditionTime
	case "-":
		return models.OperationSubtraction, c.SubtractionTime
	case "*":
		return models.OperationMultiplication, c.MultiplicationTime
	case "/":
		return models.OperationDivision, c.DivisionTime
	default:
		return "", 0
	}
}

func (c *Calculator) ExecuteOperation(op models.OperationType, arg1, arg2 float64) (float64, error) {
	switch op {
	case models.OperationAddition:
		return arg1 + arg2, nil
	case models.OperationSubtraction:
		return arg1 - arg2, nil
	case models.OperationMultiplication:
		return arg1 * arg2, nil
	case models.OperationDivision:
		if arg2 == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return arg1 / arg2, nil
	case models.OperationValue:
		return arg1, nil // Just return the value
	default:
		return 0, fmt.Errorf("unknown operation type: %s", op)
	}
}

func (c *Calculator) UpdateTaskResult(task *models.Task, result float64) {
	task.Result = &result
	task.Status = models.TaskStatusCompleted
	now := time.Now()
	task.CompletedAt = &now
}
