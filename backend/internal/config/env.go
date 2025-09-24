package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var once sync.Once

func LoadEnv() {
	once.Do(func() {
		err := godotenv.Load(".env")
		if err != nil {
			log.Println("NO env file found or error loading it", err)
		}

	})
}

func LoadJWTSecret() []byte {
	LoadEnv()
	return []byte(os.Getenv("HASURA_GRAPHQL_JWT_SECRET"))
}

func LoadADMINSecret() []byte {
	LoadEnv()
	return []byte(os.Getenv("HASURA_GRAPHQL_ADMIN_SECRET"))
}

func ENDPoint() []byte {
	LoadEnv()
	return []byte(os.Getenv("END_POINT"))
}
func CLOUDINARYCREDINTIALS() (string,string,string){
	LoadEnv()
	return os.Getenv("CLOUDINARY_CLOUD_NAME"), os.Getenv("CLOUDINARY_API_KEY"), os.Getenv("CLOUDINARY_API_SECRET")

}

func CHAPAPAYMENTCREDINTIALS()(string,string,string){
	LoadEnv()
	return os.Getenv("CHAPA_SECRET_KEY"),os.Getenv("CHAPA_CALLBACK_URL"),os.Getenv("CHAPA_RETURN_URL")
}
