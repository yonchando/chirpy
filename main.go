package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	mux := http.NewServeMux()

	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("./internal/template/"))))

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		w.WriteHeader(http.StatusOK)

		w.Write([]byte("OK"))
	})

	serve := &http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	fmt.Println("Server start at http://localhost:8080")
	err := serve.ListenAndServe()

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Server start at http://localhost:8080")
}
