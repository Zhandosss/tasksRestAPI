package verification

import (
	"restAPI/internal/model"
	"time"
)

func Task(task model.Task, callTime time.Time) bool {
	if task.Text == "" {
		return false
	}
	if callTime.Sub(task.Date) >= time.Second {
		return false
	}
	return true
}
