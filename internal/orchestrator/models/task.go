package models

import (
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "PENDING"
	TaskStatusProcessing TaskStatus = "PROCESSING"
	TaskStatusCompleted  TaskStatus = "COMPLETED"
	TaskStatusError      TaskStatus = "ERROR"
)

type OperationType string

const (
	OperationAddition       OperationType = "ADDITION"
	OperationSubtraction    OperationType = "SUBTRACTION"
	OperationMultiplication OperationType = "MULTIPLICATION"
	OperationDivision       OperationType = "DIVISION"
	OperationValue          OperationType = "VALUE" // Just a value, no operation
)

type Task struct {
	ID            uuid.UUID     `json:"id"`
	ExpressionID  uuid.UUID     `json:"-"`
	Arg1          *Task         `json:"-"`
	Arg1Value     float64       `json:"arg1"`
	Arg2          *Task         `json:"-"`  
	Arg2Value     float64       `json:"arg2"`
	Operation     OperationType `json:"operation"`
	OperationTime int           `json:"operation_time"`
	Status        TaskStatus    `json:"status"`
	Result        *float64      `json:"result,omitempty"`
	Error         string        `json:"error,omitempty"`
	Dependencies  []*uuid.UUID  `json:"-"`
	CreatedAt     time.Time     `json:"-"`
	StartedAt     *time.Time    `json:"-"`
	CompletedAt   *time.Time    `json:"-"`
}

type TaskResponse struct {
	ID            uuid.UUID     `json:"id"`
	Arg1          float64       `json:"arg1"`
	Arg2          float64       `json:"arg2"`
	Operation     OperationType `json:"operation"`
	OperationTime int           `json:"operation_time"`
}

type GetTaskResponse struct {
	Task *TaskResponse `json:"task,omitempty"`
}

type TaskResultRequest struct {
	ID     uuid.UUID `json:"id"`
	Result float64   `json:"result"`
}
