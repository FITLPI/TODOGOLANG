package tasks_postgres_repository

import (
	"context"

	"github.com/FITLPI/TODOGOLANG/internal/core/domain"
	core_errors "github.com/FITLPI/TODOGOLANG/internal/core/errors"
)

func (r *TaskRepository) GetTaskByID(
	ctx context.Context,
	id int,
) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
		SELECT id, version, title, description, completed, created_at, completed_at, author_user_id
		FROM todoapp.tasks
		WHERE id = $1
	`

	var task domain.Task
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&task.ID,
		&task.Version,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
		&task.CompletedAt,
		&task.AuthorUserID,
	)

	if err != nil {
		return domain.Task{}, core_errors.ErrNotFound
	}

	return task, nil
}
