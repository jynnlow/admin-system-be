package entity

import (
	"admin-system-be/constants"
	"errors"
)

type User struct {
	Username string
	Password string
	Role     string
	Approved bool
}

// NewUserEntity is constructor of user ...
func NewUserEntity(username, password, secretKey, dbSecretKey string) (*User, error) {
	if secretKey == "" {
		return &User{
			Username: username,
			Password: password,
			Role:     constants.USER,
			Approved: false,
		}, nil

	} else if secretKey != "" && secretKey == dbSecretKey {
		return &User{
			Username: username,
			Password: password,
			Role:     constants.ADMIN,
			Approved: true,
		}, nil
	}

	return nil, errors.New("wrong admin secret key")
}
