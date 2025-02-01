package models

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/yonchando/chirpy/internal/auth"
	"github.com/yonchando/chirpy/internal/configs"
	"github.com/yonchando/chirpy/internal/database"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Token        string    `json:"token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
}

func (u *User) FindByID(cfg *configs.Config, ID uuid.UUID) error {
	user, err := cfg.DB.FindUserByID(context.Background(), ID)

	if err != nil {
		return err
	}

	u.Setter(user)

	return nil
}

func (u *User) Authenticate(cfg *configs.Config, email, password string) (database.User, error) {

	user, err := cfg.DB.FindUserByEmail(context.Background(), email)

	if err != nil {
		log.Println(err)
		return database.User{}, err
	}

	if auth.CheckPasswordHash(password, user.HashedPassword) != nil {
		return database.User{}, errors.New("invalid password")
	}

	u.Setter(user)

	return user, nil
}

func (u *User) Setter(user database.User) {
	u.ID = user.ID
	u.Email = user.Email
	u.CreatedAt = user.CreatedAt
	u.UpdatedAt = user.UpdatedAt
}

func (u *User) CreateToken(cfg *configs.Config) error {
	var token string
	expiresInHours := time.Duration(time.Hour)

	token, err := auth.MakeJWT(u.ID, cfg.TokenSecret, expiresInHours)

	if err != nil {
		return err
	}

	u.Token = token

	return nil
}

func (u *User) CreateRefreshToken(cfg *configs.Config) error {

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		return err
	}

	refreshTokenParams := database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    u.ID,
		ExpiresAt: time.Now().UTC().AddDate(0, 0, 60),
		RevokedAt: sql.NullTime{
			Valid: false,
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	_, err = cfg.DB.CreateRefreshToken(context.Background(), refreshTokenParams)

	if err != nil {
		return err
	}

	u.RefreshToken = refreshToken

	return nil
}
