package dto

import "errors"

type DTO interface {
	Validate() error
}

type UserRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	SecretKey string `json:"secretKey"`
}

func (u *UserRequest) Validate() error {
	if u.Username == "" || u.Password == "" {
		return errors.New("username or password cannot be empty")
	}

	return nil
}
