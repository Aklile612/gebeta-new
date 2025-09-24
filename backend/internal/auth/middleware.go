package auth

import (
	"fmt"
	"net/http"

	"strings"

	"github.com/gin-gonic/gin"
)

type key int

const UserKey key =0
func JWTMiddleware()gin.HandlerFunc{
	return (func(c *gin.Context){
		authHeader:=c.GetHeader("Authorization")
		fmt.Println("Authorization Header:", authHeader)
		if authHeader == ""{
			c.JSON(http.StatusUnauthorized,gin.H{"error":"missing authorization header"})
			c.Abort()
			return 
		}
		parts:= strings.Split(authHeader," ")
		if len(parts)!= 2 || strings.ToLower(parts[0])!= "bearer"{
			c.JSON(http.StatusUnauthorized,gin.H{"error":"invalid authorization header"})
			c.Abort()
			return
		}

		tokenStr:= parts[1]
		claims,err:= ValidateJWT(tokenStr)
		if err!= nil{
			c.JSON(http.StatusUnauthorized,gin.H{"error":"invalid token"})
			c.Abort()
			return 
		}
		c.Set("user_id",claims.UserID)
		c.Set("email", claims.Email)     
		c.Set("full_name", claims.FullName)
		c.Set("jwt_token", tokenStr)
		c.Next()
		
	})
}