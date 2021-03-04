package userssvc

import (
	"gin-example/models"
	"gin-example/pkg/app"
	"gin-example/pkg/errcode"
)

type User struct {
	ID       uint
	Name     string
	Password string
	Role     string
	Email    string
	Gender   string

	PageNumber int
	PageSize   int
}

type UsersList struct {
	Users      []models.User
	TotalCount uint
}

// for swagger show Response
type UsersListResponse struct {
	*errcode.ErrorMessage
	Data *UsersList
}

func (u *User) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})

	if u.Name != "" {
		maps["name"] = u.Name
	}
	if u.ID > 0 {
		maps["id"] = u.ID
	}
	if u.Role != "" {
		maps["role"] = u.Role
	}
	if u.Email != "" {
		maps["email"] = u.Email
	}
	if u.Gender != "" {
		maps["gender"] = u.Gender
	}

	return maps
}

func (u *User) GetUsers() (*UsersListResponse, error) {
	users, err := models.GetUsers(u.PageNumber, u.PageSize, u.getMaps())
	if err != nil {
		return nil, err
	}
	usersList := &UsersList{Users: users}

	count, err := models.GetUsersTotal(u.getMaps())
	if err != nil {
		return nil, err
	}
	usersList.TotalCount = count

	return &UsersListResponse{
		ErrorMessage: errcode.Success,
		Data:         usersList,
	}, nil
}

func (u *User) ExistByName() (bool, error) {
	return models.UserExistByName(u.Name)
}

func (u *User) Add() error {
	return models.AddUser(u.Name, u.Password, u.Role, u.Email, u.Gender)
}

func (u *User) CheckPassword() error {
	id, password, err := models.GetUserPasswordByName(u.Name)
	if err != nil {
		return err
	}

	if err := app.Compare(password, u.Password); err != nil {
		return err
	}

	u.ID = id

	return nil
}
