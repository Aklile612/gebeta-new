package handlers

import (
	"net/http"

	"github.com/aklile/recipe-backend/internal/auth"
	"github.com/aklile/recipe-backend/internal/graphql"
	"github.com/gin-gonic/gin"
)



type RegisterRequest struct{
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func RegisterHandler(c *gin.Context){
	var req RegisterRequest

	err:= c.ShouldBindJSON(&req)
	
	if err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":"Invalid Input"})
		return
	}
	hashedPassword,err:= auth.HashPassword(req.Password)

	if err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error": "Error Hashing Password"})
		return
	}

	user,err:= graphql.InsertUser(req.Name,req.Email,hashedPassword)

	if err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{"Error": "Error saveing the user"})
		return
	}

	token,err:= auth.GenerateJWT(user.ID,user.Email,user.FullName)

	if err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{"Error": "Error geerating token the user"})
	}

	c.JSON(http.StatusOK,gin.H{"message":"user successfully created","user":user,"token":token})

}
