package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gin-example/middleware/session-auth"
	"gin-example/pkg/app"
	"gin-example/pkg/errcode"
	"gin-example/service/users"
)

// curl -X GET "http://127.0.0.1:8000/api/v1/users"
// curl -X GET "http://127.0.0.1:8000/api/v1/users?name=zqyangchn"
type GetUsersForm struct {
	ID     uint   `form:"id" binding:""`
	Name   string `form:"name" binding:"max=100"`
	Role   string `form:"role" binding:""`
	Email  string `form:"email" binding:""`
	Gender string `form:"gender" binding:""`

	PageNumber int `form:"pageNumber,default=1" binding:"min=1"`
	PageSize   int `form:"pageSize,default=1" binding:"min=1,max=10"`
}

func GetUsers(c *gin.Context) {
	appG := app.Gin{Context: c}

	form := GetUsersForm{}
	if err := app.BindAndValid(c, &form); err != nil {
		appG.Response(http.StatusBadRequest, errcode.InvalidParamsError.WithDetails(err.Error()), struct{}{})
		return
	}

	users := userssvc.User{
		ID:     form.ID,
		Name:   form.Name,
		Role:   form.Role,
		Email:  form.Email,
		Gender: form.Gender,

		PageNumber: form.PageNumber,
		PageSize:   form.PageSize,
	}
	usersListResponse, err := users.GetUsers()
	if err != nil {
		appG.Response(http.StatusInternalServerError, errcode.GetUserError.WithDetails(err.Error()), nil)
		return
	}
	appG.ResponseSuccess(http.StatusOK, usersListResponse)
}

//
type AddUsersForm struct {
	Name     string `form:"name" binding:"required,min=3,max=100"`
	Password string `form:"password" binding:"required,min=6,max=20"`
	Role     string `form:"role" binding:"required"`
	Email    string `form:"email" binding:""`
	Gender   string `form:"gender" binding:""`
}

func Register(c *gin.Context) {
	var (
		appG = app.Gin{Context: c}
		form AddUsersForm
	)

	if err := app.BindAndValid(c, &form); err != nil {
		appG.Response(http.StatusBadRequest, errcode.InvalidParamsError.WithDetails(err.Error()), struct{}{})
		return
	}

	password, err := app.Encrypt(form.Password)
	if err != nil {
		appG.Response(http.StatusInternalServerError, errcode.CreateUserError.WithDetails(err.Error()), struct{}{})
		return
	}

	user := userssvc.User{
		Name:     form.Name,
		Password: password,
		Role:     form.Role,
		Email:    form.Email,
		Gender:   form.Gender,
	}
	exists, err := user.ExistByName()
	if err != nil {
		appG.Response(http.StatusInternalServerError, errcode.ServerError.WithDetails(err.Error()), struct{}{})
		return
	}
	if exists {
		appG.Response(http.StatusOK, errcode.CreateUserError.WithDetails("User Name Exist"), struct{}{})
		return
	}

	err = user.Add()
	if err != nil {
		appG.Response(http.StatusInternalServerError, errcode.ServerError.WithDetails(err.Error()), struct{}{})
		return
	}

	if err := sessionauth.SaveAuthSession(c, user.ID); err != nil {
		appG.Response(http.StatusInternalServerError, errcode.CreateSessionError.WithDetails(err.Error()), struct{}{})
		return
	}

	appG.Response(http.StatusOK, errcode.Success, struct{}{})
}

type LoginForm struct {
	Name     string `form:"name" binding:"required,min=3,max=100"`
	Password string `form:"password" binding:"required,min=6,max=20"`
}

func Login(c *gin.Context) {
	appG := app.Gin{Context: c}
	form := LoginForm{}

	if err := app.BindAndValid(c, &form); err != nil {
		appG.Response(http.StatusBadRequest, errcode.InvalidParamsError.WithDetails(err.Error()), struct{}{})
		return
	}

	if hasSession := sessionauth.HasSession(c); hasSession == true {
		appG.Response(http.StatusOK, errcode.Success, struct{}{})
		return
	}

	user := userssvc.User{
		Name:     form.Name,
		Password: form.Password,
	}
	if err := user.CheckPassword(); err != nil {
		appG.Response(http.StatusUnauthorized, errcode.UserPasswordError.WithDetails(err.Error()), struct{}{})
		return
	}

	if err := sessionauth.SaveAuthSession(c, user.ID); err != nil {
		appG.Response(http.StatusInternalServerError, errcode.CreateSessionError.WithDetails(err.Error()), struct{}{})
		return
	}

	appG.Response(http.StatusOK, errcode.Success, struct{}{})
}

func Logout(c *gin.Context) {
	appG := app.Gin{Context: c}

	if hasSession := sessionauth.HasSession(c); hasSession != true {
		appG.Response(http.StatusUnauthorized, errcode.ClearSessionError.WithDetails("用户未登录"), struct{}{})
		return
	}

	if err := sessionauth.ClearAuthSession(c); err != nil {
		appG.Response(http.StatusInternalServerError, errcode.ClearSessionError.WithDetails(err.Error()), struct{}{})
		return
	}

	appG.Response(http.StatusOK, errcode.Success, struct{}{})
}
