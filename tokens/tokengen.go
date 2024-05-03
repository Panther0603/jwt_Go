package tokens

import (
	"fmt"
	customerrors "jwt_use/customErrors"
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

type SignedDetails struct {
	Email    string
	Username string
	Uid      string
	jwt.StandardClaims
}

func GenerateToken(username string, email string, userid string) (token string, err error) {

	//
	claims := &SignedDetails{
		Email:    email,
		Username: username,
		Uid:      userid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
		},
	}
	godotenv.Load()
	fmt.Print("validate token gen >>" + os.Getenv("SECRET_KEY_JWT"))
	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("SECRET_KEY_JWT")))

	if err != nil {

		log.Println("Error is coming from here  " + err.Error())
		return "", err
	}
	return token, nil
}

func ValidateToken(stoken string) (claims *SignedDetails, msg string) {

	// token, err := jwt.ParseWithClaims(token, &SignedDetails{}, func(t *jwt.Token) (interface{}, error) {
	// 	return []byte(""), nil
	// })

	godotenv.Load()

	fmt.Print("validate token sec >>" + os.Getenv("SECRET_KEY_JWT"))
	token, err := jwt.ParseWithClaims(stoken, &SignedDetails{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY_JWT")), nil
	})

	if err != nil {
		msg = err.Error()
	}
	claims, ok := token.Claims.(*SignedDetails)

	if !ok {
		msg = customerrors.ErrInValidToken
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = customerrors.ErrTokenExpire
		return
	}
	return claims, msg
}
