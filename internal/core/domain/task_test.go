package domain

import "testing"

func TestTaskDeadline(t *testing.T) {
	due := "2026-11-18"
	task := NewTaskUninitialized("Deadline test", nil, 1)
	if err := task.Validate(); err != nil {
		t.Fatalf("task without deadline: %v", err)
	}
	if err := task.ApplyPatch(TaskPatch{DueDate: Nullable[string]{Set: true, Value: &due}}); err != nil {
		t.Fatal(err)
	}

	if err := task.ApplyPatch(TaskPatch{}); err != nil || task.DueDate == nil || *task.DueDate != due {
		t.Fatalf("omitted deadline must be preserved: %v", err)
	}

	invalid := "2026-02-30"
	if err := task.ApplyPatch(TaskPatch{DueDate: Nullable[string]{Set: true, Value: &invalid}}); err == nil || task.DueDate == nil || *task.DueDate != due {
		t.Fatalf("invalid date must be rejected without changing the task: %v", err)
	}

	if err := task.ApplyPatch(TaskPatch{DueDate: Nullable[string]{Set: true}}); err != nil || task.DueDate != nil {
		t.Fatalf("null must clear the deadline: %v", err)
	}
}
