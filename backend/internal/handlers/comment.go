package handlers

import (
	"net/http"

	"github.com/aklile/recipe-backend/internal/graphql"
	"github.com/gin-gonic/gin"
)

func AddCommentHandler(c *gin.Context) {

	recipeID := c.Param("id")

	userIDInterf, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authorized"})
		return
	}

	userID, ok := userIDInterf.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user id"})
		return
	}

	comment := c.PostForm("comment")

	if comment == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "comment can't be empty"})
		return
	}

	err := graphql.InserComment(userID, recipeID, comment)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error creating comment", "details": err.Error()})
		return
	}

	 c.JSON(http.StatusOK, gin.H{"message": "comment succesfully created"})
}
