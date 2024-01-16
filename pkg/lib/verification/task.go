package verification

import (
	"restAPI/internal/model"
)

func Task(task model.Task) bool {
	if task.Text == "" {
		return false
	}
	return true
}
