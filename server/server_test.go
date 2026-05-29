package server

import (
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestNewApp(t *testing.T) {
	app := NewApp()
	if app == nil || app.mux == nil {
		t.Fatal("NewApp() debería retornar una instancia válida de App con su mux inicializado")
	}
}

func TestApp_RunServer(t *testing.T) {

	t.Run("Inicio exitoso del servidor en segundo plano", func(t *testing.T) {
		app := NewApp()

		port := ":0"

		errChan := make(chan error, 1)

		go func() {
			err := app.RunServer(port)
			if err != nil && err != http.ErrServerClosed {
				errChan <- err
			}
			close(errChan)
		}()

		time.Sleep(100 * time.Millisecond)

		select {
		case err := <-errChan:
			if err != nil {
				t.Fatalf("El servidor falló inesperadamente al iniciar: %v", err)
			}
		default:
			t.Log("Servidor levantado y escuchando correctamente")
		}
	})

	t.Run("Retorno de error ante puerto inválido", func(t *testing.T) {
		app := NewApp()

		invalidPort := "puerto_invalido_999999"

		err := app.RunServer(invalidPort)

		if err == nil {
			t.Error("se esperaba que RunServer retornara un error ante un puerto inválido, se obtuvo nil")
		}

		if err != nil && !strings.Contains(err.Error(), "address") {
			t.Logf("Aviso: El servidor retornó un error (correcto), pero con un mensaje inesperado: %v", err)
		}
	})
}
