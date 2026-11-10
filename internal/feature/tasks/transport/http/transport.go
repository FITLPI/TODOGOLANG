package tasks_transport_http

import (
	"context"
	"net/http"

	"github.com/FITLPI/TODOGOLANG/internal/core/domain"
	core_http_server "github.com/FITLPI/TODOGOLANG/internal/core/transport/http/server"
)

type TasksHTTPHandler struct {
	tasksService TasksService
}

type TasksService interface {
	CreateTask(ctx context.Context, task domain.Task) (domain.Task, error)
	GetTaskByID(ctx context.Context, id int) (domain.Task, error)
	ListAllTasks(ctx context.Context, userID int) ([]domain.Task, error)
	ListIncompleteTasks(ctx context.Context, userID int) ([]domain.Task, error)
	ListCompletedTasks(ctx context.Context, userID int) ([]domain.Task, error)
	UpdateTask(ctx context.Context, task domain.Task) (domain.Task, error)
	DeleteTask(ctx context.Context, id int) error
}

func NewTasksHTTPHandler(tasksService TasksService) *TasksHTTPHandler {
	return &TasksHTTPHandler{
		tasksService: tasksService,
	}
}

func (h *TasksHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/tasks",
			Handler: h.CreateTask,
		},
		{
			Method:  http.MethodGet,
			Path:    "/tasks",
			Handler: h.ListTasks,
		},
		{
			Method:  http.MethodGet,
			Path:    "/tasks/{id}",
			Handler: h.GetTaskByID,
		},
		{
			Method:  http.MethodPut,
			Path:    "/tasks/{id}",
			Handler: h.UpdateTask,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/tasks/{id}",
			Handler: h.DeleteTask,
		},
	}
}