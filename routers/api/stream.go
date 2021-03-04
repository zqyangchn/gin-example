package api

import (
	"gin-example/pkg/app"
	"gin-example/pkg/errcode"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// @Summary 权限验证, 获取 Token
// @Produce json
// @Param appKey body string false "appKey"
// @Param appSecret body string false "appSecret"
// @Success 200 {object} service.AuthResponse
// @Failure 500 {object} app.Response
// @Router /stream [get]
func Stream(c *gin.Context) {
	appG := app.Gin{Context: c,}
	i := 0
	for t := range time.Tick(1*time.Second) {
		if i > 5 {
			break
		}
		if _, err := c.Writer.WriteString(t.String()+"\n"); err != nil {
			appG.Response(http.StatusInternalServerError, errcode.EditTagError, struct {}{})
			return
		}
		c.Writer.Flush()
		i++
	}
	appG.ResponseSuccess(http.StatusOK, struct {}{})
}
