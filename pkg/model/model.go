package model

type Task struct {
	Id          int64
	Title       string
	Description string
	Status      string
	Assignee    *string
}
