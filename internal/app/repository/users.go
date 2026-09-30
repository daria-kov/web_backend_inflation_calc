package repository

import (
	"web_backend/internal/app/ds"
)

func (r *Repository) CreateUser(fullName string, email string) (ds.Users, error) {
	newUser := ds.Users{
		FullName: fullName,
		Email:    email,
	}

	if err := r.db.Create(&newUser).Error; err != nil {
		return ds.Users{}, err
	}

	return newUser, nil
}
