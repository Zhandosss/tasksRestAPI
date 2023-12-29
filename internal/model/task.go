package model

import (
	"fmt"
	"time"
)

type Task struct {
	ID   int64     `json:"id" db:"task_id"`
	Text string    `json:"text" db:"task"`
	Tags []string  `json:"tags" db:"omitempty"`
	Date time.Time `json:"due" db:"date"`
}

func (task *Task) String() string {
	return fmt.Sprintf("task ID: %d\n task text: %s\n task tags: %v\n task date %v\n",
		task.ID, task.Text, task.Tags, task.Date)
}
