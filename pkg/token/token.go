package token

import (
	"context"
	"google.golang.org/grpc/metadata"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Token struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func GenerateJWT(email, secretKey string) (string, error) {

	var tokenClaim = Token{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaim)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return ``, err
	}

	return tokenString, nil
}

func VerifyJWT(bearerToken, secretKey string) (*jwt.Token, error) {

	token, err := jwt.ParseWithClaims(bearerToken, &Token{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	return token, err
}

func WriteToken(t string) context.Context {
	md := metadata.New(map[string]string{"token": t})
	return metadata.NewOutgoingContext(context.Background(), md)
}
