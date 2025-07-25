package handlers

import (
	"database/sql"
	"fmt"
	"io"
	"strings"

	"bytes"

	"net/http"

	"github.com/Sasha125588/event_app/internal/models"
	"github.com/Sasha125588/event_app/internal/service"
	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	taskService *service.TaskService
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

// CreateTask handles POST /api/v1/tasks
// @Summary Create a new task
// @Description Create a new task with the provided details
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body models.CreateTaskRequest true "Task details"
// @Success 201 {object} models.Task
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /tasks [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req models.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.taskService.CreateTask(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// GetTasks handles GET /api/v1/tasks
// @Summary Get all tasks
// @Description Get all tasks with optional filtering
// @Tags tasks
// @Accept json
// @Produce json
// @Param status query string false "Filter by status"
// @Param priority query string false "Filter by priority"
// @Param limit query int false "Limit number of tasks returned (default: 50)"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} models.TasksResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /tasks [get]
func (h *TaskHandler) GetTasks(c *gin.Context) {
	var filters models.TaskFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default values if not provided
	if filters.Limit == 0 {
		filters.Limit = 50 // default limit
	}

	tasks, err := h.taskService.GetTasks(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

// GetTask handles GET /api/v1/tasks/:id
// @Summary Get a task by ID
// @Description Get detailed information about a specific task
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} models.Task
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /tasks/{id} [get]
func (h *TaskHandler) GetTask(c *gin.Context) {
	id := c.Param("id")

	task, err := h.taskService.GetTask(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

// UpdateTask handles PUT /api/v1/tasks/:id
// @Summary Update a task
// @Description Update an existing task's details
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Param task body models.UpdateTaskRequest true "Updated task details"
// @Success 200 {object} models.Task
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /tasks/{id} [put]
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.taskService.UpdateTask(id, req)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

// DeleteTask handles DELETE /api/v1/tasks/:id
// @Summary Delete a task
// @Description Delete an existing task
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} models.MessageResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id := c.Param("id")

	err := h.taskService.DeleteTask(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

// CreateSubTask handles POST /api/v1/tasks/:id/subtasks
// @Summary Create a subtask
// @Description Create a new subtask for a specific task
// @Tags subtasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Param subtask body models.CreateSubTaskRequest true "Subtask details"
// @Success 201 {object} models.SubTask
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /tasks/{id}/subtasks [post]
func (h *TaskHandler) CreateSubTask(c *gin.Context) {
	taskID := c.Param("id")

	var req models.CreateSubTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subTask, err := h.taskService.CreateSubTask(taskID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, subTask)
}

// GetSubTasksByTaskID handles GET /api/v1/tasks/:id/subtasks
// @Summary Get task's subtasks
// @Description Get all subtasks for a specific task
// @Tags subtasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} models.SubTasksResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /tasks/{id}/subtasks [get]
func (h *TaskHandler) GetSubTasksByTaskID(c *gin.Context) {
	taskID := c.Param("id")

	subTasks, err := h.taskService.GetSubTasksByTaskID(taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"subtasks": subTasks})
}

// UpdateSubTask handles PUT /api/v1/tasks/:id/subtasks/:subtask_id
// @Summary Update a subtask
// @Description Update an existing subtask's details
// @Tags subtasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Param subtask_id path string true "Subtask ID"
// @Param subtask body models.UpdateSubTaskRequest true "Updated subtask details"
// @Success 200 {object} models.SubTask
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /tasks/{id}/subtasks/{subtask_id} [put]
func (h *TaskHandler) UpdateSubTask(c *gin.Context) {
	taskID := c.Param("id")
	subtaskID := c.Param("subtask_id")

	var req models.UpdateSubTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify that the subtask belongs to the task
	subTask, err := h.taskService.GetSubTasksByTaskID(taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	found := false
	for _, st := range subTask {
		if st.ID == subtaskID {
			found = true
			break
		}
	}
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "SubTask not found in the specified task"})
		return
	}

	updatedSubTask, err := h.taskService.UpdateSubTask(subtaskID, req)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "SubTask not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedSubTask)
}

// DeleteSubTask handles DELETE /api/v1/tasks/:id/subtasks/:subtask_id
// @Summary Delete a subtask
// @Description Delete an existing subtask
// @Tags subtasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Param subtask_id path string true "Subtask ID"
// @Success 200 {object} models.MessageResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /tasks/{id}/subtasks/{subtask_id} [delete]
func (h *TaskHandler) DeleteSubTask(c *gin.Context) {
	taskID := c.Param("id")
	subtaskID := c.Param("subtask_id")

	// Verify that the subtask belongs to the task
	subTasks, err := h.taskService.GetSubTasksByTaskID(taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	found := false
	for _, st := range subTasks {
		if st.ID == subtaskID {
			found = true
			break
		}
	}
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "SubTask not found in the specified task"})
		return
	}

	err = h.taskService.DeleteSubTask(subtaskID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "SubTask not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SubTask deleted successfully"})
}

// ReorderSubTaskRequest represents the request body for reordering a subtask
// @Description Request body for reordering a subtask within its parent task
type ReorderSubTaskRequest struct {
	NewOrder int `json:"new_order" binding:"required,min=0" example:"3"`
}

// ReorderSubTask handles POST /api/v1/tasks/:id/subtasks/:subtask_id/reorder
// @Summary Reorder a subtask
// @Description Reorder a subtask within its parent task by changing its order position
// @Tags subtasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Param subtask_id path string true "Subtask ID"
// @Param request body ReorderSubTaskRequest true "Reorder request"
// @Success 200 {object} models.MessageResponse
// @Failure 400 {object} models.ErrorResponse "Invalid request body or order out of bounds"
// @Failure 404 {object} models.ErrorResponse "Task or subtask not found, or subtask does not belong to task"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /tasks/{id}/subtasks/{subtask_id}/reorder [post]
func (h *TaskHandler) ReorderSubTask(c *gin.Context) {
	taskID := c.Param("id")
	subTaskID := c.Param("subtask_id")

	fmt.Printf("ReorderSubTask request received - taskID: %s, subTaskID: %s\n", taskID, subTaskID)

	// Логируем тело запроса
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Printf("Error reading request body: %v\n", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Failed to read request body"})
		return
	}
	// Восстанавливаем тело запроса для последующего использования
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	fmt.Printf("Request body: %s\n", string(body))

	var req ReorderSubTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("Error binding JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: fmt.Sprintf("Invalid request: %v", err)})
		return
	}

	fmt.Printf("Parsed request: %+v\n", req)

	err = h.taskService.ReorderSubTask(taskID, subTaskID, req.NewOrder)
	if err != nil {
		fmt.Printf("Error in ReorderSubTask service: %v\n", err)
		switch err.Error() {
		case "subtask does not belong to the specified task":
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
			return
		case "task not found":
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
			return
		default:
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "SubTask not found"})
				return
			}
			if strings.HasPrefix(err.Error(), "invalid order:") {
				c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
			return
		}
	}

	fmt.Printf("ReorderSubTask completed successfully\n")
	c.JSON(http.StatusOK, models.MessageResponse{Message: "Subtask reordered successfully"})
}
