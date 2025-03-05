package models

import (
	"time"

	"github.com/google/uuid"
)

type ExpressionStatus string

const (
	StatusPending   ExpressionStatus = "PENDING"
	StatusComputing ExpressionStatus = "COMPUTING"
	StatusCompleted ExpressionStatus = "COMPLETED"
	StatusError     ExpressionStatus = "ERROR"
)

type Expression struct {
	ID         uuid.UUID        `json:"id"`
	Expression string           `json:"expression"`
	Status     ExpressionStatus `json:"status"`
	Result     *float64         `json:"result,omitempty"`
	Error      string           `json:"error,omitempty"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
	Tasks      []*Task          `json:"-"`
}

type ExpressionResponse struct {
	ID     uuid.UUID        `json:"id"`
	Status ExpressionStatus `json:"status"`
	Result *float64         `json:"result,omitempty"`
}

type ExpressionsResponse struct {
	Expressions []ExpressionResponse `json:"expressions"`
}

type ExpressionDetailResponse struct {
	Expression ExpressionResponse `json:"expression"`
}

type CalculateRequest struct {
	Expression string `json:"expression"`
}

type CalculateResponse struct {
	ID uuid.UUID `json:"id"`
}
