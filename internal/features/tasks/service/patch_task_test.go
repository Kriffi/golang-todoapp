package tasks_service

import (
	"context"
	"errors"
	"testing"

	"github.com/Kriffi/golang-todoapp/internal/core/domain"
	core_errors "github.com/Kriffi/golang-todoapp/internal/core/errors"
)

type patchTaskRepository struct {
	TasksRepository
	task       domain.Task
	writes     int
	writeError error
}

func (r *patchTaskRepository) GetTask(_ context.Context, _ int) (domain.Task, error) {
	return r.task, nil
}

func (r *patchTaskRepository) PatchTask(_ context.Context, _ int, task domain.Task) (domain.Task, error) {
	r.writes++
	if r.writeError != nil {
		return domain.Task{}, r.writeError
	}
	task.Version++
	r.task = task
	return task, nil
}

func TestPatchTaskChecksClientVersion(t *testing.T) {
	repository := &patchTaskRepository{task: domain.Task{ID: 1, Version: 3, Title: "Original"}}
	service := NewTasksService(repository)
	originalVersion := 3
	firstTitle := "First edit"
	secondTitle := "Stale edit"

	updated, err := service.PatchTask(context.Background(), 1, domain.TaskPatch{
		Title: domain.Nullable[string]{Value: &firstTitle, Set: true},
	}, &originalVersion)
	if err != nil || updated.Version != 4 || updated.Title != firstTitle {
		t.Fatalf("first edit: got task %+v, error %v", updated, err)
	}

	_, err = service.PatchTask(context.Background(), 1, domain.TaskPatch{
		Title: domain.Nullable[string]{Value: &secondTitle, Set: true},
	}, &originalVersion)
	if !errors.Is(err, core_errors.ErrConflict) {
		t.Fatalf("stale edit: expected a conflict, got %v", err)
	}
	if repository.writes != 1 || repository.task.Title != firstTitle {
		t.Fatalf("stale edit reached storage: %+v", repository)
	}
}

func TestPatchTaskPreservesRepositoryConflict(t *testing.T) {
	repository := &patchTaskRepository{
		task:       domain.Task{ID: 1, Version: 3, Title: "Original"},
		writeError: core_errors.ErrConflict,
	}
	service := NewTasksService(repository)
	newTitle := "Concurrent edit"

	_, err := service.PatchTask(context.Background(), 1, domain.TaskPatch{
		Title: domain.Nullable[string]{Value: &newTitle, Set: true},
	}, nil)
	if !errors.Is(err, core_errors.ErrConflict) || repository.writes != 1 {
		t.Fatalf("repository write: count %d, error %v", repository.writes, err)
	}
}
