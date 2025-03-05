package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/popvictor123/distributed-calc/internal/orchestrator/models"
	"github.com/google/uuid"
)

type Agent struct {
	OrchestratorURL string
	WorkerCount     int
	Client          *http.Client
}

type TaskResult struct {
	ID     string  `json:"id"`
	Result float64 `json:"result"`
}

func NewAgent() *Agent {
	orchestratorURL := getEnvOrDefault("ORCHESTRATOR_URL", "http://localhost:8080")
	workerCountStr := getEnvOrDefault("COMPUTING_POWER", "3")
	workerCount, err := strconv.Atoi(workerCountStr)
	if err != nil || workerCount < 1 {
		workerCount = 3
	}

	return &Agent{
		OrchestratorURL: orchestratorURL,
		WorkerCount:     workerCount,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func (a *Agent) Start() {
	log.Printf("Starting agent with %d workers, connecting to orchestrator at %s", a.WorkerCount, a.OrchestratorURL)
	
	for i := 0; i < a.WorkerCount; i++ {
		go a.worker(i)
	}

	select {}
}

func (a *Agent) worker(id int) {
	log.Printf("Worker %d started", id)
	
	for {
		task, err := a.fetchTask()
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		log.Printf("Worker %d received task %s: %v %v %v", id, task.ID, task.Arg1, task.Operation, task.Arg2)

		result, err := a.processTask(task)
		if err != nil {
			log.Printf("Worker %d error processing task %s: %v", id, task.ID, err)
			continue
		}

		err = a.submitResult(task.ID, result)
		if err != nil {
			log.Printf("Worker %d error submitting result for task %s: %v", id, task.ID, err)
			continue
		}

		log.Printf("Worker %d completed task %s with result %v", id, task.ID, result)
	}
}

func (a *Agent) fetchTask() (*models.TaskResponse, error) {
	resp, err := a.Client.Get(fmt.Sprintf("%s/internal/task", a.OrchestratorURL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("no tasks available")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response models.GetTaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Task, nil
}

func (a *Agent) processTask(task *models.TaskResponse) (float64, error) {
	var result float64

	time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)

	switch task.Operation {
	case models.OperationAddition:
		result = task.Arg1 + task.Arg2
	case models.OperationSubtraction:
		result = task.Arg1 - task.Arg2
	case models.OperationMultiplication:
		result = task.Arg1 * task.Arg2
	case models.OperationDivision:
		if task.Arg2 == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		result = task.Arg1 / task.Arg2
	case models.OperationValue:
		// Just return the value
		result = task.Arg1
	default:
		return 0, fmt.Errorf("unknown operation: %s", task.Operation)
	}

	return result, nil
}

func (a *Agent) submitResult(taskID uuid.UUID, result float64) error {
	taskResult := models.TaskResultRequest{
		ID:     taskID,
		Result: result,
	}

	jsonData, err := json.Marshal(taskResult)
	if err != nil {
		return err
	}

	resp, err := a.Client.Post(
		fmt.Sprintf("%s/internal/task", a.OrchestratorURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to submit result: status code %d", resp.StatusCode)
	}

	return nil
}
