package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/yonchando/chirpy/internal/configs"
	"github.com/yonchando/chirpy/internal/handlers"
	"github.com/yonchando/chirpy/internal/middleware"
)

type Route struct {
	Port string
	Mux  *http.ServeMux
}

func (r *Route) Hanlders(config *configs.Config) {
	mux := r.Mux

	mux.Handle("/app/", middleware.LogRequest(config.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("./internal/template"))))))

	mux.Handle("GET /admin/metrics", middleware.LogRequest(handlers.Metrics(config)))
	mux.Handle("POST /admin/reset", middleware.LogRequest(handlers.Reset(config)))

	mux.Handle("POST /api/users", middleware.LogRequest(handlers.PostUserHandler(config)))
	mux.Handle("PUT /api/users", middleware.LogRequest(handlers.PutUserHandler(config)))

	mux.Handle("POST /api/login", middleware.LogRequest(handlers.LoginUserHandler(config)))
	mux.Handle("POST /api/refresh", middleware.LogRequest(handlers.RefreshTokenHandler(config)))
	mux.Handle("POST /api/revoke", middleware.LogRequest(handlers.RevokeRefreshTokenHandler(config)))

	mux.Handle("GET /api/chirps", middleware.LogRequest(handlers.GetChirpHanlder(config)))
	mux.Handle("GET /api/chirps/{chirpID}", middleware.LogRequest(handlers.ShowChirpHanlder(config)))
	mux.Handle("POST /api/chirps", middleware.LogRequest(handlers.PostChirpHanlder(config)))
	mux.Handle("DELETE /api/chirps/{chirpID}", middleware.LogRequest(handlers.DeleteChirpHanlder(config)))

	mux.Handle("GET /api/healthz", middleware.LogRequest(handlers.Healthz()))

	serve := &http.Server{
		Addr:    fmt.Sprintf("localhost:%s", r.Port),
		Handler: mux,
	}

	fmt.Printf("Server start at http://localhost:%s\n", r.Port)
	err := serve.ListenAndServe()

	if err != nil {
		log.Fatalln(err)
	}
}
