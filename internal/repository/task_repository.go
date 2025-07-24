package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/Sasha125588/event_app/internal/models"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) CreateTask(task *models.Task) error {
	query := `
		INSERT INTO tasks (id, title, icon_name, start_time, end_time, due_date, progress, status, comments, attachments, links, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	_, err := r.db.Exec(query, task.ID, task.Title, task.IconName, task.StartTime, task.EndTime,
		task.DueDate, task.Progress, task.Status, task.Comments, task.Attachments, task.Links,
		task.CreatedAt, task.UpdatedAt)
	return err
}

func (r *TaskRepository) GetTaskByID(id string) (*models.Task, error) {
	query := `
		SELECT id, title, icon_name, start_time, end_time, due_date, progress, status, comments, attachments, links, created_at, updated_at
		FROM tasks WHERE id = $1
	`
	task := &models.Task{}
	row := r.db.QueryRow(query, id)

	err := row.Scan(
		&task.ID, &task.Title, &task.IconName, &task.StartTime, &task.EndTime,
		&task.DueDate, &task.Progress, &task.Status, &task.Comments, &task.Attachments,
		&task.Links, &task.CreatedAt, &task.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	subTasks, err := r.GetTaskSubTasks(id)
	if err != nil {
		return nil, err
	}
	task.SubTasks = subTasks

	return task, nil
}

func (r *TaskRepository) UpdateTask(id string, updates *models.UpdateTaskRequest) error {
	setParts := []string{}
	args := []any{}
	argIndex := 1

	if updates.Title != nil {
		setParts = append(setParts, fmt.Sprintf("title = $%d", argIndex))
		args = append(args, *updates.Title)
		argIndex++
	}
	if updates.IconName != nil {
		setParts = append(setParts, fmt.Sprintf("icon_name = $%d", argIndex))
		args = append(args, *updates.IconName)
		argIndex++
	}
	if updates.StartTime != nil {
		setParts = append(setParts, fmt.Sprintf("start_time = $%d", argIndex))
		args = append(args, *updates.StartTime)
		argIndex++
	}
	if updates.EndTime != nil {
		setParts = append(setParts, fmt.Sprintf("end_time = $%d", argIndex))
		args = append(args, *updates.EndTime)
		argIndex++
	}
	if updates.DueDate != nil {
		setParts = append(setParts, fmt.Sprintf("due_date = $%d", argIndex))
		args = append(args, *updates.DueDate)
		argIndex++
	}
	if updates.Progress != nil {
		setParts = append(setParts, fmt.Sprintf("progress = $%d", argIndex))
		args = append(args, *updates.Progress)
		argIndex++
	}
	if updates.Status != nil {
		setParts = append(setParts, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *updates.Status)
		argIndex++
	}
	if updates.Comments != nil {
		setParts = append(setParts, fmt.Sprintf("comments = $%d", argIndex))
		args = append(args, *updates.Comments)
		argIndex++
	}
	if updates.Attachments != nil {
		setParts = append(setParts, fmt.Sprintf("attachments = $%d", argIndex))
		args = append(args, *updates.Attachments)
		argIndex++
	}
	if updates.Links != nil {
		setParts = append(setParts, fmt.Sprintf("links = $%d", argIndex))
		args = append(args, *updates.Links)
		argIndex++
	}

	if len(setParts) == 0 {
		return fmt.Errorf("no fields to update")
	}

	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	args = append(args, id)

	query := fmt.Sprintf("UPDATE tasks SET %s WHERE id = $%d", strings.Join(setParts, ", "), argIndex)
	_, err := r.db.Exec(query, args...)
	return err
}

func (r *TaskRepository) DeleteTask(id string) error {
	query := "DELETE FROM tasks WHERE id = $1"
	_, err := r.db.Exec(query, id)
	return err
}

func (r *TaskRepository) GetTasks(filters models.TaskFilters) ([]models.Task, error) {
	query := "SELECT id, title, icon_name, start_time, end_time, due_date, progress, status, comments, attachments, links, created_at, updated_at FROM tasks"
	args := []any{}
	whereConditions := []string{}
	argIndex := 1

	if filters.Status != "" && filters.Status != "all" {
		whereConditions = append(whereConditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, filters.Status)
		argIndex++
	}

	if len(whereConditions) > 0 {
		query += " WHERE " + strings.Join(whereConditions, " AND ")
	}

	if filters.SortBy != "" {
		sortColumn := "due_date"
		sortDirection := "ASC"

		if filters.SortType == "desc" {
			sortDirection = "DESC"
		}

		query += fmt.Sprintf(" ORDER BY %s %s", sortColumn, sortDirection)
	} else {
		query += " ORDER BY created_at DESC"
	}

	if filters.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filters.Limit)
		argIndex++

		if filters.Offset > 0 {
			query += fmt.Sprintf(" OFFSET $%d", argIndex)
			args = append(args, filters.Offset)
		}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.ID, &task.Title, &task.IconName, &task.StartTime, &task.EndTime,
			&task.DueDate, &task.Progress, &task.Status, &task.Comments, &task.Attachments,
			&task.Links, &task.CreatedAt, &task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		subTasks, err := r.GetTaskSubTasks(task.ID)
		if err != nil {
			return nil, err
		}
		task.SubTasks = subTasks

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *TaskRepository) GetTaskSubTasks(taskID string) ([]models.SubTask, error) {
	query := `
		SELECT id, task_id, title, description, status, created_at, updated_at 
		FROM sub_tasks 
		WHERE task_id = $1 
		ORDER BY created_at ASC
	`
	rows, err := r.db.Query(query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subTasks []models.SubTask
	for rows.Next() {
		var subTask models.SubTask
		err := rows.Scan(&subTask.ID, &subTask.TaskID, &subTask.Title,
			&subTask.Description, &subTask.Status, &subTask.CreatedAt, &subTask.UpdatedAt)
		if err != nil {
			return nil, err
		}
		subTasks = append(subTasks, subTask)
	}

	return subTasks, nil
}
