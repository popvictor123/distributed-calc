package api

import (
	"encoding/json"
	"net/http"

	"github.com/popvictor123/distributed-calc/internal/orchestrator/models"
	"github.com/popvictor123/distributed-calc/internal/orchestrator/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) CalculateHandler(w http.ResponseWriter, r *http.Request) {
	var req models.CalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusUnprocessableEntity, "Invalid request payload")
		return
	}

	if req.Expression == "" {
		respondWithError(w, http.StatusUnprocessableEntity, "Expression is required")
		return
	}

	expr, err := h.service.CalculateExpression(req.Expression)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to process expression")
		return
	}

	resp := models.CalculateResponse{
		ID: expr.ID,
	}

	respondWithJSON(w, http.StatusCreated, resp)
}

func (h *Handler) GetExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	expressions := h.service.GetAllExpressions()
	var response models.ExpressionsResponse
	response.Expressions = make([]models.ExpressionResponse, 0, len(expressions))

	for _, expr := range expressions {
		response.Expressions = append(response.Expressions, models.ExpressionResponse{
			ID:     expr.ID,
			Status: expr.Status,
			Result: expr.Result,
		})
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (h *Handler) GetExpressionByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusUnprocessableEntity, "Invalid expression ID")
		return
	}

	expr, err := h.service.GetExpressionByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Expression not found")
		return
	}

	response := models.ExpressionDetailResponse{
		Expression: models.ExpressionResponse{
			ID:     expr.ID,
			Status: expr.Status,
			Result: expr.Result,
		},
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (h *Handler) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	task, err := h.service.GetNextTask()
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No tasks available")
		return
	}

	response := models.GetTaskResponse{
		Task: &models.TaskResponse{
			ID:            task.ID,
			Arg1:          task.Arg1Value,
			Arg2:          task.Arg2Value,
			Operation:     task.Operation,
			OperationTime: task.OperationTime,
		},
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (h *Handler) SubmitTaskResultHandler(w http.ResponseWriter, r *http.Request) {
	var req models.TaskResultRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusUnprocessableEntity, "Invalid request payload")
		return
	}

	err := h.service.UpdateTaskResult(req.ID, req.Result)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Task not found")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload) 
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
	w.Write([]byte("\n"))
}
