package tasks_service

import (
	"context"
	"fmt"

	"github.com/FITLPI/TODOGOLANG/internal/core/domain"
)

func (s *TasksService) UpdateTask(ctx context.Context, task domain.Task) (domain.Task, error) {
	if err := task.Validate(); err != nil {
		return domain.Task{}, fmt.Errorf("validate task: %w", err)
	}

	task, err := s.tasksRepository.UpdateTask(ctx, task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("update task: %w", err)
	}

	return task, nil
}

func (s *TasksService) DeleteTask(ctx context.Context, id int) error {
	err := s.tasksRepository.DeleteTask(ctx, id)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	return nil
}