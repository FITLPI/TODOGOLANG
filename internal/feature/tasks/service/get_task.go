package tasks_service

import (
	"context"
	"fmt"

	"github.com/FITLPI/TODOGOLANG/internal/core/domain"
)

func (s *TasksService) GetTaskByID(ctx context.Context, id int) (domain.Task, error) {
	task, err := s.tasksRepository.GetTaskByID(ctx, id)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get task by id: %w", err)
	}
	return task, nil
}