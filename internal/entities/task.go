package entities

import (
	"fmt"
	"time"
)

type TaskWithTag struct {
	ID      int64     `db:"id"`
	Task    string    `db:"task"`
	Date    time.Time `db:"date"`
	Tag     *string   `db:"tag,omitempty"`
	OwnerID int64     `db:"owner_id"`
}

func (task *TaskWithTag) String() string {
	return fmt.Sprintf("task ID: %d\n task text: %s\n task tag: %s\n task date %v\n",
		task.ID, task.Task, task.Tag, task.Date)
}
