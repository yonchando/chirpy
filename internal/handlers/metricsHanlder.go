package handlers

import (
	"fmt"
	"log"
	"net/http"

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

		w.WriteHeader(http.StatusOK)

		apiCfg.FileserverHits.Store(0)

		w.Write([]byte("Reset hit...\n"))
	})

}
