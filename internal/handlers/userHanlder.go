package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/yonchando/chirpy/internal/database"
	"github.com/yonchando/chirpy/internal/helper"
	"github.com/yonchando/chirpy/internal/models"
)

func PostUserHanlder(cfg *models.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Email string `json:"email"`
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

		if params.Email == "" {
			helper.ResponseWithJson(w, http.StatusUnprocessableEntity, models.ErrorBody{
				Message: "Invalid data",
				Errors: parameters{
					Email: "Email is required",
				},
			})
			return
		}

		userParams := database.CreateUserParams{
			ID:        uuid.New(),
			Email:     params.Email,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		}

		var user database.User
		user, err = cfg.DB.CreateUser(context.Background(), userParams)

		if err != nil {
			log.Println(err)
			helper.ResponseWithError(w, http.StatusInternalServerError, "Created failed")
			return
		}

		u := models.User{}
		u.Map(user)

		helper.ResponseWithJson(w, http.StatusCreated, u)
	})
}
