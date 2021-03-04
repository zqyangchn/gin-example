package service

import (
	"errors"
	"fmt"

	"gin-example/models"
	"gin-example/pkg/app"
	"gin-example/pkg/errcode"
)

type Token struct {
	Token string
}

type AuthResponse struct {
	*errcode.ErrorMessage
	Data Token
}

func CheckAuth(appKey, appSecret string) error {
	auth := models.Auth{AppKey: appKey, AppSecret: appSecret}
	auth, err := auth.Get()
	if err != nil {
		return err
	}
	if auth.ID > 0 {
		return nil
	}
	return errors.New("auth info not exist")
}

func GenerateToken(appKey, appSecret string) (*AuthResponse, error) {
	token, err := app.GenerateToken(appKey, appSecret)
	fmt.Println(err)
	if err != nil {
		return nil, err
	}
	return &AuthResponse{
		ErrorMessage: errcode.Success,
		Data:         Token{Token: token},
	}, nil
}

type ReverseSolutionJWTResponse struct {
	*errcode.ErrorMessage
	Data *app.Claims
}

func ReverseSolutionJWT(token string) (*ReverseSolutionJWTResponse, error) {
	claims, err := app.ParseTokenWithoutValid(token)
	if err != nil {
		return nil, err
	}

	return &ReverseSolutionJWTResponse{
		ErrorMessage: errcode.Success,
		Data:         claims,
	}, nil
}
