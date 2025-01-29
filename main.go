package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/yonchando/chirpy/internal/configs"
	"github.com/yonchando/chirpy/internal/database"
	"github.com/yonchando/chirpy/internal/routes"
)

func main() {
	godotenv.Load()

	dbUrl := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbUrl)

	if err != nil {
		log.Fatalln(err)
	}

	dbQueries := database.New(db)

	apiCfg := configs.Config{}

	apiCfg.DB = *dbQueries

	mux := http.NewServeMux()

	route := routes.Route{
		Port: "8081",
		Mux:  mux,
	}

	route.Hanlders(&apiCfg)

}
