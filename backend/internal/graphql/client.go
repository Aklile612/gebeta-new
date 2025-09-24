package graphql

import (
    "fmt"
    
    "sync"

    "github.com/machinebox/graphql"
)

var (
    Client *graphql.Client
    once   sync.Once
)

func InitClient(endpoint string) {
    once.Do(func() {
        fmt.Printf("GraphQL Client pointer: %p\n", &Client)
        fmt.Println("Initializing GraphQL client with endpoint:", endpoint)
        Client = graphql.NewClient(endpoint)
        if Client == nil {
            panic("Failed to create graphql.Client")
        }
        fmt.Println("graphql.Client created successfully")
    })
}

