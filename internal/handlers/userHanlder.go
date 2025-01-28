package handlers

import (
	"net/http"

	"github.com/yonchando/chirpy/internal/models"
)

func PostUserHanlder(cfg *models.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
