package domain

import (
	"fmt"
	"time"

	core_errors "github.com/FITLPI/TODOGOLANG/internal/core/errors"
)

type Task struct {
	ID            int
	Version       int
	Title         string
	Description   *string
	Completed     bool
	CreatedAt     time.Time
	CompletedAt   *time.Time
	AuthorUserID  int
}

func NewTask(
	id int,
	version int,
	title string,
	description *string,
	completed bool,
	createdAt time.Time,
	completedAt *time.Time,
	authorUserID int,
) Task {
	return Task{
		ID:            id,
		Version:       version,
		Title:         title,
		Description:   description,
		Completed:     completed,
		CreatedAt:     createdAt,
		CompletedAt:   completedAt,
		AuthorUserID:  authorUserID,
	}
}

func NewTaskUninitialized(
	title string,
	description *string,
	completed bool,
	authorUserID int,
) Task {
	now := time.Now()
	return NewTask(
		UninitializedID,
		UninitializedVersion,
		title,
		description,
		completed,
		now,
		nil,
		authorUserID,
	)
}

func (t Task) Validate() error {
	titleLength := len([]rune(t.Title))
	if titleLength < 1 || titleLength > 100 {
		return fmt.Errorf("invalid title length: %d: %w", titleLength, core_errors.ErrInvalidArgument)
	}

	if t.Description != nil {
		descLength := len([]rune(*t.Description))
		if descLength < 1 || descLength > 1000 {
			return fmt.Errorf("invalid description length: %d: %w", descLength, core_errors.ErrInvalidArgument)
		}
	}

	if t.Completed {
		if t.CompletedAt == nil {
			return fmt.Errorf("completed task must have completed_at set: %w", core_errors.ErrInvalidArgument)
		}
		if t.CompletedAt.Before(t.CreatedAt) {
			return fmt.Errorf("completed_at cannot be before created_at: %w", core_errors.ErrInvalidArgument)
		}
	} else {
		if t.CompletedAt != nil {
			return fmt.Errorf("incomplete task must not have completed_at set: %w", core_errors.ErrInvalidArgument)
		}
	}

	if t.AuthorUserID == UninitializedID {
		return fmt.Errorf("author_user_id is required: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}

func (t *Task) MarkCompleted() error {
	if t.Completed {
		return fmt.Errorf("task already completed: %w", core_errors.ErrConflict)
	}
	now := time.Now()
	t.Completed = true
	t.CompletedAt = &now
	t.Version++
	return nil
}

func (t *Task) MarkIncomplete() error {
	if !t.Completed {
		return fmt.Errorf("task already incomplete: %w", core_errors.ErrConflict)
	}
	t.Completed = false
	t.CompletedAt = nil
	t.Version++
	return nil
}

func (t *Task) UpdateTitle(title string) error {
	if len([]rune(title)) < 1 || len([]rune(title)) > 100 {
		return fmt.Errorf("invalid title length: %w", core_errors.ErrInvalidArgument)
	}
	t.Title = title
	t.Version++
	return nil
}

func (t *Task) UpdateDescription(description *string) error {
	if description != nil && (len([]rune(*description)) < 1 || len([]rune(*description)) > 1000) {
		return fmt.Errorf("invalid description length: %w", core_errors.ErrInvalidArgument)
	}
	t.Description = description
	t.Version++
	return nil
}