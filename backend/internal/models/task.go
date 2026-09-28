package models

import "time"

type TaskType string

const (
	TaskTypeRoutine TaskType = "routine"
	TaskTypeAdhoc   TaskType = "adhoc"
)

type Task struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Type        TaskType   `json:"type"`
	Frequency   string     `json:"frequency,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
