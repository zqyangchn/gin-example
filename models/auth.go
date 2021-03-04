package models

import (
	"gorm.io/gorm"

	"gin-example/pkg/database"
)

type Auth struct {
	gorm.Model

	AppKey    string `json:"app_key"`
	AppSecret string `json:"app_secret"`
}

func (a *Auth) Get() (Auth, error) {
	var auth Auth

	err := database.GetGormDB().Where("app_key = ? AND app_secret = ? AND is_del = ?", a.AppKey, a.AppSecret, 0).First(&auth).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return auth, err
	}
	return auth, nil
}
