package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/yonchando/chirpy/internal/auth"
	"github.com/yonchando/chirpy/internal/configs"
	"github.com/yonchando/chirpy/internal/database"
	"github.com/yonchando/chirpy/internal/helper"
	"github.com/yonchando/chirpy/internal/models"
)

func PostUserHandler(cfg *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		w.Header().Set("Content-Type", "application/json")

		decode := json.NewDecoder(r.Body)
		params := parameters{}

		err := decode.Decode(&params)
		if err != nil {
			log.Println(err)
			helper.ResponseWithError(w, http.StatusInternalServerError, "Something wrong with decode")
			return
		}

		errors := map[string]string{}

		if params.Email == "" {
			errors["email"] = "Email is required"
		}

		if params.Password == "" {
			errors["password"] = "Password is required"
		}

		if len(errors) > 0 {
			helper.ResponseWithJson(w, http.StatusUnprocessableEntity, models.ErrorBody{
				Message: "Invalid data",
				Errors:  errors,
			})
			return
		}

		var password string

		password, err = auth.HashPassword(params.Password)
		if err != nil {
			log.Println(err)

			helper.ResponseWithError(w, http.StatusInternalServerError, "Something went wrong")
			return
		}

		userParams := database.CreateUserParams{
			ID:             uuid.New(),
			Email:          params.Email,
			HashedPassword: password,
			CreatedAt:      time.Now().UTC(),
			UpdatedAt:      time.Now().UTC(),
		}

		var user database.User
		user, err = cfg.DB.CreateUser(context.Background(), userParams)

		if err != nil {
			log.Println(err)
			helper.ResponseWithError(w, http.StatusInternalServerError, "Something went wrong.")
			return
		}

		u := models.User{}
		u.Setter(user)

		helper.ResponseWithJson(w, http.StatusCreated, u)
	})
}

func PutUserHandler(cfg *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			helper.ResponseWithError(w, 401, "Unauthorized")
			return
		}

		var userId uuid.UUID

		userId, err = auth.ValidateJWT(token, cfg.TokenSecret)
		if err != nil {
			helper.ResponseWithError(w, 401, "Unauthorized")
			return
		}

		type parameters struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		w.Header().Set("Content-Type", "application/json")

		decode := json.NewDecoder(r.Body)
		params := parameters{}

		err = decode.Decode(&params)
		if err != nil {
			log.Println(err)
			helper.ResponseWithError(w, http.StatusInternalServerError, "Something wrong with decode")
			return
		}

		errors := map[string]string{}

		if params.Email == "" {
			errors["email"] = "Email is required"
		}

		if params.Password == "" {
			errors["password"] = "Password is required"
		}

		if len(errors) > 0 {
			helper.ResponseWithJson(w, http.StatusUnprocessableEntity, models.ErrorBody{
				Message: "Invalid data",
				Errors:  errors,
			})
			return
		}

		var password string

		password, err = auth.HashPassword(params.Password)
		if err != nil {
			log.Println(err)

			helper.ResponseWithError(w, http.StatusInternalServerError, "Something went wrong")
			return
		}

		userParams := database.UpdateUserByIDParams{
			ID:             userId,
			Email:          params.Email,
			HashedPassword: password,
			UpdatedAt:      time.Now().UTC(),
		}

		err = cfg.DB.UpdateUserByID(context.Background(), userParams)

		if err != nil {
			log.Println(err)
			helper.ResponseWithError(w, http.StatusInternalServerError, "Something went wrong.")
			return
		}

		user := models.User{}

		err = user.FindByID(cfg, userId)

		if err != nil {
			log.Println(err)
			helper.ResponseWithError(w, http.StatusNotFound, "User not found.")
			return
		}

		helper.ResponseWithJson(w, http.StatusOK, user)
	})
}

func LoginUserHandler(cfg *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		w.Header().Set("Content-Type", "application/json")

		decode := json.NewDecoder(r.Body)
		params := parameters{}

		err := decode.Decode(&params)
		if err != nil {
			log.Println(err)
			helper.ResponseWithError(w, http.StatusInternalServerError, "Something wrong with decode")
			return
		}

		errors := map[string]string{}

		if params.Email == "" {
			errors["email"] = "Email is required"
		}

		if params.Password == "" {
			errors["password"] = "Password is required"
		}

		if len(errors) > 0 {
			helper.ResponseWithJson(w, http.StatusUnprocessableEntity, models.ErrorBody{
				Message: "Invalid data",
				Errors:  errors,
			})
			return
		}

		u := models.User{}
		_, err = u.Authenticate(cfg, params.Email, params.Password)

		if err != nil {
			helper.ResponseWithError(w, http.StatusUnauthorized, "Email or password is invalid.")
			return
		}

		err = u.CreateToken(cfg)

		if err != nil {
			helper.ResponseWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		err = u.CreateRefreshToken(cfg)

		if err != nil {
			helper.ResponseWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		helper.ResponseWithJson(w, http.StatusOK, u)
	})
}

func RefreshTokenHandler(cfg *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		token, err := auth.GetBearerToken(r.Header)

		if err != nil {
			log.Println(err)
			helper.ResponseWithError(w, 401, "Unauthorized")
			return
		}

		refreshToken := database.RefreshToken{
			RevokedAt: sql.NullTime{
				Valid: false,
			},
		}
		contxt := context.Background()
		refreshToken, err = cfg.DB.FindRefreshToken(contxt, token)

		if err != nil {
			log.Println(err)
			helper.ResponseWithError(w, 401, "Unauthorized")
			return
		}

		if refreshToken.ExpiresAt.Before(time.Now().UTC()) {
			helper.ResponseWithJson(w, 401, models.ErrorBody{
				Message: "Refresh token is expired.",
			})
			return
		}

		user := models.User{}

		err = user.FindByID(cfg, refreshToken.UserID)
		if err != nil {
			log.Println(err)
			helper.ResponseWithError(w, 401, "Unauthorized")
			return
		}

		err = user.CreateToken(cfg)

		if err != nil {
			log.Println(err)
			helper.ResponseWithError(w, 401, "Unauthorized")
			return
		}

		helper.ResponseWithJson(w, 200, map[string]string{
			"token": user.Token,
		})

	})
}

func RevokeRefreshTokenHandler(cfg *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		token, err := auth.GetBearerToken(r.Header)

		if err != nil {
			log.Println(err)
			helper.ResponseWithError(w, 401, "Unauthorized")
			return
		}

		contxt := context.Background()

		err = cfg.DB.RevokeRefreshToken(contxt, database.RevokeRefreshTokenParams{
			RevokedAt: sql.NullTime{
				Time:  time.Now().UTC(),
				Valid: true,
			},
			Token: token,
		})

		if err != nil {
			log.Println(err)
			helper.ResponseWithError(w, 401, "Unauthorized")
			return
		}

		helper.ResponseWithJson(w, 204, "")

	})
}
