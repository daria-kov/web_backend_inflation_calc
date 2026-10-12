package repository

import (
	"web_backend/internal/app/ds"
)

func (r *Repository) CreateUser(fullName string, password string) (ds.Users, error) {
	newUser := ds.Users{
		FullName: fullName,
		Password: password,
	}

	if err := r.db.Create(&newUser).Error; err != nil {
		return ds.Users{}, err
	}

	return newUser, nil
}
