package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sort"
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

		qAuthorID := r.URL.Query().Get("author_id")

		qSort := r.URL.Query().Get("sort")

		if qSort == "" || qSort != "desc" {
			qSort = "asc"
		}

		var err error
		var chirps []database.Chirp

		if qAuthorID == "" {
			chirps, err = cfg.DB.GetAllChirps(context.Background())
		} else {
			var authorID uuid.UUID
			authorID, err = uuid.Parse(qAuthorID)
			chirps, err = cfg.DB.GetAllChirpByAuthor(context.Background(), authorID)
		}

		if err != nil {
			log.Println(err)
			helper.ResponseWithError(w, http.StatusInternalServerError, "Something went wrong!")
			return
		}

		sort.Slice(chirps, func(i, j int) bool {
			if qSort == "asc" {
				return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
			}

			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})

		helper.ResponseWithJson(w, http.StatusOK, chirps)
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

		chirpID, err := uuid.Parse(r.PathValue("chirpID"))

		var chirp database.Chirp
		chirp, err = cfg.DB.FindChirpByID(context.Background(), chirpID)

		if err != nil {
			log.Println(err)
			helper.ResponseWithError(w, http.StatusNotFound, "Chirp not found.")
			return
		}

		if chirp.UserID != userID {
			helper.ResponseWithError(w, http.StatusForbidden, "Unauthorized")
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

func PostWebHook(cfg *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		apiKey, err := auth.GetAPIKey(r.Header)

		if err != nil || apiKey != cfg.PolkaKey {
			helper.ResponseWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		type parameters struct {
			Event string `json:"event"`
			Data  struct {
				UserID string `json:"user_id"`
			} `json:"data"`
		}

		decode := json.NewDecoder(r.Body)
		params := parameters{}

		err = decode.Decode(&params)
		if err != nil {
			log.Println(err)
			helper.ResponseWithError(w, http.StatusInternalServerError, "Something went wrong!.")
			return
		}

		if params.Event != "user.upgraded" {
			helper.ResponseWithJson(w, http.StatusNoContent, "")
			return
		}

		userId, err := uuid.Parse(params.Data.UserID)
		if err != nil {
			log.Println(err)
			helper.ResponseWithJson(w, http.StatusNoContent, "")
			return
		}

		_, err = cfg.DB.FindUserByID(context.Background(), userId)
		if err != nil {
			helper.ResponseWithError(w, http.StatusNotFound, "User not found.")
			return
		}

		cfg.DB.UpdateUsertoChirpRed(context.Background(), userId)

		helper.ResponseWithJson(w, http.StatusNoContent, "")
	})
}
