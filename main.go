package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"

	_ "github.com/lib/pq"

	"github.com/yonchando/chirpy/internal/helper"
	"github.com/yonchando/chirpy/internal/middleware"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cfg.fileserverHits.Add(1)

		log.Printf("Change Hits: %v", cfg.fileserverHits.Load())

		next.ServeHTTP(w, r)
	})
}

func healthz() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		w.WriteHeader(http.StatusOK)

		w.Write([]byte("OK\n"))
	})
}

func metrics(apiCfg *apiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		w.WriteHeader(http.StatusOK)

		hit := apiCfg.fileserverHits.Load()
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
func reset(apiCfg *apiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		w.WriteHeader(http.StatusOK)

		apiCfg.fileserverHits.Store(0)

		w.Write([]byte("Reset hit...\n"))
	})

}

func validateChirp() http.Handler {
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

func main() {

	mux := http.NewServeMux()

	apiCfg := apiConfig{}

	mux.Handle("/app/", middleware.LogRequest(apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("./internal/template"))))))

	mux.Handle("GET /api/healthz", middleware.LogRequest(healthz()))

	mux.Handle("GET /admin/metrics", middleware.LogRequest(metrics(&apiCfg)))

	mux.Handle("POST /admin/reset", middleware.LogRequest(reset(&apiCfg)))

	mux.Handle("POST /api/validate_chirp", middleware.LogRequest(validateChirp()))

	port := "8081"
	serve := &http.Server{
		Addr:    fmt.Sprintf("localhost:%s", port),
		Handler: mux,
	}

	fmt.Printf("Server start at http://localhost:%s\n", port)
	err := serve.ListenAndServe()

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Server start at http://localhost:8080")
}
