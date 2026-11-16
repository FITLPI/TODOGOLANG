package tasks_service

import (
	"context"

	"github.com/FITLPI/TODOGOLANG/internal/core/domain"
)

type TasksService struct {
	tasksRepository TasksRepository
}

type TasksRepository interface {
	CreateTask(ctx context.Context, task domain.Task) (domain.Task, error)
	GetTaskByID(ctx context.Context, id int) (domain.Task, error)
	ListAllTasks(ctx context.Context, userID int) ([]domain.Task, error)
	ListIncompleteTasks(ctx context.Context, userID int) ([]domain.Task, error)
	ListCompletedTasks(ctx context.Context, userID int) ([]domain.Task, error)
	UpdateTask(ctx context.Context, task domain.Task) (domain.Task, error)
	DeleteTask(ctx context.Context, id int) error
}

func NewTasksService(tasksRepository TasksRepository) *TasksService {
	return &TasksService{
		tasksRepository: tasksRepository,
	}
}