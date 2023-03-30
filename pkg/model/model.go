package model

import (
	"time"
)

type Task struct {
	Id          int64
	Title       string
	Description string
	Status      string
	Assignee    *string
	LastUpdated *time.Time `db:"last_updated"`
}
