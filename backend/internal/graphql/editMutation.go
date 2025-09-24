package graphql

import (
	"context"

	"github.com/aklile/recipe-backend/internal/config"
	"github.com/aklile/recipe-backend/internal/models"
	hasura "github.com/machinebox/graphql"
)

func CheckRecipeOwnership(userID,recipeId,jwtToken string)(bool,error){

	// adminSecret:= config.LoadADMINSecret()

	req:= hasura.NewRequest(`
		query($user_id: uuid!, $recipe_id: uuid!) {
			recipes_by_pk(id: $recipe_id) {
				user_id
			}
		}
	`)
	req.Var("user_id",userID)
	req.Var("recipe_id",recipeId)
	req.Header.Set("Authorization", "Bearer "+jwtToken)


	var resp struct {
		RecipesByPk *struct {
			UserID string `json:"user_id"`
		} `json:"recipes_by_pk"`
	}


	err:= Client.Run(context.Background(),req,&resp)

	if err!= nil || resp.RecipesByPk==nil{
		return false,err
	}

	return resp.RecipesByPk.UserID == userID,nil

}

func UpdateRecipe(id, title, description, imageURL, difficulty string, prepTime int, cookTime int, categoryID string, isPaid bool, price float64,jwtToken string)error{
	// adminSecret:= config.LoadADMINSecret()

	req:= hasura.NewRequest(`
		mutation($id: uuid!, $title: String!, $description: String!, $image_url: String, $difficulty: String!,
		         $prep_time: Int!, $cook_time: Int!, $category_id: uuid!, $is_paid: Boolean!, $price: float8!) {
			update_recipes_by_pk(pk_columns: {id: $id}, _set: {
				title: $title,
				description: $description,
				featured_image: $image_url,
				difficulty: $difficulty,
				prep_time_minutes: $prep_time,
				cook_time_minutes: $cook_time,
				category_id: $category_id,
				is_paid: $is_paid,
				price: $price
			}) {
				id
			}
		}
	`)
	req.Var("id", id)
	req.Var("title", title)
	req.Var("description", description)
	req.Var("image_url", imageURL)
	req.Var("difficulty", difficulty)
	req.Var("prep_time", prepTime)
	req.Var("cook_time", cookTime)
	req.Var("category_id", categoryID)
	req.Var("is_paid", isPaid)
	req.Var("price", price)
	req.Header.Set("Authorization", "Bearer "+jwtToken)


	var resp struct {
		UpdateRecipesByPk struct {
			ID string `json:"id"`
		} `json:"update_recipes_by_pk"`
	}

	return Client.Run(context.Background(), req, &resp)
}

func UpdateRecipeSteps(recipeID string, steps []models.StepInput)error{

	adminSecret:= config.LoadADMINSecret()

	delReq:= hasura.NewRequest(`
		mutation(recipe_id: $uuid!){
			delete_recipe_steps(where:{recipe_id:{_eq: $recipe_id}}){
				affected_rows
			}
		}
	`)

	delReq.Var("recipe_id",recipeID)
	delReq.Header.Set("x-hasura-admin-secret",string(adminSecret))


	err:= Client.Run(context.Background(),delReq,nil)
	if err!= nil{
		return nil
	}

	return InsertRecipeSteps(recipeID,steps)
}

func UpdateIngredients(recipeID string, ingredients []models.IngredientInput) (error){
	adminSecret:= config.LoadADMINSecret()
	
	delReq:= hasura.NewRequest(`
		mutation(recipe_id: $uuid){
			delete_recipe_ingredients(where: {recipe_id:{_eq: $recipe_id}}){
				affected_rows
			}
		}
	`)

	delReq.Var("recipe_id",recipeID)
	delReq.Header.Set("x-hasura-admin-secret",string(adminSecret))

	err:= Client.Run(context.Background(),delReq,nil)

	if err!= nil{
		return err
	}

	return InsertIngredients(recipeID,ingredients)

}