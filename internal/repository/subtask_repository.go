package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/Sasha125588/event_app/internal/models"
)

type SubTaskRepository struct {
	db *sql.DB
}

func NewSubTaskRepository(db *sql.DB) *SubTaskRepository {
	return &SubTaskRepository{db: db}
}

func (r *SubTaskRepository) CreateSubTask(subTask *models.SubTask) error {
	query := `
		INSERT INTO sub_tasks (id, task_id, title, description, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(query, subTask.ID, subTask.TaskID, subTask.Title,
		subTask.Description, subTask.Status, subTask.CreatedAt, subTask.UpdatedAt)
	return err
}

func (r *SubTaskRepository) GetSubTaskByID(id string) (*models.SubTask, error) {
	query := `
		SELECT id, task_id, title, description, status, created_at, updated_at
		FROM sub_tasks WHERE id = $1
	`
	subTask := &models.SubTask{}
	row := r.db.QueryRow(query, id)

	err := row.Scan(
		&subTask.ID, &subTask.TaskID, &subTask.Title, &subTask.Description,
		&subTask.Status, &subTask.CreatedAt, &subTask.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return subTask, nil
}

func (r *SubTaskRepository) UpdateSubTask(id string, updates *models.UpdateSubTaskRequest) error {
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if updates.Title != nil {
		setParts = append(setParts, fmt.Sprintf("title = $%d", argIndex))
		args = append(args, *updates.Title)
		argIndex++
	}
	if updates.Description != nil {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, *updates.Description)
		argIndex++
	}
	if updates.Status != nil {
		setParts = append(setParts, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *updates.Status)
		argIndex++
	}

	if len(setParts) == 0 {
		return fmt.Errorf("no fields to update")
	}

	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	args = append(args, id)

	query := fmt.Sprintf("UPDATE sub_tasks SET %s WHERE id = $%d", strings.Join(setParts, ", "), argIndex)
	_, err := r.db.Exec(query, args...)
	return err
}

func (r *SubTaskRepository) DeleteSubTask(id string) error {
	query := "DELETE FROM sub_tasks WHERE id = $1"
	_, err := r.db.Exec(query, id)
	return err
}

func (r *SubTaskRepository) GetSubTasksByTaskID(taskID string) ([]models.SubTask, error) {
	query := `
		SELECT id, task_id, title, description, status, created_at, updated_at
		FROM sub_tasks WHERE task_id = $1 ORDER BY created_at ASC
	`
	rows, err := r.db.Query(query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subTasks []models.SubTask
	for rows.Next() {
		var subTask models.SubTask
		err := rows.Scan(
			&subTask.ID, &subTask.TaskID, &subTask.Title, &subTask.Description,
			&subTask.Status, &subTask.CreatedAt, &subTask.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		subTasks = append(subTasks, subTask)
	}

	return subTasks, nil
}
