package jwtauth

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"gin-example/pkg/app"
	"gin-example/pkg/errcode"
)

// curl -X GET "http://127.0.0.1:8000/api/v1/tags?pageNumber=1&pageSize=10"
// curl -X GET "http://127.0.0.1:8000/api/v1/tags?pageNumber=1&pageSize=10" -H 'token:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhcHBLZXkiOiJ6cXlhbmdjaG4iLCJhcHBTZWNyZXQiOiJnby1wcm9tZ3JhbW1pbmctdG91ci1ib29rIiwiZXhwIjoxNTk4NTAwMzIwLCJpc3MiOiJodHRwLXNlcnZpY2UifQ.HxA0waRpdCpeiSK2P1qyFgBcz4O3kP_chrbF2UJ1oLY'

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			token string
			eCode = errcode.Success
		)

		if s, exist := c.GetQuery("token"); exist {
			token = s
		} else {
			token = c.GetHeader("token")
		}

		if token == "" {
			eCode = errcode.AuthTokenNotObtained
		} else {
			_, err := app.ParseToken(token)
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					eCode = errcode.AuthTokenTimeout
				default:
					eCode = errcode.AuthTokenError
				}
			}
		}

		if eCode.Code != errcode.Success.Code {
			appG := app.Gin{Context: c}
			appG.Response(http.StatusUnauthorized, eCode, struct{}{})
			c.Abort()
			return
		}
		c.Next()
	}
}
