package models

import (
	"gorm.io/gorm"

	"gin-example/pkg/app"
	"gin-example/pkg/database"
)

type User struct {
	gorm.Model

	Name     string `json:"name"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Email    string `json:"email"`
	Gender   string `json:"gender"`
}

func GetUsers(pageNumber, pageSize int, maps interface{}) ([]User, error) {
	var users []User
	pageOffset := app.GetPageOffset(pageNumber, pageSize)

	db := database.GetGormDB().Model(&User{}).Select("id, name, role, created_at, updated_at, deleted_at, email, gender")
	if err := db.Offset(pageOffset).Limit(pageSize).Where(maps).Scan(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func GetUsersTotal(maps map[string]interface{}) (uint, error) {
	var count int64

	if err := database.GetGormDB().Model(&User{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}

	return uint(count), nil
}

// UserExistByName checks if there is a tag with the same name
func UserExistByName(name string) (bool, error) {
	var user User
	err := database.GetGormDB().Select("id").Where("name = ?", name).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if user.ID > 0 {
		return true, nil
	}

	return false, nil
}

func UserDetailByName(name string) (*User, error) {
	var user User
	if err := database.GetGormDB().Where("name = ?", name).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func UserDetail(id uint) (*User, error) {
	var user User
	if err := database.GetGormDB().Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserPasswordByName(name string) (uint, string, error) {
	var user User
	if err := database.GetGormDB().Select("id, password").Where("name = ?", name).First(&user).Error; err != nil {
		return 0, "", err
	}
	return user.ID, user.Password, nil
}

// AddUser
func AddUser(name, password, role, email, gender string) error {
	user := User{
		Name:     name,
		Password: password,
		Role:     role,
		Email:    email,
		Gender:   gender,
	}
	db := database.GetGormDB()
	if err := db.Create(&user).Error; err != nil {
		return err
	}

	return nil
}
