package auth

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pass string) (string, error) {
	if len([]byte(pass)) > 20 {
		return "", errors.New("password cant exceed 20 bytes")
	}
	if len([]byte(pass)) < 10 {
		return "", errors.New("password cant be less than 8 bytes")
	}
	byteHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(byteHash), nil
}

func CheckHashedPassword(password string, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	now := &jwt.NumericDate{Time: time.Now()}
	expires := &jwt.NumericDate{Time: (time.Now().Add(expiresIn))}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{Issuer: "twitter", IssuedAt: now, ExpiresAt: expires, Subject: userID.String()})
	signedToken, err := token.SignedString([]byte(tokenSecret))

	if err != nil {
		return "", err
	}

	return signedToken, nil
}
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid token")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid subject UUID: %v", err)
	}

	return userID, nil

}


func GetBearerToken(headers http.Header)(string , error){

    bearerToken :=headers.Get("Authorization")
    bearerToken  =strings.Replace(bearerToken,"Bearer ","" , 1)
    log.Printf("bear is:%v" , bearerToken)



    return bearerToken , nil
}
















