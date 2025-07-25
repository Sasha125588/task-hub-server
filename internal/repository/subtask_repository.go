package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/Sasha125588/event_app/internal/models"
)

// SubTaskRepository handles database operations for subtasks
type SubTaskRepository struct {
	db *sql.DB
}

// NewSubTaskRepository creates a new instance of SubTaskRepository
func NewSubTaskRepository(db *sql.DB) *SubTaskRepository {
	return &SubTaskRepository{db: db}
}

// CreateSubTask creates a new subtask in the database
// The order is automatically set to be the last in the task's subtask list
func (r *SubTaskRepository) CreateSubTask(subTask *models.SubTask) error {
	// Get the maximum order for the task
	var maxOrder int
	err := r.db.QueryRow("SELECT COALESCE(MAX(\"order\"), -1) FROM sub_tasks WHERE task_id = $1", subTask.TaskID).Scan(&maxOrder)
	if err != nil {
		return err
	}

	// Set the new subtask's order to be last
	subTask.Order = maxOrder + 1

	query := `
		INSERT INTO sub_tasks (id, task_id, title, description, status, "order", created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`

	return r.db.QueryRow(
		query,
		subTask.ID,
		subTask.TaskID,
		subTask.Title,
		subTask.Description,
		subTask.Status,
		subTask.Order,
		subTask.CreatedAt,
		subTask.UpdatedAt,
	).Scan(&subTask.ID)
}

// GetSubTaskByID retrieves a subtask by its ID
func (r *SubTaskRepository) GetSubTaskByID(id string) (*models.SubTask, error) {
	query := `
		SELECT id, task_id, title, description, status, "order", created_at, updated_at
		FROM sub_tasks WHERE id = $1
	`
	fmt.Printf("GetSubTaskByID query: %s with id: %s\n", query, id)

	subTask := &models.SubTask{}
	row := r.db.QueryRow(query, id)

	err := row.Scan(
		&subTask.ID,
		&subTask.TaskID,
		&subTask.Title,
		&subTask.Description,
		&subTask.Status,
		&subTask.Order,
		&subTask.CreatedAt,
		&subTask.UpdatedAt,
	)

	if err != nil {
		fmt.Printf("GetSubTaskByID error: %v\n", err)
		return nil, err
	}

	fmt.Printf("GetSubTaskByID result: %+v\n", subTask)
	return subTask, nil
}

// UpdateSubTask updates an existing subtask with the provided changes
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

// DeleteSubTask removes a subtask from the database
func (r *SubTaskRepository) DeleteSubTask(id string) error {
	query := "DELETE FROM sub_tasks WHERE id = $1"
	_, err := r.db.Exec(query, id)
	return err
}

// GetSubTasksByTaskID retrieves all subtasks for a specific task, ordered by their order field
func (r *SubTaskRepository) GetSubTasksByTaskID(taskID string) ([]models.SubTask, error) {
	query := `
		SELECT id, task_id, title, description, status, "order", created_at, updated_at
		FROM sub_tasks
		WHERE task_id = $1
		ORDER BY "order" ASC`

	rows, err := r.db.Query(query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subTasks []models.SubTask
	for rows.Next() {
		var subTask models.SubTask
		err := rows.Scan(
			&subTask.ID,
			&subTask.TaskID,
			&subTask.Title,
			&subTask.Description,
			&subTask.Status,
			&subTask.Order,
			&subTask.CreatedAt,
			&subTask.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		subTasks = append(subTasks, subTask)
	}

	return subTasks, nil
}

// ReorderSubTask updates the order of a subtask and adjusts other subtasks' orders accordingly
func (r *SubTaskRepository) ReorderSubTask(taskID string, subTaskID string, newOrder int) error {
	fmt.Printf("ReorderSubTask called with taskID: %s, subTaskID: %s, newOrder: %d\n", taskID, subTaskID, newOrder)

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get current order of the subtask
	var currentOrder int
	err = tx.QueryRow("SELECT \"order\" FROM sub_tasks WHERE id = $1", subTaskID).Scan(&currentOrder)
	if err != nil {
		return fmt.Errorf("failed to get current order: %w", err)
	}

	// Временно установим очень большой order для нашего элемента, чтобы освободить текущую позицию
	_, err = tx.Exec("UPDATE sub_tasks SET \"order\" = -1 WHERE id = $1", subTaskID)
	if err != nil {
		return fmt.Errorf("failed to set temporary order: %w", err)
	}

	if currentOrder < newOrder {
		// Сдвигаем элементы вверх
		_, err = tx.Exec(`
            UPDATE sub_tasks 
            SET "order" = "order" - 1 
            WHERE task_id = $1 
            AND "order" > $2 
            AND "order" <= $3
            AND id != $4`,
			taskID, currentOrder, newOrder, subTaskID)
	} else {
		// Сдвигаем элементы вниз
		_, err = tx.Exec(`
            UPDATE sub_tasks 
            SET "order" = "order" + 1 
            WHERE task_id = $1 
            AND "order" >= $2 
            AND "order" < $3
            AND id != $4`,
			taskID, newOrder, currentOrder, subTaskID)
	}
	if err != nil {
		return fmt.Errorf("failed to update other tasks order: %w", err)
	}

	// Устанавливаем финальный порядок для нашего элемента
	_, err = tx.Exec("UPDATE sub_tasks SET \"order\" = $1 WHERE id = $2", newOrder, subTaskID)
	if err != nil {
		return fmt.Errorf("failed to set final order: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
