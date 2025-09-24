package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aklile/recipe-backend/internal/graphql"
	"github.com/aklile/recipe-backend/internal/media"
	"github.com/aklile/recipe-backend/internal/models"
	"github.com/gin-gonic/gin"
)


func EditRecipesHandler(c *gin.Context){
	recipeID := c.Param("id")

	userIDInterf, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}
	userID, ok := userIDInterf.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}	

	jwtTokenInterf, exists := c.Get("jwt_token")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "JWT token not found"})
		return
	}
	jwtToken, ok := jwtTokenInterf.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid JWT token"})
		return
	}

	isOwner,err:= graphql.CheckRecipeOwnership(userID,recipeID,jwtToken)

	if err!=nil || !isOwner{
		c.JSON(http.StatusForbidden,gin.H{"error":"Didn't have access to edit"})
		return
	}

	title := c.PostForm("title")
	description := c.PostForm("description")
	prepTimeStr := c.PostForm("prep_time_minutes")
	cookTimeStr := c.PostForm("cook_time_minutes")
	difficulty := c.PostForm("difficulty")
	categoryName := c.PostForm("category_name")
	isPaid := c.PostForm("is_paid") == "true"
	priceStr := c.PostForm("price")


	categoryID, err := graphql.GetOrCreateCatagoryID(categoryName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get/create category"})
		return
	}

	file,fileHeader,err:= c.Request.FormFile("image")
	if err!= nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":"error uploading image"})
		return
	}
	defer file.Close()
	imageURL,err:= media.UploadImage(file,fileHeader)
	if err!= nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":"Failed to upload image","detail": err.Error()})
		return
	}

	prepTime, _ := strconv.Atoi(prepTimeStr)
	cookTime, _ := strconv.Atoi(cookTimeStr)
	price, _ := strconv.ParseFloat(priceStr, 64)

	err = graphql.UpdateRecipe(recipeID, title, description, imageURL, difficulty, prepTime, cookTime, categoryID, isPaid, price,jwtToken)

	if err!= nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to update a recipe"})
		return
	}

	stepsJson:=c.PostForm("steps")

	var steps []models.StepInput
	err = json.Unmarshal([]byte(stepsJson), &steps)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid steps format", "detail": err.Error()})
		return
	}

	err= graphql.UpdateRecipeSteps(recipeID,steps)

	if err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"internal server error","detail":err.Error()})
		return
	}

	
	ingredientsJSON := c.PostForm("ingredients")
	var ingredients []models.IngredientInput
	err= json.Unmarshal([]byte(ingredientsJSON),&ingredients)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ingredients format", "detail": err.Error()})
		return
	}

	err = graphql.UpdateIngredients(recipeID,ingredients)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ingredients", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recipe updated successfully"})
}