package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 16)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func GetBearerToken(headers http.Header) (string, error) {
	authorization := headers.Get("Authorization")

	if authorization == "" {
		return "", errors.New("Authorization not exists")
	}

	token := strings.Split(authorization, "Bearer ")

	return strings.TrimSpace(token[1]), nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	claim := &jwt.RegisteredClaims{
		Issuer:    "chipry",
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	ss, err := token.SignedString([]byte(tokenSecret))

	if err != nil {
		return "", err
	}

	return ss, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claim := jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(tokenString, &claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	}, jwt.WithLeeway(5*time.Second))

	if err != nil {
		log.Println(err)
		return uuid.Nil, err
	}

	var ID uuid.UUID
	ID, err = uuid.Parse(claim.Subject)

	if err != nil {
		log.Println(err)
		return uuid.Nil, err
	}

	return ID, nil

}

func MakeRefreshToken() (string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		log.Println("Error: ", err)
		return "", err
	}

	refresh_token := hex.EncodeToString(b)

	return refresh_token, nil
}
