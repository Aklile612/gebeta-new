package auth

import (
	"errors"
	"time"

	"github.com/aklile/recipe-backend/internal/config"
	"github.com/golang-jwt/jwt/v5"
)


var jwtSecret= config.LoadJWTSecret()

type Claims struct{
	UserID string `json:"sub"`
	Email     string `json:"email"`      
	FullName  string `json:"full_name"`
	jwt.RegisteredClaims 
	HasuraClaims   HasuraClaims  `json:"https://hasura.io/jwt/claims"` 
}

type HasuraClaims struct{
	DefaultRole string `json:"x-hasura-default-role"`
	AllowedRoles []string `json:"x-hasura-allowed-roles"`
	UserID string `json:"x-hasura-user-id"`
}

func GenerateJWT (userID,email,fullName string) (string,error){
	claims:= Claims{
		UserID: userID,
		Email:    email,
		FullName: fullName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
		HasuraClaims: HasuraClaims{
			DefaultRole: "user",
			AllowedRoles: []string{"user"},
			UserID: userID,
		},
	}
	token:= jwt.NewWithClaims(jwt.SigningMethodHS256,claims)

	return token.SignedString(jwtSecret)
}


func ValidateJWT(tokenStr string) (*Claims,error){
	token,err:= jwt.ParseWithClaims(tokenStr,&Claims{},func(token *jwt.Token)(interface{},error){
		if _,ok := token.Method.(*jwt.SigningMethodHMAC); !ok{
			return nil,errors.New("unexpected signing method")
		}
		return  jwtSecret,nil
	})
	if err!= nil{
		return  nil,err
	}
	claims,ok:= token.Claims.(*Claims)

	if !ok || !token.Valid{
		return nil,errors.New("invalid token")
	}

	return claims,nil
}