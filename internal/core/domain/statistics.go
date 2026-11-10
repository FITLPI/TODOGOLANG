package domain

import (
	"time"

	core_errors "github.com/FITLPI/TODOGOLANG/internal/core/errors"
)

// Statistics — агрегат по задачам за период. Считается в домене
// (CreateStatistics), а не в SQL.
type Statistics struct {
	TasksCreated               int
	TasksCompleted             int
	TasksCompletedRate         *float64       // nil, если создано 0 задач
	TasksAverageCompletionTime *time.Duration // nil, если выполнено 0 задач
}

func CreateStatistics(tasks []Task) Statistics {
	var (
		tasksCreated   int
		tasksCompleted int
		totalDuration  time.Duration
	)

	for _, task := range tasks {
		tasksCreated++

		if task.Completed && task.CompletedAt != nil {
			tasksCompleted++

			duration := task.CompletedAt.Sub(task.CreatedAt)
			if duration > 0 {
				totalDuration += duration
			}
		}
	}

	stats := Statistics{
		TasksCreated:   tasksCreated,
		TasksCompleted: tasksCompleted,
	}

	if tasksCreated > 0 {
		rate := float64(tasksCompleted) / float64(tasksCreated) * 100
		stats.TasksCompletedRate = &rate
	}

	if tasksCompleted > 0 && totalDuration > 0 {
		avg := totalDuration / time.Duration(tasksCompleted)
		stats.TasksAverageCompletionTime = &avg
	}

	return stats
}

// ValidateRange проверяет смысл диапазона: если заданы обе границы,
// to должен быть строго после from. Формат дат проверяет транспорт.
func ValidateRange(from, to *time.Time) error {
	if from != nil && to != nil {
		if to.Before(*from) || to.Equal(*from) {
			return core_errors.ErrInvalidArgument
		}
	}

	return nil
}
