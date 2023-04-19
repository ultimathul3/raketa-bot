package service

import (
	"context"

	"github.com/vanyaio/raketa-bot/internal/domain"
)

type Service interface {
	SignUp(ctx context.Context, id int64) error
	CreateTask(ctx context.Context, url string) error
	DeleteTask(ctx context.Context, url string) error
	AssignWorker(ctx context.Context, url string, userID int64) error
	CloseTask(ctx context.Context, url string) error
	GetOpenTasks(ctx context.Context) ([]domain.Task, error)
}
