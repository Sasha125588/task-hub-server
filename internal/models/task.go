package models

import (
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	StatusNotStarted TaskStatus = "not-started"
	StatusCompleted  TaskStatus = "completed"
	StatusInProgress TaskStatus = "in-progress"
)

type User struct {
	ID   string `json:"id" db:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name string `json:"name" db:"name" example:"John Doe"`
	Src  string `json:"src" db:"src" example:"https://avatars.githubusercontent.com/u/124599?v=4"`
}

// SubTask represents a subtask within a task
// @Description A subtask that belongs to a parent task
type SubTask struct {
	ID          string     `json:"id" db:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	TaskID      string     `json:"task_id" db:"task_id" example:"123e4567-e89b-12d3-a456-426614174001"`
	Title       string     `json:"title" db:"title" example:"Implement user authentication"`
	Description *string    `json:"description,omitempty" db:"description" example:"Add JWT token authentication"`
	Status      TaskStatus `json:"status" db:"status" example:"not-started"`
	Order       int        `json:"order" db:"order" example:"1"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at" example:"2024-01-01T00:00:00Z"`
}

type Task struct {
	ID          string     `json:"id" db:"id"`
	Title       string     `json:"title" db:"title"`
	IconName    string     `json:"icon_name" db:"icon_name"`
	StartTime   *string    `json:"start_time,omitempty" db:"start_time"`
	EndTime     *string    `json:"end_time,omitempty" db:"end_time"`
	DueDate     time.Time  `json:"due_date" db:"due_date"`
	Progress    int        `json:"progress" db:"progress"`
	Status      TaskStatus `json:"status" db:"status"`
	Comments    int        `json:"comments" db:"comments"`
	Attachments int        `json:"attachments" db:"attachments"`
	Links       int        `json:"links" db:"links"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`

	Users    []User    `json:"users,omitempty"`
	SubTasks []SubTask `json:"sub_tasks,omitempty"`
}

type CreateTaskRequest struct {
	Title     string     `json:"title" binding:"required"`
	IconName  string     `json:"icon_name" binding:"required"`
	StartTime *string    `json:"start_time,omitempty"`
	EndTime   *string    `json:"end_time,omitempty"`
	DueDate   time.Time  `json:"due_date" binding:"required"`
	Status    TaskStatus `json:"status" binding:"required"`
	UserIDs   []string   `json:"user_ids,omitempty"`
}

type UpdateTaskRequest struct {
	Title       *string     `json:"title,omitempty"`
	IconName    *string     `json:"icon_name,omitempty"`
	StartTime   *string     `json:"start_time,omitempty"`
	EndTime     *string     `json:"end_time,omitempty"`
	DueDate     *time.Time  `json:"due_date,omitempty"`
	Progress    *int        `json:"progress,omitempty"`
	Status      *TaskStatus `json:"status,omitempty"`
	Comments    *int        `json:"comments,omitempty"`
	Attachments *int        `json:"attachments,omitempty"`
	Links       *int        `json:"links,omitempty"`
}

// CreateSubTaskRequest represents the request body for creating a new subtask
// @Description Request body for creating a new subtask
type CreateSubTaskRequest struct {
	Title       string     `json:"title" binding:"required" example:"Implement user authentication"`
	Description *string    `json:"description,omitempty" example:"Add JWT token authentication"`
	Status      TaskStatus `json:"status" binding:"required" example:"not-started"`
}

// UpdateSubTaskRequest represents the request body for updating an existing subtask
// @Description Request body for updating an existing subtask
type UpdateSubTaskRequest struct {
	Title       *string     `json:"title,omitempty" example:"Implement user authentication"`
	Description *string     `json:"description,omitempty" example:"Add JWT token authentication"`
	Status      *TaskStatus `json:"status,omitempty" example:"in-progress"`
}

type TaskFilters struct {
	Status   TaskStatus `form:"status"`
	SortBy   string     `form:"sort_by"`
	SortType string     `form:"sort_type"`
	Limit    int        `form:"limit"`
	Offset   int        `form:"offset"`
}

func NewTask(req CreateTaskRequest) *Task {
	now := time.Now()
	return &Task{
		ID:          uuid.New().String(),
		Title:       req.Title,
		IconName:    req.IconName,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		DueDate:     req.DueDate,
		Progress:    0,
		Status:      req.Status,
		Comments:    0,
		Attachments: 0,
		Links:       0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func NewSubTask(taskID string, req CreateSubTaskRequest) *SubTask {
	now := time.Now()
	return &SubTask{
		ID:          uuid.New().String(),
		TaskID:      taskID,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
