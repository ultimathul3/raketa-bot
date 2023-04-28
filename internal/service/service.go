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

func (r *RaketaService) SignUp(ctx context.Context, id int64, username string) error {
	_, err := r.client.SignUp(ctx, &raketapb.SignUpRequest{Id: id, Username: username})
	return err
}

func (r *RaketaService) GetUserRole(ctx context.Context, username string) (types.Role, error) {
	response, err := r.client.GetUserRole(ctx, &raketapb.GetUserRoleRequest{Username: username})
	if err != nil {
		return types.UnknownRole, err
	}
	return convertProtoRoleToTypes(response.Role), err
}

func (r *RaketaService) CreateTask(ctx context.Context, url string) error {
	_, err := r.client.CreateTask(ctx, &raketapb.CreateTaskRequest{Url: url})
	return err
}

func (r *RaketaService) SetTaskPrice(ctx context.Context, url string, price uint64) error {
	_, err := r.client.SetTaskPrice(ctx, &raketapb.SetTaskPriceRequest{Url: url, Price: price})
	return err
}

func (r *RaketaService) DeleteTask(ctx context.Context, url string) error {
	_, err := r.client.DeleteTask(ctx, &raketapb.DeleteTaskRequest{Url: url})
	return err
}

func (r *RaketaService) AssignUser(ctx context.Context, url, username string) error {
	_, err := r.client.AssignUser(ctx, &raketapb.AssignUserRequest{
		Url:      url,
		Username: username,
	})
	return err
}

func (r *RaketaService) CloseTask(ctx context.Context, url string) error {
	_, err := r.client.CloseTask(ctx, &raketapb.CloseTaskRequest{
		Url: url,
	})
	return err
}

func (r *RaketaService) GetUnassignTasks(ctx context.Context) ([]types.Task, error) {
	var tasks []types.Task
	response, err := r.client.GetUnassignTasks(ctx, &raketapb.GetUnassignTasksRequest{})

	for _, task := range response.Tasks {
		tasks = append(tasks, types.Task{
			Url:    task.Url,
			Status: convertProtoStatusToTypes(task.Status),
			UserID: task.UserId,
			Price:  task.Price,
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

func convertProtoRoleToTypes(role raketapb.GetUserRoleResponse_Role) types.Role {
	switch role {
	case raketapb.GetUserRoleResponse_ADMIN:
		return types.AdminRole
	case raketapb.GetUserRoleResponse_REGULAR:
		return types.RegularRole
	default:
		return types.UnknownRole
	}
}
