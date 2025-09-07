package middleware

import (
	"blogAggregator/internal/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc{
	return func(c *gin.Context){
		authHeader:=c.GetHeader("Authorization")
		if authHeader==""{
			c.JSON(http.StatusBadRequest,gin.H{
				"error":"missing token",
			})
			c.Abort()
			return
		}
		tokenStr:=strings.TrimPrefix(authHeader,"Bearer ")
		userID,err:=auth.ParseToken(tokenStr)
		if err!=nil{
			c.JSON(http.StatusUnauthorized,gin.H{
				"error":"invalid token",
			})
			c.Abort()
			return
		}
		c.Set("User_id",userID)
		c.Next()
	}
}