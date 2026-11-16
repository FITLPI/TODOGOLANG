package tasks_postgres_repository

import (
	"context"

	"github.com/FITLPI/TODOGOLANG/internal/core/domain"
	core_errors "github.com/FITLPI/TODOGOLANG/internal/core/errors"
)

func (r *TaskRepository) UpdateTask(
	ctx context.Context,
	task domain.Task,
) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
		UPDATE todoapp.tasks
		SET title = $2, description = $3, completed = $4, completed_at = $5, version = version + 1
		WHERE id = $1 AND version = $6
		RETURNING id, version, title, description, completed, created_at, completed_at, author_user_id
	`

	var updatedTask domain.Task
	err := r.pool.QueryRow(ctx, query,
		task.ID,
		task.Title,
		task.Description,
		task.Completed,
		task.CompletedAt,
		task.Version,
	).Scan(
		&updatedTask.ID,
		&updatedTask.Version,
		&updatedTask.Title,
		&updatedTask.Description,
		&updatedTask.Completed,
		&updatedTask.CreatedAt,
		&updatedTask.CompletedAt,
		&updatedTask.AuthorUserID,
	)

	if err != nil {
		return domain.Task{}, core_errors.ErrConflict
	}

	return updatedTask, nil
}

func (r *TaskRepository) DeleteTask(
	ctx context.Context,
	id int,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `DELETE FROM todoapp.tasks WHERE id = $1`

	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return core_errors.ErrNotFound
	}

	return nil
}