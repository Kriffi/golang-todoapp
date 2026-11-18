package tasks_postgres_repository

import (
	"context"
	"fmt"

	"github.com/Kriffi/golang-todoapp/internal/core/domain"
)

func (r *TasksRepository) GetTasks(
	ctx context.Context,
	userID *int,
	limit *int,
	offset *int,
	overdue bool,
) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
 SELECT id, version, title, description, completed, created_at, completed_at, author_user_id, to_char(due_date, 'YYYY-MM-DD')
 FROM todoapp.tasks
 WHERE ($3::int IS NULL OR author_user_id=$3)
   AND (NOT $4::boolean OR (NOT completed AND due_date < (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')::date))
 ORDER BY id ASC
 LIMIT $1
 OFFSET $2
 `

	rows, err := r.pool.Query(
		ctx,
		query,
		limit,
		offset,
		userID,
		overdue,
	)
	if err != nil {
		return nil, fmt.Errorf("select tasks: %w", err)
	}
	defer rows.Close()

	var taskModels []TaskModel

	for rows.Next() {
		var taskModel TaskModel

		err := rows.Scan(
			&taskModel.ID,
			&taskModel.Version,
			&taskModel.Title,
			&taskModel.Description,
			&taskModel.Completed,
			&taskModel.CreatedAt,
			&taskModel.CompletedAt,
			&taskModel.AuthorUserID,
			&taskModel.DueDate,
		)
		if err != nil {
			return nil, fmt.Errorf("scan tasks: %w", err)
		}

		taskModels = append(taskModels, taskModel)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	taskDomains := taskDomainsFromModels(taskModels)

	return taskDomains, nil
}
