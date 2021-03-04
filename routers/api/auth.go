package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gin-example/pkg/app"
	"gin-example/pkg/errcode"
	"gin-example/service"
)

type AuthForm struct {
	AuthKey    string `form:"appKey" binding:"required"`
	AuthSecret string `form:"appSecret" binding:"required"`
}

// @Summary 权限验证, 获取 Token
// @Produce json
// @Param appKey body string false "appKey"
// @Param appSecret body string false "appSecret"
// @Success 200 {object} service.AuthResponse
// @Failure 500 {object} app.Response
// @Router /auth [post]
func GetAuth(c *gin.Context) {
	appG := app.Gin{Context: c}
	form := AuthForm{}

	if err := app.BindAndValid(c, &form); err != nil {
		appG.Response(http.StatusBadRequest, errcode.InvalidParamsError.WithDetails(err.Error()), struct{}{})
		return
	}

	if err := service.CheckAuth(form.AuthKey, form.AuthSecret); err != nil {
		appG.Response(http.StatusUnauthorized, errcode.AuthNotExistError.WithDetails(err.Error()), struct{}{})
		return
	}

	authResponse, err := service.GenerateToken(form.AuthKey, form.AuthSecret)
	if err != nil {
		appG.Response(http.StatusInternalServerError, errcode.AuthTokenGenerateError.WithDetails(err.Error()), struct{}{})
		return
	}
	appG.ResponseSuccess(http.StatusOK, authResponse)
}

type ReverseSolutionJWTForm struct {
	Token string `form:"token" binding:"required"`
}

// curl -X GET "http://127.0.0.1:8000/reverse/solution/jwt?token=token"

// @Summary 解析Token成json格式
// @Produce json
// @Param token query string false "token"
// @Success 200 {object} service.ReverseSolutionJWTResponse
// @Failure 500 {object} app.Response
// @Router /reverse/solution/jwt [get]
func ReverseSolutionJWT(c *gin.Context) {
	appG := app.Gin{Context: c}
	form := ReverseSolutionJWTForm{}

	if err := app.BindAndValid(c, &form); err != nil {
		appG.Response(http.StatusBadRequest, errcode.InvalidParamsError.WithDetails(err.Error()), struct{}{})
		return
	}

	reverseSolutionJWTResponse, err := service.ReverseSolutionJWT(form.Token)
	if err != nil {
		appG.Response(http.StatusInternalServerError, errcode.AuthTokenParseError.WithDetails(err.Error()), struct{}{})
		return
	}
	appG.ResponseSuccess(http.StatusOK, reverseSolutionJWTResponse)
}
