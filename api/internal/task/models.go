package task

import (
	"time"
)

type Task struct {
	ID      string    `db:"task_id" json:"id"`
	Title   string    `db:"title" json:"title"`
	Content *string   `db:"content" json:"content"`
	Created time.Time `db:"created" json:"created"`
}

type NewTask struct {
	Title   string  `db:"title" json:"title"`
	Content *string `db:"content" json:"content"`
}

type UpdateTask struct {
	Title   string  `json:"title" validate:"required"`
	Content *string `json:"content"`
}
