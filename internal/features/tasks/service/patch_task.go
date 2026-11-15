package tasks_service

import (
	"context"
	"fmt"

	"github.com/Kriffi/golang-todoapp/internal/core/domain"
	core_errors "github.com/Kriffi/golang-todoapp/internal/core/errors"
)

func (s *TasksService) PatchTask(
	ctx context.Context,
	id int,
	patch domain.TaskPatch,
	expectedVersion *int,
) (domain.Task, error) {
	task, err := s.tasksRepository.GetTask(ctx, id)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get task: %w", err)
	}

	if expectedVersion != nil && task.Version != *expectedVersion {
		return domain.Task{}, fmt.Errorf("task with id='%d' has changed: %w", id, core_errors.ErrConflict)
	}

	if err := task.ApplyPatch(patch); err != nil {
		return domain.Task{}, fmt.Errorf("apply task patch: %w", err)
	}

	patchedTask, err := s.tasksRepository.PatchTask(ctx, id, task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("patch task: %w", err)
	}

	return patchedTask, nil
}
