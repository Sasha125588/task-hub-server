package service

import (
	"fmt"

	"github.com/Sasha125588/event_app/internal/models"
	"github.com/Sasha125588/event_app/internal/repository"
)

type TaskService struct {
	taskRepo    *repository.TaskRepository
	subTaskRepo *repository.SubTaskRepository
}

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

func (s *TaskService) DeleteSubTask(id string) error {
	_, err := s.subTaskRepo.GetSubTaskByID(id)
	if err != nil {
		return fmt.Errorf("subtask not found: %w", err)
	}

	return s.subTaskRepo.DeleteSubTask(id)
}

func (s *TaskService) GetSubTasksByTaskID(taskID string) ([]models.SubTask, error) {
	return s.subTaskRepo.GetSubTasksByTaskID(taskID)
}
