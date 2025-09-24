package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aklile/recipe-backend/internal/auth"
	"github.com/aklile/recipe-backend/internal/config"
	"github.com/aklile/recipe-backend/internal/graphql"
	"github.com/aklile/recipe-backend/internal/handlers"
	"github.com/aklile/recipe-backend/internal/payments"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//docker start docker-hasura-1 docker-postgres-1 portainer
func main() {
	
	config.LoadEnv()
	fmt.Println("END_POINT from os.Getenv:", os.Getenv("END_POINT"))
	endpoint := string(config.ENDPoint())
	fmt.Println("Loaded END_POINT from config.ENDPoint():", endpoint)

	if endpoint == "" {
		log.Fatal("❌ END_POINT is empty. Please check your .env and config.")
	}

	graphql.InitClient(endpoint)
	if graphql.Client == nil {
		log.Fatal("❌ GraphQL client is still nil after InitClient")
	} else {
		fmt.Println("✅ GraphQL client initialized correctly")
	}
	// go build -o servercmd && ./servercmd 
	//du -sh ~/.cache/go-build
	//sudo du -sh /var/log
	// sudo journalctl --vacuum-time=3d
	//go clean -cache -testcache -modcache

	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001", "http://127.0.0.1:3000", "http://127.0.0.1:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.POST("/register", handlers.RegisterHandler)
	router.POST("/login",handlers.LoginHandler)
	router.GET("/recipes", handlers.GetAllRecipesHandler)
	router.POST("/recipes/webhook/purchase",handlers.PurchaseWebhookHandler)
	router.GET("/api/payments/callback", payments.ChapaCallbackHandler)
	authGroup:= router.Group("/")
	authGroup.Use(auth.JWTMiddleware())
	authGroup.POST("/add_recipes",handlers.AddRecipeHandler)
	authGroup.PUT("/recipes/:id",handlers.EditRecipesHandler)
	authGroup.POST("/comment_recipes/:id",handlers.AddCommentHandler)
	authGroup.POST("/recipes/likes/:id",handlers.AddRecipeLikesHandler)
	authGroup.POST("/recipes/bookmark/:id",handlers.ADDRecipeBookmarksHandler)
	authGroup.POST("/recipes/rating/:id",handlers.AddRatingtoRecipeHandler)
	authGroup.DELETE("/recipes/likes/:id",handlers.DeleteRecipeLikehandler)
	authGroup.POST("recipes/buy/:id",payments.InitRecipePurchaseHandler)
	err := router.Run(":8081")

	if err != nil { 
		log.Fatal("Failed to run the server", err)
	}
}