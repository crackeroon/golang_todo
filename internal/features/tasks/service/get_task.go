package tasks_service

import (
	"context"
	"fmt"

	"github.com/crackeroon/golang_todo/internal/core/domain"
)

func (s *TasksService) GetTask(
	ctx context.Context,
	id int,
) (domain.Task, error) {
	task, err := s.tasksRepository.GetTask(ctx, id)
	if err != nil {
		return domain.Task{}, fmt.Errorf("Get Task from repository: %w", err)
	}
	return task, nil
}
