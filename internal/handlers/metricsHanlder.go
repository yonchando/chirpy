package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/yonchando/chirpy/internal/helper"
	"github.com/yonchando/chirpy/internal/models"
)

func Metrics(apiCfg *models.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		w.WriteHeader(http.StatusOK)

		hit := apiCfg.FileserverHits.Load()
		log.Printf("Show Hits: %v", hit)

		body := `
        <html lang="en">
            <head>
                <meta charset="UTF-8" />
                <meta name="viewport" content="width=device-width, initial-scale=1.0" />
                <title>Metrics</title>
            </head>
            <body>
                <h1>Welcome, Chirpy Admin</h1>
                <p>Chirpy has been visited %d times!</p>
            </body>
        </html>
        `

		w.Write([]byte(fmt.Sprintf(body, hit)))
	})
}

func Reset(apiCfg *models.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		platform := os.Getenv("PLATFORM")

		if strings.ToLower(platform) != "dev" {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("You don't have permissions."))
			return
		}

		w.WriteHeader(http.StatusOK)

		apiCfg.FileserverHits.Store(0)

		err := apiCfg.DB.DeleteAllUser(context.Background())

		if err != nil {
			helper.ResponseWithError(w, http.StatusInternalServerError, "Something went wrong!")
			return
		}

		w.Write([]byte("Reset"))
	})

}
