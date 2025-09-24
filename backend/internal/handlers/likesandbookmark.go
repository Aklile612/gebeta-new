package handlers

import (
	"net/http"

	"github.com/aklile/recipe-backend/internal/graphql"
	"github.com/gin-gonic/gin"
)

func AddRecipeLikesHandler(c *gin.Context){
	recipeID := c.Param("id")

	userIDInterf, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}
	userID, ok := userIDInterf.(string)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id not valied"})
		return
	}
	err:= graphql.InesrtRecipeLikes(userID,recipeID)

	if err!= nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":"cannot like the recipe","detail":err.Error()})
		return
	}

	c.JSON(http.StatusOK,gin.H{"message":"the recipe liked succefully"})

}


func ADDRecipeBookmarksHandler(c *gin.Context){
	recipeID:= c.Param("id")

	userIDInterf,exists:= c.Get("user_id")

	if !exists{
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	userID ,ok := userIDInterf.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id not valied"})
		return
	}

	err:= graphql.InsertUserRecipeBookmark(userID,recipeID)
	if err!= nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":"cannot like the recipe","detail":err.Error()})
		return
	}

	c.JSON(http.StatusOK,gin.H{"message":"Bookmard Added succefully"})

}

func AddRatingtoRecipeHandler(c *gin.Context){
	recipeID := c.Param("id")

	userIDInterf, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}
	userID, ok := userIDInterf.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID"})
		return
	}

	var body struct{
		Rating int `json:"rating"`
	}

	err:= c.BindJSON(&body)

	if err != nil || body.Rating < 1 || body.Rating > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "rating must be between 1 and 5"})
		return
	}

	err = graphql.InsertRatingByUser(userID,recipeID,body.Rating)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not rate recipe"})
		return
	}

	c.JSON(http.StatusOK,gin.H{"message":"recpe succesfully rated"})
}

func DeleteRecipeLikehandler(c *gin.Context){
	recipeID := c.Param("id")

	userIDInterf, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}
	userID, ok := userIDInterf.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID"})
		return
	}

	err := graphql.DeleteRecipeLike(userID,recipeID)
	if err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"can't unlike the recipe"})
		return
	}

	c.JSON(http.StatusOK,gin.H{"message":"succesfully unliked recipe"})
}