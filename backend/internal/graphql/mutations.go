	package graphql

	import (
		"context"
		"fmt"
		

		"github.com/aklile/recipe-backend/internal/config"
		
		"github.com/aklile/recipe-backend/internal/models"
		hasura "github.com/machinebox/graphql"
	)
	func InsertUser(name, email, password string) (models.User, error) {
		fmt.Printf("InsertUser called. Client: %v\n", Client)
		
		if Client == nil {
			return models.User{}, fmt.Errorf("graphql client not initialized")
		}
		var adminSecret = config.LoadADMINSecret()
		
		fmt.Printf("Admin Secret: %s\n", string(adminSecret))

		req := hasura.NewRequest(`
			mutation($full_name:String!,$email:String!,$password:String!){
				insert_users_one(object:{
				full_name:$full_name,
				email:$email,
				password:$password
				}){
				id
				email
				full_name
				}
			}
		`)
		req.Var("full_name",name)
		req.Var("email",email)
		req.Var("password",password)
		req.Header.Set("x-hasura-admin-secret",string(adminSecret))

		var resp struct{
			InsertUser models.User  `json:"insert_users_one"`

		}
		err:= Client.Run(context.Background(),req,&resp)
		if err != nil {
			fmt.Printf("❌ InsertUser GraphQL error: %v\n", err)
			return models.User{}, fmt.Errorf("insert user failed: %w", err)
		}

		return resp.InsertUser, nil

	}
	func InsertRecipe(title,desc,imageURL, difficulty string, prepTime, cookTime int, userID, categoryID string, isPaid bool, price float64)(models.Recipe,error){

		var adminSecret = config.LoadADMINSecret()
		req:= hasura.NewRequest(`
			mutation($title: String!, $description: String!, $featured_image: String!, $difficulty: String!, $prep_time_minutes: Int!, $cook_time_minutes: Int!, $user_id: uuid!, $category_id: uuid!, $is_paid: Boolean!, $price: numeric!){
				insert_recipes_one(object:{
					title: $title,
					description: $description,
					featured_image: $featured_image,
					difficulty: $difficulty,
					prep_time_minutes: $prep_time_minutes,
					cook_time_minutes: $cook_time_minutes,
					user_id: $user_id,
					category_id: $category_id,
					is_paid: $is_paid,
					price: $price
				}){
					id
					title
					featured_image	
				}
			}
		`)
		req.Var("title", title)
		req.Var("description", desc)
		req.Var("featured_image", imageURL)
		req.Var("difficulty", difficulty)
		req.Var("prep_time_minutes", prepTime)
		req.Var("cook_time_minutes", cookTime)
		req.Var("user_id", userID)
		req.Var("category_id", categoryID)
		req.Var("is_paid", isPaid)
		req.Var("price", price)
		req.Header.Set("x-hasura-admin-secret", string(adminSecret))

		var resp struct{
			InsertRecipe models.Recipe `json:"insert_recipes_one"`
		}

		err := Client.Run(context.Background(),req,&resp)

		return resp.InsertRecipe,err
	}
	func InsertRecipeSteps(recipeID string, steps []models.StepInput) error{
		adminSecret:= config.LoadADMINSecret()

		var inputSteps []map [string]interface{}
		for _,step := range steps{
			inputSteps = append(inputSteps, map[string]interface{}{
				"recipe_id":   recipeID,
            	"step_number": step.StepNumber,
            	"description": step.Description,
			})

		}
		fmt.Printf("Steps to insert: %+v\n", inputSteps)
		req:= hasura.NewRequest(`
			mutation($steps:[recipe_steps_insert_input!]!){
				insert_recipe_steps(objects:$steps){
					affected_rows
				}
			}
		`)

		req.Var("steps",inputSteps)
		req.Header.Set("x-hasura-admin-secret", string(adminSecret))
		var resp struct {
        	InsertRecipeSteps struct {
        	    AffectedRows int `json:"affected_rows"`
        	} `json:"insert_recipe_steps"`
    	}

		return Client.Run(context.Background(), req, &resp)
	}

	func InserComment(userID, recipeID , comment string) error {
		adminSecret:= config.LoadADMINSecret()	

		req:= hasura.NewRequest(`
			mutation($user_id: uuid!,$recipe_id: uuid!,$comment: String!){
				insert_recipe_comments_one(object:{
					user_id:$user_id,
					recipe_id:$recipe_id,
					comment:$comment
				}){
					id	
				}
			}
		`)

		req.Var("user_id",userID)
		req.Var("recipe_id",recipeID)
		req.Var("comment",comment)
		req.Header.Set("x-hasura-admin-secret", string(adminSecret))
		
		var resp struct{
			InsertComment struct{
				ID string `json:"id"`
			} `json:"insert_recipe_comments_one"`
		}
		return  Client.Run(context.Background(),req,&resp)
	}

	func InsertIngredients(recipeID string, ingredients []models.IngredientInput) error{
		adminSecret:= config.LoadADMINSecret()

		var inputIngredients []map [string]interface{}

		for _,ingredient := range ingredients{
			inputIngredients = append(inputIngredients, map[string]interface{}{
				"recipe_id":recipeID,
				"name":ingredient.Name,
				"quantity":ingredient.Quantity,
			})
		}
		req:= hasura.NewRequest(`
			mutation ($objects: [ingredients_insert_input!]!) {
				insert_ingredients(objects: $objects) {
					affected_rows
				}
			}	
		`)

		req.Var("objects",inputIngredients)
		req.Header.Set("x-hasura-admin-secret", string(adminSecret))

		var resp struct{
			InsertIngredients struct{
				AffectedRows int `json:"affected_rows"`
			} `json:"insert_ingredients"`
		}

		return  Client.Run(context.Background(),req,&resp)


	}


	func InesrtRecipeLikes(userID,recipeID string)error{
		adminSecret:= config.LoadADMINSecret()

		req:= hasura.NewRequest(`
			mutation($user_id: uuid!, $recipe_id: uuid!){
				insert_recipe_likes_one(object:{
					user_id: $user_id,
					recipe_id: $recipe_id,
				}){
					user_id	
				}
			}
		`)

		req.Var("user_id",userID)
		req.Var("recipe_id",recipeID)
		req.Header.Set("x-hasura-admin-secret", string(adminSecret))

		var resp struct{
			InsertLike struct{
				ID string `json:"id"`
			} `json:"insert_recipe_likes_one"`
		}

		return Client.Run(context.Background(),req,&resp)
	}

	func InsertUserRecipeBookmark(userID,recipeID string) error{
		adminSecret:= config.LoadADMINSecret()

		req:= hasura.NewRequest(`
			mutation($user_id: uuid!, $recipe_id: uuid!){
				insert_recipe_bookmarks_one(object:{
					user_id: $user_id,
					recipe_id: $recipe_id
				}){
					user_id	
				}
			}
		`)

		req.Var("user_id",userID)
		req.Var("recipe_id",recipeID)
		req.Header.Set("x-hasura-admin-secret", string(adminSecret))

		var resp struct{
			InsertBookmark struct{
				ID string `json:"id"`
			} `json:"insert_recipe_bookmarks_one"`
		}

		return  Client.Run(context.Background(),req,&resp)


	}


	func InsertRatingByUser(userID string,recipeID string,rating int) error{
		adminSecret:= config.LoadADMINSecret()

		req:= hasura.NewRequest(`
			mutation($user_id: uuid!, $recipe_id: uuid!, $rating: Int!){
				insert_recipe_ratings_one(
					object: {
					user_id: $user_id,
					recipe_id: $recipe_id,
					rating: $rating
					},
					on_conflict: {
						constraint: recipe_ratings_pkey,
						update_columns: [rating]
					}
				){
					recipe_id	
				}
			}
		`)

		req.Var("user_id", userID)
		req.Var("recipe_id", recipeID)
		req.Var("rating", rating)
		req.Header.Set("x-hasura-admin-secret", string(adminSecret))

		var resp struct {
			InsertRecipeRatingsOne struct {
				RecipeID string `json:"recipe_id"`
			} `json:"insert_recipe_ratings_one"`
		}

		return Client.Run(context.Background(),req,&resp)
	}

	func DeleteRecipeLike(userID , recipeID string) error{

		adminSecret:= config.LoadADMINSecret()

		req:= hasura.NewRequest(`
			mutation($user_id: uuid! , $recipe_id: uuid!){
				delete_recipe_likes(where: {
					user_id: {_eq: $user_id},
					recipe_id: {_eq: $recipe_id}
				}){
					affected_rows	
				}
			}
		`)

		req.Var("user_id",userID)
		req.Var("recipe_id",recipeID)
		req.Header.Set("x-hasura-admin-secret", string(adminSecret))

		var resp struct{

			DeleteRecipeLikes struct{
				AffectedRows int `json:"affected_rows"`
			}`json:"delete_recipe_likes"`
		}

		return Client.Run(context.Background(),req,&resp)
	}
	func InsertRecipeImages(recipeID string, imageURLs []string) error {
		adminSecret := config.LoadADMINSecret()
	
		
		var objects []map[string]interface{}
		for _, url := range imageURLs {
			objects = append(objects, map[string]interface{}{
				"recipe_id": recipeID,
				"image_url": url,
			})
		}
	
		req := hasura.NewRequest(`
			mutation($objects: [recipe_images_insert_input!]!) {
				insert_recipe_images(objects: $objects) {
					affected_rows
				}
			}
		`)
	
		req.Var("objects", objects)
		req.Header.Set("x-hasura-admin-secret", string(adminSecret))
	
		var resp struct {
			Insert struct {
				AffectedRows int `json:"affected_rows"`
			} `json:"insert_recipe_images"`
		}
	
		return Client.Run(context.Background(), req, &resp)
	}
	