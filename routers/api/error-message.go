package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gin-example/pkg/app"
	"gin-example/service"
)

// @Summary 获取错误码
// @Produce json
// @Success 200 {object} service.ErrorMessageResponse
// @Router /error/message [get]
func GetErrorMessages(c *gin.Context) {
	appG := app.Gin{Context: c}
	appG.ResponseSuccess(http.StatusOK, service.GetAllErrorMessage())
}
