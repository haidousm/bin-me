package mocks

import (
	"time"

	"binme.haido.us/internal/models"
)

type UserModel struct{}

var mockUser = models.User{
	ID:             1,
	Name:           "Alice Elica",
	Email:          "alice@example.com",
	HashedPassword: []byte("pa$$word"),
	Created:        time.Now(),
}

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "dupe@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Get(id int) (models.User, error) {
	switch id {
	case 1:
		return mockUser, nil
	default:
		return models.User{}, models.ErrNoRecord
	}
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	if email == mockUser.Email && password == string(mockUser.HashedPassword) {
		return 1, nil
	}

	return 0, models.ErrInvalidCredentials
}

func (m *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}

func (m *UserModel) PasswordUpdate(id int, currentPassword, newPassword string) error {
	if id == 1 {
		if currentPassword != string(mockUser.HashedPassword) {
			return models.ErrInvalidCredentials
		}

		return nil
	}

	return models.ErrNoRecord
}
