package graphql

import (
	"context"
	"fmt"
	"log"

	"github.com/aklile/recipe-backend/internal/config"
	hasura "github.com/machinebox/graphql"
)


type RecipePurchaseInfo struct {
	ID          string
	IsPaid      bool
	Price       float64
	OwnerUserID string
}

func InsertRecipePurchase(userID, recipeID string, amount float64) error{
	adminSecret:= config.LoadADMINSecret()

	req:= hasura.NewRequest(`
		mutation($user_id: uuid!, $recipe_id: uuid,$amount: numeric!){
			insert_recipe_purchases_one(object:{
				user_id: $user_id,
				recipe_id: $recipe_id,
				amount: $amount,
			}){
				id	
			}
		}
	`)
	req.Var("user_id",userID)
	req.Var("recipe_id",recipeID)
	req.Var("amount",amount)

	req.Header.Set("x-hasura-admin-secret", string(adminSecret))
	var resp struct{
		InsertRecipePurchasesOne struct{
			ID string `json:"id"`
		}`json:"insert_recipe_purchases_one"`
	}

	err:= Client.Run(context.Background(),req,&resp)

	if err != nil{
		return err
	}
	log.Println("Added to the recipe purshase table")
	return nil
}

func GetRecipePurchaseInfo(recipeID,requesterUserID string) (*RecipePurchaseInfo,error){

	adminSecret := config.LoadADMINSecret()

	req:= hasura.NewRequest(`
		query($id: uuid!){
			recipes_by_pk(id: $id){
				id
				is_paid
				price
				user_id
			}
		}
	`)
	req.Var("id", recipeID)
	req.Header.Set("x-hasura-admin-secret", string(adminSecret))


	var resp struct {
		RecipesByPk *struct {
			ID     string  `json:"id"`
			IsPaid bool    `json:"is_paid"`
			Price  float64 `json:"price"`
			UserID string  `json:"user_id"`
		} `json:"recipes_by_pk"`
	}

	err:= Client.Run(context.Background(),req,&resp)
	
	if err!= nil{

		return nil, fmt.Errorf("GetRecipePurchaseInfo gql error: %w", err)
	}

	if resp.RecipesByPk == nil {
		return nil, nil
	}
	return &RecipePurchaseInfo{
		ID:          resp.RecipesByPk.ID,
		IsPaid:      resp.RecipesByPk.IsPaid,
		Price:       resp.RecipesByPk.Price,
		OwnerUserID: resp.RecipesByPk.UserID,
	}, nil
}

func UserHasRecipeAccess(userID,recipeID string) (bool,error){
	adminSecret := config.LoadADMINSecret()

	req:= hasura.NewRequest(`
		query ($user_id: uuid!, $recipe_id: uuid!) {
  			recipe_access_aggregate(
  			  where: { user_id: { _eq: $user_id }, recipe_id: { _eq: $recipe_id } }
  			) {
  			  aggregate {
  			    count
  			  }
  			}
		}
	`)

	req.Var("user_id", userID)
	req.Var("recipe_id", recipeID)
	req.Header.Set("x-hasura-admin-secret", string(adminSecret))

	var resp struct {
		RecipeAccessAggregate struct {
			Aggregate struct {
				Count int `json:"count"`
			} `json:"aggregate"`
		} `json:"recipe_access_aggregate"`
	}
	err := Client.Run(context.Background(), req, &resp)

	if err!= nil{
		return false,err
	}

	return resp.RecipeAccessAggregate.Aggregate.Count > 0, nil

}