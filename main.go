package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/yonchando/chirpy/internal/database"
	"github.com/yonchando/chirpy/internal/handlers"
	"github.com/yonchando/chirpy/internal/middleware"
	"github.com/yonchando/chirpy/internal/models"
)

func main() {

	dbUrl := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbUrl)

	if err != nil {
		log.Fatalln(err)
	}

	dbQueries := database.New(db)

	apiCfg := models.Config{}

	apiCfg.DB = *dbQueries

	mux := http.NewServeMux()

	mux.Handle("/app/", middleware.LogRequest(apiCfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("./internal/template"))))))

	mux.Handle("GET /api/healthz", middleware.LogRequest(handlers.Healthz()))

	mux.Handle("GET /admin/metrics", middleware.LogRequest(handlers.Metrics(&apiCfg)))

	mux.Handle("POST /admin/reset", middleware.LogRequest(handlers.Reset(&apiCfg)))

	mux.Handle("POST /api/validate_chirp", middleware.LogRequest(handlers.ValidateChirp()))

	mux.Handle("POST /api/users", middleware.LogRequest(handlers.PostUserHanlder(&apiCfg)))

	port := "8081"
	serve := &http.Server{
		Addr:    fmt.Sprintf("localhost:%s", port),
		Handler: mux,
	}

	fmt.Printf("Server start at http://localhost:%s\n", port)
	err = serve.ListenAndServe()

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Server start at http://localhost:8080")
}
