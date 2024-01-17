package verification

import "restAPI/internal/model"

func User(user model.User) bool {
	if user.FirstName == "" || user.LastName == "" || user.Login == "" || user.Password == "" {
		return false
	}
	return true
}
