package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {

	password := "password"

	hash, _ := HashPassword(password)

	if err := CheckPasswordHash(password, hash); err != nil {
		t.Errorf("Expects <nil>, got %v", err)
	}

	if err := CheckPasswordHash("wrong-password", hash); err == nil {
		t.Errorf("Expects not <nil>, got %v", err)
	}
}

func TestGetBearerToken(t *testing.T) {
	headers := http.Header{}

	headers.Set("Authorization", "Bearer token")

	token, err := GetBearerToken(headers)

	if err != nil {
		t.Errorf("Expected token, got %v", err)
	}

	if token != "token" {
		t.Errorf("Expected token, got %v", token)
	}
}

func TestMakeJWT(t *testing.T) {
	tests := struct {
		ID          uuid.UUID
		tokenSecret string
		expiresIn   time.Duration
	}{
		ID:          uuid.New(),
		tokenSecret: "random string token",
		expiresIn:   time.Hour * 24,
	}

	ss, _ := MakeJWT(tests.ID, tests.tokenSecret, tests.expiresIn)

	if ss == "" {
		t.Errorf("Expect return token")
	}
}

func TestValidateJWT(t *testing.T) {
	cases := struct {
		ID          uuid.UUID
		tokenSecret string
		expiresIn   time.Duration
	}{
		ID:          uuid.New(),
		tokenSecret: "random string token",
		expiresIn:   time.Hour * 24,
	}

	ss, err := MakeJWT(cases.ID, cases.tokenSecret, cases.expiresIn)

	if err != nil {
		t.Errorf("Expect not err but got %s", err)
	}

	var id uuid.UUID
	id, err = ValidateJWT(ss, cases.tokenSecret)

	if id != cases.ID {
		t.Errorf("Expect %v to match %v", id, cases.ID)
	}
}

func TestMakeRefreshToken(t *testing.T) {
	refreh_token, err := MakeRefreshToken()

	if refreh_token == "" || err != nil {
		t.Error("Failed to generate refresh token")
	}
}
