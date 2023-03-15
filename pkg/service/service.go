package service

import (
	"context"
	"time"

	api "github.com/Wattpad/TaskManager/api/task_manager"
	"github.com/Wattpad/TaskManager/pkg/model"
	"github.com/Wattpad/TaskManager/pkg/storage"
	"github.com/Wattpad/wsl/log"
)

type Service struct {
	logger log.Logger
	db     *storage.Storage
	api.UnimplementedTaskServiceServer
}

func NewService(logger log.Logger, db *storage.Storage) *Service {
	return &Service{logger: logger, db: db}
}

func (s *Service) CreateTask(ctx context.Context, request *api.CreateTaskRequest) (*api.CreateTaskResponse, error) {
	var lastUpdated time.Time
	switch request.GetStatus() {
	case api.TaskStatus_TASK_STATUS_NOT_STARTED:
		lastUpdated = request.GetDetails().GetNotStartedDetails().GetCreatedDate().AsTime()
	case api.TaskStatus_TASK_STATUS_IN_PROGRESS:
		lastUpdated = request.GetDetails().GetInProgressDetails().GetStartedDate().AsTime()
		// TODO add DONE handling
	}

	switch request.GetDetails().GetDetails() {
	case api.TaskDetails_InProgressDetails:
	}

	// TODO include last updated in DB call
	task, err := s.db.CreateTask(ctx, request.GetTitle(), request.GetDescription(), request.GetStatus().String())
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

	return &api.GetTaskResponse{
		Task: task,
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
