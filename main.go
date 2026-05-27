package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/wyneg/prueba_go/database"
	"github.com/wyneg/prueba_go/server"
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

	app := server.NewApp()

	if err := app.RunServer(os.Getenv("PORT")); err != nil {
		log.Fatal(err)
	}

}
