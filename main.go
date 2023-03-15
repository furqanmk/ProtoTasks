package main

import (
	"google.golang.org/grpc/reflection"

	api "github.com/Wattpad/TaskManager/api/task_manager"
	"github.com/Wattpad/TaskManager/internal/config"
	"github.com/Wattpad/TaskManager/pkg/service"
	"github.com/Wattpad/TaskManager/pkg/storage"
	"github.com/Wattpad/wsl"
	"github.com/Wattpad/wsl/grpc"
	"github.com/Wattpad/wsl/sql"
)

// EntryPoint is typically where you will, for instance, attach HTTP
// handlers to a router or register a GRPC API implementation with a server.
func EntryPoint(grpcServer *grpc.Server, taskService *service.Service) {
	api.RegisterTaskServiceServer(grpcServer, taskService)
	reflection.Register(grpcServer)
}

func main() {
	wsl.New[config.Config]("TaskManager",
		EntryPoint,
		grpc.ServerModule,
		service.NewService,
		sql.Module("taskmanager"),
		storage.NewStorage,
	).Run()
}
