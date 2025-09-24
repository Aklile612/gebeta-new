package graphql

import (
	"context"
	"fmt"

	"github.com/aklile/recipe-backend/internal/config"
	"github.com/aklile/recipe-backend/internal/models"
	hasura "github.com/machinebox/graphql"
)


func GetUserByEmail(email string)(models.User,error){
	if Client == nil{
		return models.User{}, fmt.Errorf("GraphQL client not initialized")
	}
	adminSecret:= config.LoadADMINSecret()

	req:= hasura.NewRequest(`
		query($email:String!){
			users(where:{email:{_eq:$email}}){
				id
				email
				full_name
				password
			}
		}
	`)
	req.Var("email",email)
	req.Header.Set("x-hasura-admin-secret",string(adminSecret))

	var resp struct{
		Users  []models.User   `json:"users"`
	}

	err := Client.Run(context.Background(), req, &resp)

	if err != nil {
		fmt.Printf("❌ InsertUser GraphQL error: %v\n", err)
		return models.User{}, fmt.Errorf("insert user failed: %w", err)
	}

	if len(resp.Users)==0{
		return  models.User{}, fmt.Errorf("user not found")
	}

	return resp.Users[0],nil
}

func GetAllFullRecipes() ([]models.FullRecipe, error) {
	adminSecret := config.LoadADMINSecret()
	req := hasura.NewRequest(`
		query {
			 recipes {
  					  title
  					  description
  					  difficulty
  					  featured_image
  					  is_paid
  					  prep_time_minutes
  					  price
  					  USERNAME {
  					    full_name
  					  }
  					  commentfromuser {
  					    comment
  					    usercommented {
  					      full_name
  					    }
  					  }
  					  RECIPESTEPS {
  					    step_number
  					    description
  					  }
  					  RECIPECATAGORIES {
  					    name
  					  }
  					  RECIPEINGREDIENTS {
  					    name
  					    quantity
  					  }
  					}
		}
	`)

	req.Header.Set("x-hasura-admin-secret", string(adminSecret))

	var resp struct {
		Recipes []models.FullRecipe `json:"recipes"`
	}

	err := Client.Run(context.Background(), req, &resp)
	return resp.Recipes, err
}
