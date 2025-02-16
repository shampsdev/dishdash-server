package repo

import (
	"context"

	"dashboard.dishdash.ru/pkg/domain"
)

type Task interface {
	CreateTask(ctx context.Context, task *domain.Task) (int64, error)
	GetTaskByID(ctx context.Context, taskId int64) (*domain.Task, error)
	GetAllTasks(ctx context.Context) ([]*domain.Task, error)
	UpdateTask(ctx context.Context, task *domain.Task) (*domain.Task, error)
	DeleteTask(ctx context.Context, taskId int64) error
}
