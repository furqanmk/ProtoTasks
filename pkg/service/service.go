package service

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	api "github.com/Wattpad/TaskManager/api/task_manager"
	"github.com/Wattpad/TaskManager/pkg/model"
	"github.com/Wattpad/TaskManager/pkg/storage"
	"github.com/Wattpad/wsl/log"
)

type Service struct {
	logger log.Logger
	db     *storage.Storage
	cache  *storage.Cache
	api.UnimplementedTaskServiceServer
}

func NewService(logger log.Logger, db *storage.Storage, cache *storage.Cache) *Service {
	return &Service{logger: logger, db: db, cache: cache}
}

func (s *Service) CreateTask(ctx context.Context, request *api.CreateTaskRequest) (*api.CreateTaskResponse, error) {
	var lastUpdated time.Time

	switch details := request.GetDetails().GetDetails().(type) {
	case *api.TaskDetails_NotStartedDetails:
		lastUpdated = details.NotStartedDetails.GetCreatedDate().AsTime()
	case *api.TaskDetails_InProgressDetails:
		lastUpdated = details.InProgressDetails.GetStartedDate().AsTime()
	case *api.TaskDetails_DoneDetails:
		lastUpdated = details.DoneDetails.GetCompletedDate().AsTime()
	}

	task, err := s.db.CreateTask(ctx, request.GetTitle(), request.GetDescription(), request.GetStatus().String(), lastUpdated)
	if err != nil {
		return nil, err
	}

	return &api.CreateTaskResponse{
		Task: task,
	}, nil
}

func (s *Service) GetTask(ctx context.Context, request *api.GetTaskRequest) (*api.GetTaskResponse, error) {
	task, err := s.db.GetTaskByID(ctx, request.GetId())
	if err != nil {
		return nil, err
	}

	status := api.TaskStatus(api.TaskStatus_value[task.Status])
	var details api.TaskDetails

	switch status {
	case api.TaskStatus_TASK_STATUS_NOT_STARTED:
		details = api.TaskDetails{
			Details: &api.TaskDetails_NotStartedDetails{
				NotStartedDetails: &api.NotStartedDetails{CreatedDate: timestamppb.New(*task.LastUpdated)},
			},
		}
	case api.TaskStatus_TASK_STATUS_IN_PROGRESS:
		details = api.TaskDetails{
			Details: &api.TaskDetails_InProgressDetails{
				InProgressDetails: &api.InProgressDetails{StartedDate: timestamppb.New(*task.LastUpdated)},
			},
		}
	case api.TaskStatus_TASK_STATUS_DONE:
		details = api.TaskDetails{
			Details: &api.TaskDetails_DoneDetails{
				DoneDetails: &api.DoneTaskDetails{CompletedDate: timestamppb.New(*task.LastUpdated)},
			},
		}
	}

	apiTask := api.Task{
		Id:          task.Id,
		Title:       task.Title,
		Description: task.Description,
		Status:      status,
		Details:     &details,
	}

	return &api.GetTaskResponse{
		Task: &apiTask,
	}, nil
}

func (s *Service) UpdateTask(ctx context.Context, request *api.UpdateTaskRequest) (*api.UpdateTaskResponse, error) {
	task, err := s.db.UpdateTask(
		ctx,
		request.GetTask().GetId(),
		request.GetTask().GetTitle(),
		request.GetTask().GetDescription(),
		request.GetTask().GetStatus().String(),
	)
	if err != nil {
		return nil, err
	}

	return &api.UpdateTaskResponse{
		Task: task,
	}, nil
}

func (s *Service) DeleteTask(ctx context.Context, request *api.DeleteTaskRequest) (*api.DeleteTaskResponse, error) {
	err := s.db.DeleteTask(ctx, request.GetId())

	return &api.DeleteTaskResponse{
		Success: err == nil,
	}, nil
}

func (s *Service) GetAllTasks(request *api.GetAllTasksRequest, tasksStream api.TaskService_GetAllTasksServer) error {
	rows, err := s.db.GetAllTasks()
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var task model.Task
		err := rows.StructScan(&task)
		if err != nil {
			return err
		}

		err = tasksStream.Send(&api.GetAllTasksResponse{
			Task: &api.Task{
				Id:          task.Id,
				Title:       task.Title,
				Description: task.Description,
				Status:      api.TaskStatus(api.TaskStatus_value[task.Status]),
				Assignee:    task.Assignee,
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) GetCache(ctx context.Context, request *api.GetCacheRequest) (*api.GetCacheResponse, error) {
	template, err := s.cache.FetchTemplate(request.GetLastName())
	if err != nil {
		return nil, err
	}

	return &api.GetCacheResponse{
		LastName: template.LastName,
		Birthday: template.Birthday,
	}, nil
}

func (s *Service) SetCache(ctx context.Context, request *api.SetCacheRequest) (*api.SetCacheResponse, error) {
	err := s.cache.CacheTemplate(&storage.Template{
		LastName: request.LastName,
		Birthday: request.Birthday,
	})

	return &api.SetCacheResponse{
		Success: err != nil,
	}, nil
}
