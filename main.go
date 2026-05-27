package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/wyneg/prueba_go/database"
	"github.com/wyneg/prueba_go/handlers"
	"github.com/wyneg/prueba_go/server"
	"github.com/wyneg/prueba_go/services"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error cargando archivo .env")
	}

	db, err := database.Connect()

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close(context.Background())

	rawgService := services.NewRAWGService(os.Getenv("RAWG_API_KEY"), os.Getenv("RAWG_BASE_URL"))

	restHandler := handlers.NewRestHandler(rawgService)

	app := server.NewApp()

	app.HttpMethods("GET", "/api/search", restHandler.GetGameHandler)
	app.HttpMethods("GET", "/api/games/{id}", restHandler.GetGameByIDHandler)

	if err := app.RunServer(os.Getenv("PORT")); err != nil {
		log.Fatal(err)
	}

}
