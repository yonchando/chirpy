package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/yonchando/chirpy/internal/helper"
)

func ValidateChirp() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Body string `json:"body"`
		}

		decoder := json.NewDecoder(r.Body)
		params := parameters{}

		err := decoder.Decode(&params)

		w.Header().Set("Content-Type", "application/json")

		if err != nil {
			log.Println(err)
			helper.RespsonseWithError(w, 500, "Something went wrong")
			return
		}

		if len(params.Body) > 140 {
			helper.RespsonseWithError(w, 400, "Chirp is too long")
			return
		}

		body := strings.Split(params.Body, " ")
		resBody := make([]string, len(body))
		for i, v := range body {
			value := strings.ToLower(v)
			if value == "kerfuffle" || value == "sharbert" || value == "fornax" {
				resBody[i] = "****"
			} else {
				resBody[i] = v
			}
		}

		helper.ResponseWithJson(w, 200, struct {
			CleanedBody string `json:"cleaned_body"`
		}{
			CleanedBody: strings.Join(resBody, " "),
		})

	})
}
