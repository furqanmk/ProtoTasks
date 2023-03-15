package storage_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"

	"github.com/Wattpad/TaskManager/pkg/storage"
	"github.com/Wattpad/wsl/log"
)

type StorageTestSuite struct {
	suite.Suite
	mock   sqlmock.Sqlmock
	db     *sqlx.DB
	repo   storage.Storage
	logger log.Logger
}

func TestStorageTestSuite(t *testing.T) {
	suite.Run(t, new(StorageTestSuite))
}

func (suite *StorageTestSuite) SetupTest() {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	suite.Require().NoError(err)

	suite.logger = log.NewNopLogger()
	suite.db = sqlx.NewDb(db, "sqlmock")
	suite.mock = mock
	suite.repo = storage.Storage{DB: suite.db, Logger: suite.logger}
}

func (suite *StorageTestSuite) TestCreateTask() {
	// TODO Set expectations
	expectedId := int64(1000001)
	title := "test_title"
	description := "test_desc"
	status := "test_status"

	suite.mock.ExpectExec("INSERT INTO task(title, description, status) VALUES(?, ?, ?)").
		WithArgs(title, description, status).
		WillReturnResult(sqlmock.NewResult(expectedId, 1))

	// TODO Call the method under test
	task, err := suite.repo.CreateTask(context.Background(), title, description, status)

	// TODO Verify method acted as expected
	suite.NoError(err)
	suite.NotNil(task)

	suite.Equal(expectedId, task.Id)
	suite.Equal(title, task.Title)
	suite.Equal(description, task.Description)
	suite.Equal(status, task.Status)
}

func (suite *StorageTestSuite) TestGetTaskById() {
	id := int64(1000001)
	title := "test_title"
	description := "test_desc"
	status := "test_status"

	rows := sqlmock.NewRows([]string{"id", "title", "description", "status", "assignee"}).
		AddRow(id, title, description, status, nil)

	suite.mock.ExpectQuery("SELECT * FROM task WHERE id = ?").
		WithArgs(id).
		WillReturnRows(rows)

	task, err := suite.repo.GetTaskByID(context.Background(), id)

	suite.NoError(err)
	suite.NotNil(task)

	suite.Equal(id, task.Id)
	suite.Equal(title, task.Title)
	suite.Equal(description, task.Description)
	suite.Equal(status, task.Status)
	suite.Nil(task.Assignee)
}

func (suite *StorageTestSuite) TestUpdateTask() {
	id := int64(1000001)
	title := "updated_title"
	description := "updated_desc"
	status := "updated_status"

	suite.mock.ExpectExec("UPDATE task SET title = ?, description = ?, status = ? WHERE id = ?").
		WithArgs(title, description, status, id).
		WillReturnResult(sqlmock.NewResult(id, 1))

	task, err := suite.repo.UpdateTask(context.Background(), id, title, description, status)

	suite.NoError(err)
	suite.NotNil(task)

	suite.Equal(id, task.Id)
	suite.Equal(title, task.Title)
	suite.Equal(description, task.Description)
	suite.Equal(status, task.Status)
	suite.Nil(task.Assignee)
}

func (suite *StorageTestSuite) TestDeleteTask() {
	id := int64(1000001)

	suite.mock.ExpectExec("DELETE FROM task WHERE id = ?").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(id, 1))

	err := suite.repo.DeleteTask(context.Background(), id)

	suite.NoError(err)
}
