package service

import (
	"database/sql"
	"fmt"

	"github.com/Sasha125588/event_app/internal/models"
	"github.com/Sasha125588/event_app/internal/repository"
)

// TaskService handles business logic for tasks and subtasks
type TaskService struct {
	taskRepo    *repository.TaskRepository
	subTaskRepo *repository.SubTaskRepository
}

// NewTaskService creates a new instance of TaskService
func NewTaskService(taskRepo *repository.TaskRepository, subTaskRepo *repository.SubTaskRepository) *TaskService {
	return &TaskService{
		taskRepo:    taskRepo,
		subTaskRepo: subTaskRepo,
	}
}

func (s *TaskService) CreateTask(req models.CreateTaskRequest) (*models.Task, error) {

	task := models.NewTask(req)
	err := s.taskRepo.CreateTask(task)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return s.taskRepo.GetTaskByID(task.ID)
}

func (s *TaskService) GetTask(id string) (*models.Task, error) {
	return s.taskRepo.GetTaskByID(id)
}

func (s *TaskService) UpdateTask(id string, req models.UpdateTaskRequest) (*models.Task, error) {
	_, err := s.taskRepo.GetTaskByID(id)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	err = s.taskRepo.UpdateTask(id, &req)
	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	return s.taskRepo.GetTaskByID(id)
}

func (s *TaskService) DeleteTask(id string) error {
	_, err := s.taskRepo.GetTaskByID(id)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}

	return s.taskRepo.DeleteTask(id)
}

func (s *TaskService) GetTasks(filters models.TaskFilters) ([]models.Task, error) {
	return s.taskRepo.GetTasks(filters)
}

// CreateSubTask creates a new subtask for a specific task
// It validates that the parent task exists before creating the subtask
func (s *TaskService) CreateSubTask(taskID string, req models.CreateSubTaskRequest) (*models.SubTask, error) {
	_, err := s.taskRepo.GetTaskByID(taskID)
	if err != nil {
		return nil, fmt.Errorf("parent task not found: %w", err)
	}

	subTask := models.NewSubTask(taskID, req)
	err = s.subTaskRepo.CreateSubTask(subTask)
	if err != nil {
		return nil, fmt.Errorf("failed to create subtask: %w", err)
	}

	return s.subTaskRepo.GetSubTaskByID(subTask.ID)
}

// UpdateSubTask updates an existing subtask
// It validates that the subtask exists before updating it
func (s *TaskService) UpdateSubTask(id string, req models.UpdateSubTaskRequest) (*models.SubTask, error) {
	_, err := s.subTaskRepo.GetSubTaskByID(id)
	if err != nil {
		return nil, fmt.Errorf("subtask not found: %w", err)
	}

	err = s.subTaskRepo.UpdateSubTask(id, &req)
	if err != nil {
		return nil, fmt.Errorf("failed to update subtask: %w", err)
	}

	return s.subTaskRepo.GetSubTaskByID(id)
}

// DeleteSubTask removes a subtask from the database
// It validates that the subtask exists before deleting it
func (s *TaskService) DeleteSubTask(id string) error {
	_, err := s.subTaskRepo.GetSubTaskByID(id)
	if err != nil {
		return fmt.Errorf("subtask not found: %w", err)
	}

	return s.subTaskRepo.DeleteSubTask(id)
}

// GetSubTasksByTaskID retrieves all subtasks for a specific task
func (s *TaskService) GetSubTasksByTaskID(taskID string) ([]models.SubTask, error) {
	return s.subTaskRepo.GetSubTasksByTaskID(taskID)
}

// ReorderSubTask reorders a subtask within its parent task
// It validates that the subtask belongs to the specified task before reordering
func (s *TaskService) ReorderSubTask(taskID string, subTaskID string, newOrder int) error {
	// Verify that the task exists
	_, err := s.taskRepo.GetTaskByID(taskID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("task not found")
		}
		return err
	}

	// Verify that the subtask exists and belongs to the task
	subTask, err := s.subTaskRepo.GetSubTaskByID(subTaskID)
	if err != nil {
		if err == sql.ErrNoRows {
			return err
		}
		return err
	}
	if subTask.TaskID != taskID {
		return fmt.Errorf("subtask does not belong to the specified task")
	}

	// Get all subtasks to validate the new order
	subtasks, err := s.subTaskRepo.GetSubTasksByTaskID(taskID)
	if err != nil {
		return err
	}

	// Validate that the new order is within bounds
	if newOrder < 0 || newOrder >= len(subtasks) {
		return fmt.Errorf("invalid order: must be between 0 and %d", len(subtasks)-1)
	}

	return s.subTaskRepo.ReorderSubTask(taskID, subTaskID, newOrder)
}
