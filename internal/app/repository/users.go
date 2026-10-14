package repository

import (
	"web_backend/internal/app/ds"

	"golang.org/x/crypto/bcrypt"
)

func (r *Repository) CreateUser(login string, fullName string, password string) (ds.Users, error) {
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return ds.Users{}, err
	}

	newUser := ds.Users{
		Login:    login,
		FullName: fullName,
		Password: string(passwordHash),
	}

	if err := r.db.Create(&newUser).Error; err != nil {
		return ds.Users{}, err
	}

	return newUser, nil
}

func (r *Repository) GetUserByLogin(login string) (ds.Users, error) {
	var foundUser ds.Users

	err := r.db.Where("login = ?", login).First(&foundUser).Error
	if err != nil {
		return ds.Users{}, err
	}

	return foundUser, nil
}

func (r *Repository) UpdateUserPassword(userID uint, passwordHash string) error {
	return r.db.Model(&ds.Users{}).
		Where("user_id = ?", userID).
		Update("password", passwordHash).Error
}
