package handlers

import (
	"log"
	"net/http"

	"github.com/aklile/recipe-backend/internal/graphql"
	"github.com/gin-gonic/gin"
)

type PurchaseEventPayload struct {
	Event struct {
		Op   string `json:"op"`
		Data struct {
			New struct {
				ID          string `json:"id"`
				UserID      string `json:"user_id"`
				RecipeID    string `json:"recipe_id"`
				PurchasedAt string `json:"purchased_at"`
				Amount      string `json:"amount"`
			} `json:"new"`
		} `json:"data"`
	} `json:"event"`
}

func PurchaseWebhookHandler(c *gin.Context){

	var payload PurchaseEventPayload
	err:= c.BindJSON(&payload)

	if err!= nil{
		log.Printf("Error parsing purchase webhook: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	if payload.Event.Op != "INSERT"{
		c.JSON(http.StatusBadRequest, gin.H{"error": "not an insert event"})
		return
	}

	purchase:= payload.Event.Data.New

	err = graphql.GrantAccess(purchase.UserID,purchase.RecipeID)
	if err != nil{
		log.Printf("Failed to grant access: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to grant access"})
		return
	}
	c.JSON(http.StatusOK,gin.H{"message":"access granted"})
}
