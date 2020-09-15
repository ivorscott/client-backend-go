package column

import (
	"time"
)

type Column struct {
	ID        string    `db:"column_id" json:"id"`
	ProjectID string    `db:"project_id" json:"projectId"`
	Title     string    `db:"title" json:"title"`
	Column    string    `db:"column" json:"column"`
	TaskIDS   []string  `db:"task_ids" json:"taskIds"`
	Created   time.Time `db:"created" json:"created"`
}

type NewColumn struct {
	ProjectID string   `json:"projectId"`
	Title     string   `json:"title"`
	Column    string   `json:"column"`
	TaskIDS   []string `json:"taskIds"`
}

type UpdateColumn struct {
	Title   string   `json:"title"`
	Column  string   `json:"column"`
	TaskIDS []string `json:"taskIds"`
}
