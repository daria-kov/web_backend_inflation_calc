package user

import "web_backend/internal/app/ds"

const currentUserID uint = 1

var currentUser *ds.Users

func GetCurrentUser() *ds.Users {
	if currentUser == nil {
		currentUser = &ds.Users{
			UserID: currentUserID,
		}
	}

	return currentUser
}
