package handlers

import (
	
	"net/http"

	"github.com/aklile/recipe-backend/internal/auth"
	"github.com/aklile/recipe-backend/internal/graphql"
	"github.com/gin-gonic/gin"
)


type LoginRequest struct{
	Email string  `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct{

	Token  string `json:"token"`
	Error string `json:"error,omitempty"`
}

func LoginHandler(c *gin.Context){
	var req LoginRequest
	err:= c.ShouldBindJSON(&req)

	if err!= nil{
		c.JSON(http.StatusBadRequest,gin.H{"Error":"Bad request"})
		return
	}

	user,err:= graphql.GetUserByEmail(req.Email)

	if err!= nil{
		c.JSON(http.StatusUnauthorized,gin.H{"error":"invalid credentials"})
		return
	}

	if !auth.CheckHashPassword(req.Password,user.Password){
		c.JSON(http.StatusUnauthorized,gin.H{"Error":"invalid credientials"})
		return
	}

	token,err:= auth.GenerateJWT(user.ID,user.Email,user.FullName)

	if err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"failed to generate token"})
		return
	}
	user.Password=""
	c.JSON(http.StatusOK,gin.H{"Token":token,"messgae":"user succesfully loged it","user":user})
	

}