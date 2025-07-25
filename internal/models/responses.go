package models

type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}

type MessageResponse struct {
	Message string `json:"message" example:"operation completed successfully"`
}

type TasksResponse struct {
	Tasks []Task `json:"tasks"`
}

// SubTasksResponse represents the response body for getting subtasks
// @Description Response body containing a list of subtasks
type SubTasksResponse struct {
	SubTasks []SubTask `json:"subtasks"`
}
