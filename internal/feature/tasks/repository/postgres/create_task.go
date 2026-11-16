package tasks_postgres_repository

import (
	"context"

	"github.com/FITLPI/TODOGOLANG/internal/core/domain"
)

func (r *TaskRepository) CreateTask(
	ctx context.Context,
	task domain.Task,
) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
		INSERT INTO todoapp.tasks (title, description, completed, created_at, completed_at, author_user_id, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, version, title, description, completed, created_at, completed_at, author_user_id
	`

	var createdTask domain.Task
	err := r.pool.QueryRow(ctx, query,
		task.Title,
		task.Description,
		task.Completed,
		task.CreatedAt,
		task.CompletedAt,
		task.AuthorUserID,
		task.Version,
	).Scan(
		&createdTask.ID,
		&createdTask.Version,
		&createdTask.Title,
		&createdTask.Description,
		&createdTask.Completed,
		&createdTask.CreatedAt,
		&createdTask.CompletedAt,
		&createdTask.AuthorUserID,
	)

	if err != nil {
		return domain.Task{}, err
	}

	return createdTask, nil
}
