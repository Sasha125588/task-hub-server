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

type SubTasksResponse struct {
	SubTasks []SubTask `json:"subtasks"`
}
