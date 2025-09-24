package payments

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/aklile/recipe-backend/internal/config"
	"github.com/aklile/recipe-backend/internal/graphql"
	"github.com/gin-gonic/gin"
)

func ChapaCallbackHandler(c *gin.Context) {
	chapaSecret, _, _ := config.CHAPAPAYMENTCREDINTIALS()

	// Get tx_ref from query
	txRef := c.Query("tx_ref")
	if  strings.TrimSpace(txRef) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or empty tx_ref"})
		return
	}

	// Call Chapa's verification API
	verifyURL := fmt.Sprintf("https://api.chapa.co/v1/transaction/verify/%s", txRef)
	req, err := http.NewRequest("GET", verifyURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create verification request"})
		return
	}
	req.Header.Set("Authorization", "Bearer "+chapaSecret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Verification request failed"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read verification response"})
		return
	}

	// Parse JSON response
	var verifyResp map[string]interface{}
	if err := json.Unmarshal(body, &verifyResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
		return
	}

	// Defensive checks
	if verifyResp["status"] == "success" {
		dataInterface, ok := verifyResp["data"]
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No data in verification response"})
			return
		}

		data, ok := dataInterface.(map[string]interface{})
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid data format"})
			return
		}

		if data["status"] == "success" {
			if !ok || data["status"] != "success" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Payment not verified"})
				return
			}
			txRefStr, ok := data["tx_ref"].(string)
			if !ok {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tx_ref in response"})
				return
			}

			parts := strings.Split(txRefStr, "-")
			if len(parts) != 3 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tx_ref format"})
				return
			}

			recipeID := parts[1]
			userID := parts[2]



			recipe, err := graphql.GetRecipePurchaseInfo(recipeID, userID)
			if err !=nil{
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch recipe purchase info", "detail": err.Error()})
				return
			}
			if recipe == nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
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
			// Insert record of purchase
			err =graphql.InsertRecipePurchase(userID, recipeID, recipe.Price)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record purchase"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "Payment verified. Access granted."})
			return
		}
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "Payment not verified"})
}
