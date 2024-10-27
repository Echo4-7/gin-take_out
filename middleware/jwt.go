package middleware

import (
	"Take_Out/pkg/e"
	"Take_Out/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		code = 200

		token := c.GetHeader("Authorization")
		if token == "" {
			code = 404
		} else {
			claims, err := util.ParseToken(token)
			if err != nil {
				code = e.ErrorAuthToken
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = e.ErrorAuthCheckTokenTimeout
			}
		}
		if code != e.SUCCESS {
			c.JSON(http.StatusOK, gin.H{
				"status ": code,
				"msg":     e.GetMsg(code),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
