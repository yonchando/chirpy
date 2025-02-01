package configs

import (
	"log"
	"net/http"
	"sync/atomic"

	"github.com/yonchando/chirpy/internal/database"
)

type Config struct {
	FileserverHits atomic.Int32
	DB             database.Queries
	TokenSecret    string
	PolkaKey       string
}

func (cfg *Config) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cfg.FileserverHits.Add(1)

		log.Printf("Change Hits: %v", cfg.FileserverHits.Load())

		next.ServeHTTP(w, r)
	})
}
