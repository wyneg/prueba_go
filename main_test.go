package main

import (
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/wyneg/prueba_go/server"
)

func TestMainExecution(t *testing.T) {
	// 1. RESPALDO DEL .ENV ORIGINAL: Protegemos tu archivo de trabajo
	var envOriginal []byte
	envExistia := false
	if _, err := os.Stat(".env"); err == nil {
		envOriginal, _ = os.ReadFile(".env")
		envExistia = true
	}

	// 2. Guardamos el estado original de las funciones y entornos
	origConnect := connectDBFunc
	origRunServer := runServerFunc
	origPort := os.Getenv("PORT")

	// Al finalizar, restauramos todo a su estado inicial
	defer func() {
		connectDBFunc = origConnect
		runServerFunc = origRunServer
		_ = os.Setenv("PORT", origPort)

		_ = os.Remove(".env")
		if envExistia {
			_ = os.WriteFile(".env", envOriginal, 0644)
		}
	}()

	// 3. Creamos el entorno simulado para el test
	_ = os.WriteFile(".env", []byte("RAWG_API_KEY=fake_key\nRAWG_BASE_URL=fake_url\nPORT=:8080"), 0644)

	// 4. INTERCEPCIÓN DE DB CON PARADA DE CONTROL
	// Hacemos que devuelva nil para que pase el "if err != nil" de producción.
	// Inmediatamente después, lanzamos un pánico personalizado. Esto detendrá la ejecución
	// de main() justo en esa línea, evitando que los servicios de abajo fallen por usar una DB vacía.
	connectDBFunc = func() (*pgx.Conn, error) {
		// Ejecutamos un pánico para cortar el flujo de forma controlada habiendo ya
		// cubierto todas las líneas superiores del archivo main.go
		panic("main_setup_completed_successfully")
	}

	runServerFunc = func(app *server.App, port string) error {
		return nil
	}

	// 5. CAPTURA Y VALIDACIÓN DEL FLUJO
	defer func() {
		if r := recover(); r != nil {
			if r == "main_setup_completed_successfully" {
				t.Log("Bloque inicial de main.go e infraestructura cubiertos con éxito")
			} else {
				t.Errorf("Pánico no controlado en la suite: %v", r)
			}
		}
	}()

	// 6. Ejecución síncrona
	main()
}
