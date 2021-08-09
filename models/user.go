package models

import (
	"gorm.io/gorm"
)

type UserCRUDOperations interface {
	Insert(*User) error
	Delete(uint) error
	Update(uint, string, string) error
	GetByID(uint) (*User, error)
	GetByUsername(string) (*User, error)
	GetAll() ([]*User, error)
}

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Approved bool   `json:"approved"`
}

type UserCRUDOperationsImpl struct {
	DbConn *gorm.DB
}

func (u *UserCRUDOperationsImpl) Insert(user *User) error {
	err := u.DbConn.Select("username", "password", "role", "approved").Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *UserCRUDOperationsImpl) Delete(id uint) error {
	foundUser, err := u.GetByID(id)
	if err != nil {
		return err
	}
	//delete only if the user exists
	//permanently deleted with Unscoped().Delete()
	err = u.DbConn.Unscoped().Delete(foundUser, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *UserCRUDOperationsImpl) Update(id uint, username string, password string) error {
	foundUser, err := u.GetByID(id)
	if err != nil {
		return err
	}
	if username != "" && foundUser.Username != username {
		foundUser.Username = username
	}
	if password != "" && foundUser.Password != password {
		foundUser.Password = password
	}
	//update user with all field
	err = u.DbConn.Save(foundUser).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *UserCRUDOperationsImpl) GetAll() ([]*User, error) {
	var users []*User
	err := u.DbConn.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *UserCRUDOperationsImpl) GetByID(id uint) (*User, error) {
	user := &User{}
	err := u.DbConn.First(user, id).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserCRUDOperationsImpl) GetByUsername(username string) (*User, error) {
	user := &User{}
	err := u.DbConn.Where("username = ?", username).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
