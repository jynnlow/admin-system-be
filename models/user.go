package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
}

type UserCRUDOperations interface {
	InsertUsers(string, string) error
	DeleteUsers(int) error
	UpdateUsers(int, string, string) error
	GetUserByID(int) (*User, error)
	GetAllUsers() ([]*User, error)
	GetPwdByUsername(string) (*string, error)
}

type UserCRUDOperationsImpl struct {
	DbConn *gorm.DB
}

func (u *UserCRUDOperationsImpl) InsertUsers(username, password string) error {
	newUser := &User{
		Username: username,
		Password: password,
	}

	err := u.DbConn.Select("username", "password").Create(newUser).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *UserCRUDOperationsImpl) DeleteUsers(id uint) error {
	//find user with given id
	foundUser, err := u.GetUserByID(id)
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

func (u *UserCRUDOperationsImpl) UpdateUsers(id uint, username string, password string) error {
	//find user with given id
	foundUser, err := u.GetUserByID(id)
	if err != nil {
		return err
	}
	//check if username and password is empty
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

func (u *UserCRUDOperationsImpl) GetAllUsers() ([]*User, error) {
	var users []*User
	err := u.DbConn.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *UserCRUDOperationsImpl) GetUserByID(id uint) (*User, error) {
	user := &User{}
	err := u.DbConn.First(user, id).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserCRUDOperationsImpl) GetPwdByUsername(username string) (*string, error) {
	user := &User{}
	err := u.DbConn.Where("username = ?", username).First(user).Error
	if err != nil {
		return nil, err
	}
	return &user.Password, nil
}
