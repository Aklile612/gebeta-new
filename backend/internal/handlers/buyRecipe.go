package handlers

import (
	"net/http"

	"github.com/aklile/recipe-backend/internal/graphql"
	"github.com/gin-gonic/gin"
)


func BuyRecipeHandler(c *gin.Context){
	recipeID:= c.Param("id")

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

	recipe,err:= graphql.GetRecipePurchaseInfo(recipeID,userID)

	if err !=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch recipe purchase info", "detail": err.Error()})
		return
	}

	if recipe==nil{
		c.JSON(http.StatusNotFound,gin.H{"error":"recipe not found"})
		return
	}
	if !recipe.IsPaid{
		c.JSON(http.StatusOK, gin.H{
			"message":        "recipe is not paid; access granted",
			"access_granted": true,
		})
		return
	}

	hasAccess,err:= graphql.UserHasRecipeAccess(userID,recipeID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check access", "detail": err.Error()})
		return
	}

	if hasAccess{
		c.JSON(http.StatusOK, gin.H{
			"message":        "already purchased; access granted",
			"access_granted": true,
		})
		return
	}
	

	c.JSON(http.StatusOK,gin.H{"message":"payment required",
		
		"access_granted":      false,})
}