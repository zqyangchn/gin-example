package app

import (
	"github.com/gin-gonic/gin"

	"gin-example/pkg/errcode"
)

type Gin struct {
	Context *gin.Context
}

type Response struct {
	*errcode.ErrorMessage
	Data interface{} `json:"data"`
}

// Response setting gin.JSON
func (g *Gin) Response(httpCode int, eMsg *errcode.ErrorMessage, data interface{}) {
	g.Context.JSON(httpCode, Response{
		ErrorMessage: eMsg,
		Data:         data,
	})
	return
}

func (g *Gin) ResponseSuccess(httpCode int, data interface{}) {
	g.Context.JSON(httpCode, data)
	return
}
