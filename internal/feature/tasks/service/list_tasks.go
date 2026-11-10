package tasks_service

import (
	"context"
	"fmt"

	"github.com/FITLPI/TODOGOLANG/internal/core/domain"
)

func (s *TasksService) ListAllTasks(ctx context.Context, userID int) ([]domain.Task, error) {
	tasks, err := s.tasksRepository.ListAllTasks(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list all tasks: %w", err)
	}
	return tasks, nil
}

func (s *TasksService) ListIncompleteTasks(ctx context.Context, userID int) ([]domain.Task, error) {
	tasks, err := s.tasksRepository.ListIncompleteTasks(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list incomplete tasks: %w", err)
	}
	return tasks, nil
}

func (s *TasksService) ListCompletedTasks(ctx context.Context, userID int) ([]domain.Task, error) {
	tasks, err := s.tasksRepository.ListCompletedTasks(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list completed tasks: %w", err)
	}
	return tasks, nil
}