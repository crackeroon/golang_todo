package statistics_postgres_repository

import (
	"time"

	"github.com/crackeroon/golang_todo/internal/core/domain"
)

type TaskModel struct {
	ID           int
	Version      int
	Title        string
	Description  *string
	Completed    bool
	CreatedAt    time.Time
	CompletedAt  *time.Time
	AuthorUserId int
}

func taskDomainFromModel(taskModel TaskModel) domain.Task {
	return domain.NewTask(
		taskModel.ID,
		taskModel.Version,
		taskModel.Title,
		taskModel.Description,
		taskModel.Completed,
		taskModel.CreatedAt,
		taskModel.CompletedAt,
		taskModel.AuthorUserId,
	)
}

func taskDomainsFromModels(tasks []TaskModel) []domain.Task {
	domains := make([]domain.Task, len(tasks))
	for i, t := range tasks {
		domains[i] = taskDomainFromModel(t)
	}
	return domains
}
