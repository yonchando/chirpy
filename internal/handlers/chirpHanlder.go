package handlers

import (
	"context"
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

func GetChirpHanlder(cfg *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		chirps, err := cfg.DB.GetAllChirps(context.Background())

		if err != nil {
			log.Println(err)
			helper.ResponseWithError(w, http.StatusInternalServerError, "Created failed")
			return
		}

		c := make([]models.Chirp, len(chirps))

		for i, v := range chirps {

			item := models.Chirp{}
			item.Map(v)

			c[i] = item
		}

		helper.ResponseWithJson(w, http.StatusOK, c)
	})
}

func PostChirpHanlder(cfg *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		token, err := auth.GetBearerToken(r.Header)

		if err != nil {
			log.Println(err)
			helper.ResponseWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		var userID uuid.UUID

		userID, err = auth.ValidateJWT(token, cfg.TokenSecret)

		if err != nil {
			log.Println(err)
			helper.ResponseWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		type parameters struct {
			Body string `json:"body"`
		}

		decode := json.NewDecoder(r.Body)
		params := parameters{}

		err = decode.Decode(&params)
		if err != nil {
			log.Println(err)
			helper.ResponseWithError(w, http.StatusInternalServerError, "Something wrong with decode")
			return
		}

		errors := models.ErrorBody{}

		if params.Body == "" {
			errors.Message = "Invalid data"
			errors.Errors = map[string]string{
				"body": "Body is required",
			}
			helper.ResponseWithJson(w, 422, errors)
			return
		}

		if len(params.Body) > 140 {
			errors.Message = "Invalid data"
			errors.Errors = map[string]string{
				"body": "Body is too long!, Body must be less 140",
			}
			helper.ResponseWithJson(w, 422, errors)
			return
		}

		createParams := database.CreateChirpParams{
			ID:        uuid.New(),
			Body:      params.Body,
			UserID:    userID,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}

		var chirp database.Chirp
		chirp, err = cfg.DB.CreateChirp(context.Background(), createParams)

		if err != nil {
			log.Println(err)
			helper.ResponseWithError(w, http.StatusInternalServerError, "Created failed")
			return
		}

		c := models.Chirp{}
		c.Map(chirp)

		helper.ResponseWithJson(w, http.StatusCreated, c)

	})
}

func DeleteChirpHanlder(cfg *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		token, err := auth.GetBearerToken(r.Header)

		if err != nil {
			log.Println(err)
			helper.ResponseWithError(w, http.StatusUnauthorized, "Unauthenticate")
			return
		}

		var userID uuid.UUID

		userID, err = auth.ValidateJWT(token, cfg.TokenSecret)

		if err != nil {
			log.Println(err)
			helper.ResponseWithError(w, http.StatusUnauthorized, "Unauthenticate")
			return
		}

		chirpID, err := uuid.Parse(r.PathValue("chirpID"))

		var chirp database.Chirp
		chirp, err = cfg.DB.FindChirpByID(context.Background(), chirpID)

		if err != nil {
			log.Println(err)
			helper.ResponseWithError(w, http.StatusNotFound, "Chirp not found.")
			return
		}

		if chirp.UserID != userID {
			helper.ResponseWithError(w, http.StatusForbidden, "Unauthorizate")
			return
		}

		err = cfg.DB.DeleteChirpByID(context.Background(), chirpID)

		if err != nil {
			log.Println(err)
			helper.ResponseWithError(w, http.StatusInternalServerError, "Something went wrong!.")
			return
		}

		helper.ResponseWithJson(w, http.StatusNoContent, "")

	})
}

func ShowChirpHanlder(cfg *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		chirpID, err := uuid.Parse(r.PathValue("chirpID"))

		if err != nil {
			log.Println(err)
			helper.ResponseWithError(w, 422, "Invalid id")
			return
		}

		var chirp database.Chirp
		chirp, err = cfg.DB.FindChirpByID(context.Background(), chirpID)

		if err != nil {
			log.Println(err)
			helper.ResponseWithError(w, 404, "Not found!")
			return
		}

		c := models.Chirp{}

		c.Map(chirp)

		helper.ResponseWithJson(w, 200, c)
	})
}
