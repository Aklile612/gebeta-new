package handlers

import (
	"encoding/json"
	"fmt"

	"net/http"
	"strconv"

	"github.com/aklile/recipe-backend/internal/graphql"
	"github.com/aklile/recipe-backend/internal/media"
	"github.com/aklile/recipe-backend/internal/models"
	"github.com/gin-gonic/gin"
)


func AddRecipeHandler(c *gin.Context){
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
	title:= c.PostForm("title")
	description := c.PostForm("description")
	prepTimeStr := c.PostForm("prep_time_minutes")
	cookTimeStr := c.PostForm("cook_time_minutes")
	difficulty := c.PostForm("difficulty")
	
	categoryName := c.PostForm("category_name")
	categoryID, err := graphql.GetOrCreateCatagoryID(categoryName)
if err != nil {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get/create category"})
	return
}
	isPaid := c.PostForm("is_paid") == "true"
	priceStr := c.PostForm("price")

	form , err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form", "detail": err.Error()})
		return
	}
	files:= form.File["image"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No images provided"})
		return
	}
	
	var imageURLs []string
	for _,fileHeader := range files{
		file,err:=fileHeader.Open()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open image", "detail": err.Error()})
			return
		}
		defer file.Close()

		url, err := media.UploadImage(file, fileHeader)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image", "detail": err.Error()})
			return
		}

		imageURLs = append(imageURLs, url)

	}

	if len(imageURLs)==0{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"No image found"})
		return
	}
	
	featuredImage := imageURLs[0]

	prepTime,_:= strconv.Atoi(prepTimeStr)
	cookTime,_:= strconv.Atoi(cookTimeStr)
	price,_ := strconv.ParseFloat(priceStr,64)


	recipe,err:= graphql.InsertRecipe(title, description, featuredImage, difficulty, prepTime, cookTime, userID, categoryID, isPaid, price)

	if err!= nil{
		 fmt.Println("❌ InsertRecipe error:", err)
    	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    	return
	}

	err = graphql.InsertRecipeImages(recipe.ID,imageURLs)

	if err!= nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"can't save the images"})
		return
	}
	stepsJson:=c.PostForm("steps")
	
	var steps []models.StepInput

	err = json.Unmarshal([]byte(stepsJson),&steps)

	if err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":"Invalid steps format","detail":err.Error()})
		return
	}
	err= graphql.InsertRecipeSteps(recipe.ID,steps)
	if err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"failed to save Steps","details":err.Error()})
		return
	}
	ingredientsJSON:= c.PostForm("ingredients")
	var ingredients []models.IngredientInput

	err= json.Unmarshal([]byte(ingredientsJSON),&ingredients)

	if err!= nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":"Invalid ingredients format format","detail":err.Error()})
		return
	}
	err=graphql.InsertIngredients(recipe.ID,ingredients)
	if err !=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"failed to save ingredients","details":err.Error()})
		return
	}

	if isPaid{
		err= graphql.InsertRecipePurchase(userID,recipe.ID,price)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":  "Failed to record initial recipe purchase",
				"detail": err.Error(),
			})
			return
		}
	}
	c.JSON(http.StatusOK,gin.H{"recipe":recipe,"steps":steps,"ingredients":ingredients,"imageURLs":imageURLs})
}

func GetAllRecipesHandler(c *gin.Context) {
	recipes, err := graphql.GetAllFullRecipes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recipes", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"recipes": recipes})
}
