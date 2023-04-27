package service

import (
	"context"

	raketapb "github.com/vanyaio/raketa-backend/proto"
	"github.com/vanyaio/raketa-bot/internal/types"
)

type RaketaService struct {
	client raketapb.RaketaServiceClient
}

func NewRaketaService(client raketapb.RaketaServiceClient) *RaketaService {
	return &RaketaService{
		client: client,
	}
}

func (r *RaketaService) SignUp(ctx context.Context, id int64) error {
	_, err := r.client.SignUp(ctx, &raketapb.SignUpRequest{Id: id})
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

func (r *RaketaService) GetOpenTasks(ctx context.Context) ([]types.Task, error) {
	var tasks []types.Task
	response, err := r.client.GetOpenTasks(ctx, &raketapb.GetOpenTasksRequest{})

	for _, task := range response.Tasks {
		tasks = append(tasks, types.Task{
			Url:    task.Url,
			Status: convertProtoStatusToTypes(task.Status),
			UserID: task.UserId,
		})
	}

	return tasks, err
}

func convertProtoStatusToTypes(status raketapb.Task_Status) types.Status {
	switch status {
	case raketapb.Task_OPEN:
		return types.TaskOpen
	case raketapb.Task_CLOSED:
		return types.TaskClosed
	case raketapb.Task_DECLINED:
		return types.TaskDeclined
	default:
		return types.TaskUnknown
	}
}
