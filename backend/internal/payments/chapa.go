package payments

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aklile/recipe-backend/internal/config"
	"github.com/aklile/recipe-backend/internal/graphql"
	"github.com/gin-gonic/gin"
)

func InitRecipePurchaseHandler(c *gin.Context) {
	chapaSecret,chapaCallBack,chapaReturn:= config.CHAPAPAYMENTCREDINTIALS()
	recipeID := c.Param("id")

	
	userID := c.GetString("user_id")
	email := c.GetString("email")
	fullName := c.GetString("full_name")
	

	if userID == "" || email == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user info"})
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

	// Get recipe info
	recipe, err := graphql.GetRecipePurchaseInfo(recipeID, userID)
	if err != nil || recipe == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}
	if !recipe.IsPaid {
		c.JSON(http.StatusOK, gin.H{"access_granted": true, "message": "Free recipe"})
		return
	}

	txRef := fmt.Sprintf("recipe-%s-%s", recipeID, userID)

	payload := map[string]interface{}{
		"amount":       fmt.Sprintf("%.2f", recipe.Price),
		"currency":     "ETB",
		"email":        email,
		"first_name":   fullName,
		"last_name":    "",
		"tx_ref":       txRef,
		"callback_url": chapaCallBack,
		"return_url":   chapaReturn,
		"customization[title]":       "Recipe Purchase",
		"customization[description]": "Buying premium recipe",
	}

	bodyBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "https://api.chapa.co/v1/transaction/initialize", bytes.NewReader(bodyBytes))
	req.Header.Set("Authorization", "Bearer "+chapaSecret)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize payment"})
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var chapaResp map[string]interface{}
	json.Unmarshal(respBody, &chapaResp)

	if chapaResp["status"] == "success" {
		data := chapaResp["data"].(map[string]interface{})
		c.JSON(http.StatusOK, gin.H{"checkout_url": data["checkout_url"]})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get checkout link", "detail": chapaResp})
}

