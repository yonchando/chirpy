package helper

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespsonseWithError(w http.ResponseWriter, code int, msg string) {

	w.WriteHeader(code)

	res, err := json.Marshal(struct {
		Error string `json:"error"`
	}{Error: msg})

	if err != nil {
		log.Println(err)
		return
	}

	w.Write(res)

}

func ResponseWithJson(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	dat, err := json.Marshal(payload)

	if err != nil {
		log.Println(err)
		return
	}

	w.Write(dat)
}
