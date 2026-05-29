package database

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
)

func TestConnect(t *testing.T) {

	origDBURL := os.Getenv("DB_URL")
	origConnectFunc := pgxConnectFunc

	defer func() {
		_ = os.Setenv("DB_URL", origDBURL)
		pgxConnectFunc = origConnectFunc // Siempre restauramos la función original
	}()

	t.Run("Error al cargar .env", func(t *testing.T) {
		_ = os.Remove(".env")

		conn, err := Connect()
		if err == nil {
			t.Errorf("se esperaba un error debido a la falta del archivo .env")
		}
		if conn != nil {
			t.Errorf("se esperaba que la conexión fuera nil")
		}

		file, _ := os.Create(".env")
		if file != nil {
			_ = file.Close()
		}
	})

	t.Run("Error de conexión - URL inválida", func(t *testing.T) {
		_ = os.WriteFile(".env", []byte(""), 0644)
		_ = os.Setenv("DB_URL", "postgres://usuario_invalido:clave_mala@localhost:9999/db_inexistente")

		conn, err := Connect()
		if err == nil {
			t.Errorf("se esperaba un error de conexión")
		}
		if conn != nil {
			t.Errorf("la conexión debió ser nil ante un fallo")
		}
	})

	t.Run("Conexión Exitosa - Cubre Println y Return", func(t *testing.T) {
		_ = os.WriteFile(".env", []byte(""), 0644)
		_ = os.Setenv("DB_URL", "postgres://localhost:5432/fake_db")

		pgxConnectFunc = func(ctx context.Context, connString string) (*pgx.Conn, error) {
			return &pgx.Conn{}, nil
		}

		conn, err := Connect()

		if err != nil {
			t.Errorf("no se esperaba ningún error, se obtuvo: %v", err)
		}
		if conn == nil {
			t.Errorf("se esperaba una instancia de conexión, se obtuvo nil")
		}
	})
}

func TestMain(m *testing.M) {
	file, err := os.Create(".env")
	if err == nil {
		_ = file.Close()
	}

	code := m.Run()
	_ = os.Remove(".env")
	os.Exit(code)
}
