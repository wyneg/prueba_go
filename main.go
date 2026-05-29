package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/wyneg/prueba_go/database"
	"github.com/wyneg/prueba_go/handlers"
	"github.com/wyneg/prueba_go/repositories"
	"github.com/wyneg/prueba_go/server"
	"github.com/wyneg/prueba_go/services"
)

var (
	connectDBFunc = database.Connect
	runServerFunc = func(app *server.App, port string) error {
		return app.RunServer(port)
	}
)

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error cargando archivo .env")
	}

	// db, err := database.Connect()
	db, err := connectDBFunc()

	if err != nil {
		log.Fatal(err)
	}

	// defer db.Close(context.Background())

	defer func() {
		if db != nil {
			db.Close(context.Background())
		}
	}()

	rawgService := services.NewRAWGService(os.Getenv("RAWG_API_KEY"), os.Getenv("RAWG_BASE_URL"))
	restHandler := handlers.NewRestHandler(rawgService)

	dbService := services.NewDBService(repositories.NewRestRepository(db))
	repositoryHandler := handlers.NewRepositoryHandler(dbService)

	app := server.NewApp()

	app.HttpMethods("GET", "/api/search", restHandler.GetGameHandler)
	app.HttpMethods("GET", "/api/games/{id}", restHandler.GetGameByIDHandler)
	app.HttpMethods("GET", "/api/library", repositoryHandler.GetGameHandler)
	app.HttpMethods("POST", "/api/library", repositoryHandler.CreateGameHandler)
	app.HttpMethods("PUT", "/api/library/{id}", repositoryHandler.UpdateGameHandler)
	app.HttpMethods("DELETE", "/api/library/{id}", repositoryHandler.DeleteGameHandler)
	app.HttpMethods("GET", "/api/library/stats", repositoryHandler.StatsGameHandler)

	// if err := app.RunServer(os.Getenv("PORT")); err != nil {
	// 	log.Fatal(err)
	// }

	if err := runServerFunc(app, os.Getenv("PORT")); err != nil {
		log.Fatal(err)
	}

}
