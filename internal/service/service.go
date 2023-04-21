package service

import (
	"context"
	"fmt"

	raketapb "github.com/vanyaio/raketa-backend/proto"
)

type Service interface {
	SignUp(ctx context.Context, id int64) error
	CreateTask(ctx context.Context, url string) error
	DeleteTask(ctx context.Context, url string) error
	AssignUser(ctx context.Context, url string, userID int64) error
	CloseTask(ctx context.Context, url string) error
	GetOpenTasks(ctx context.Context) ([]*raketapb.Task, error)
}

type RaketaService struct {
	client raketapb.RaketaServiceClient
}

func NewRaketaService(client raketapb.RaketaServiceClient) *RaketaService {
	return &RaketaService{
		client: client,
	}
}

func (r *RaketaService) SignUp(ctx context.Context, id int64) error {
	a, err := r.client.SignUp(ctx, &raketapb.SignUpRequest{Id: id})
	fmt.Println(a, err.Error())
	return err
}

func (r *RaketaService) CreateTask(ctx context.Context, url string) error {
	_, err := r.client.CreateTask(ctx, &raketapb.CreateTaskRequest{Url: url})
	return err
}

func (r *RaketaService) DeleteTask(ctx context.Context, url string) error {
	_, err := r.client.DeleteTask(ctx, &raketapb.DeleteTaskRequest{Url: url})
	return err
}

func (r *RaketaService) AssignUser(ctx context.Context, url string, userID int64) error {
	_, err := r.client.AssignUser(ctx, &raketapb.AssignUserRequest{
		Url:    url,
		UserId: userID,
	})
	return err
}

func (r *RaketaService) CloseTask(ctx context.Context, url string) error {
	_, err := r.client.CloseTask(ctx, &raketapb.CloseTaskRequest{
		Url: url,
	})
	return err
}

func (r *RaketaService) GetOpenTasks(ctx context.Context) ([]*raketapb.Task, error) {
	response, err := r.client.GetOpenTasks(ctx, &raketapb.GetOpenTasksRequest{})
	return response.Tasks, err
}
