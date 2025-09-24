package graphql

import (
	"context"
	"fmt"

	"github.com/aklile/recipe-backend/internal/config"
	
	 "github.com/machinebox/graphql"
)

func GetOrCreateCatagoryID(name string) (string, error) {
	adminSecret := config.LoadADMINSecret()

	checkQuery := graphql.NewRequest(`
		query($name:String!){
			categories(where:{name: {_eq:$name}}){
				id
			}
		}
	`)
	checkQuery.Var("name", name)
	checkQuery.Header.Set("x-hasura-admin-secret", string(adminSecret))

	var checkResp struct {
		Categories []struct {
			ID string `json:"id"`
		} `json:"categories"`
	}

	err := Client.Run(context.Background(), checkQuery, &checkResp)

	if err != nil {
		return "", fmt.Errorf("failed to check category: %w", err)
	}

	if len(checkResp.Categories) > 0 {
		// Exists
		return checkResp.Categories[0].ID, nil
	}

	insertReq:= graphql.NewRequest(`
		mutation($name:String!){
			insert_categories_one(object:{name:$name}){
				id
			}
		}
	`)

	insertReq.Var("name",name)
	insertReq.Header.Set("x-hasura-admin-secret", string(adminSecret))

	var insertResp struct {
		InsertCategory struct {
			ID string `json:"id"`
		} `json:"insert_categories_one"`
	}

	err = Client.Run(context.Background(),insertReq,&insertResp)

	if err != nil {
		return "", fmt.Errorf("failed to insert category: %w", err)
	}

	return insertResp.InsertCategory.ID,nil
}
