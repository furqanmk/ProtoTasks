package storage

import (
	"context"
	"time"

	_ "github.com/go-sql-driver/mysql" // required for connecting to mysql
	"github.com/jmoiron/sqlx"

	api "github.com/Wattpad/TaskManager/api/task_manager"
	"github.com/Wattpad/TaskManager/pkg/model"
	"github.com/Wattpad/wsl/log"
	"github.com/Wattpad/wsl/sql"
)

type Storage struct {
	DB     *sqlx.DB
	Logger log.Logger
}

func NewStorage(db *sql.DB, logger log.Logger, driverName sql.DriverName) (*Storage, error) {
	dbConn := sqlx.NewDb(db, driverName.String())

	return &Storage{
		Logger: logger,
		DB:     dbConn,
	}, nil
}

func (s *Storage) CreateTask(ctx context.Context, title, description, status string, lastUpdated time.Time) (*api.Task, error) {
	query := "INSERT INTO task(title, description, status, last_updated) VALUES(?, ?, ?, ?)"
	result, err := s.DB.ExecContext(ctx, query, title, description, status, lastUpdated)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	task := api.Task{
		Id:          id,
		Title:       title,
		Description: description,
		Status:      api.TaskStatus(api.TaskStatus_value[status]),
	}

	return &task, nil
}

func (s *Storage) GetTaskByID(ctx context.Context, id int64) (*model.Task, error) {
	var task model.Task
	// functionally equivalent to task := api.Task{}

	query := "SELECT * FROM task WHERE id = ?"

	// SelectContext returns multiple rows
	// GetContext returns one row
	if err := s.DB.GetContext(ctx, &task, query, id); err != nil {
		return nil, err
	}

	return &task, nil
}

func (s *Storage) UpdateTask(ctx context.Context, id int64, title, description, status string) (*api.Task, error) {
	query := "UPDATE task SET title = ?, description = ?, status = ? WHERE id = ?"
	_, err := s.DB.ExecContext(ctx, query, title, description, status, id)
	if err != nil {
		return nil, err
	}

	task := api.Task{
		Id:          id,
		Title:       title,
		Description: description,
		Status:      api.TaskStatus(api.TaskStatus_value[status]),
	}

	return &task, nil
}

func (s *Storage) DeleteTask(ctx context.Context, id int64) error {
	query := "DELETE FROM task WHERE id = ?"
	_, err := s.DB.ExecContext(ctx, query, id)
	return err
}

func (s *Storage) GetAllTasks() (*sqlx.Rows, error) {
	return s.DB.Queryx("SELECT * FROM task")
}
