package helper

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/yonchando/chirpy/internal/models"
)

func ResponseWithError(w http.ResponseWriter, code int, msg string) {

	w.WriteHeader(code)

	res, err := json.Marshal(models.ErrorBody{
		Message: msg,
	})

	if err != nil {
		log.Println(err)
		return
	}

	w.Write(res)

}

func ResponseWithJson(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)

	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(code)
	w.Write(dat)
}
