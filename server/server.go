package server

import (
	"fmt"
	"net/http"
)

type App struct {
	mux *http.ServeMux
}

func NewApp() *App {
	return &App{
		mux: http.NewServeMux(),
	}
}

func (app *App) RunServer(port string) error {

	server := http.Server{
		Addr:    port,
		Handler: app.mux,
	}

	fmt.Printf("Servidor escuchando en el puerto %s\n", port)

	return server.ListenAndServe()
}
