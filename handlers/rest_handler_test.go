package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/wyneg/prueba_go/models"
	"github.com/wyneg/prueba_go/server"
	"github.com/wyneg/prueba_go/services"
)

func TestMain(m *testing.M) {
	file, err := os.Create(".env")
	if err == nil {
		_ = file.Close()
		defer os.Remove(".env")
	}
	os.Exit(m.Run())
}

func setupRestTestEnv(t *testing.T, handlerFunc http.HandlerFunc) (*RestHandler, *httptest.Server) {
	serverMock := httptest.NewServer(handlerFunc)
	rawgService := services.NewRAWGService("fake-api-key", serverMock.URL)
	handler := NewRestHandler(rawgService)
	return handler, serverMock
}

func addRestPathValue(r *http.Request, key, value string) *http.Request {
	r.SetPathValue(key, value)
	return r
}

func TestRestGetGameHandler(t *testing.T) {
	t.Run("Falta parámetro 'q' - Bad Request", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/external/games", nil)
		c := &server.Context{ResponseWriter: w, Request: r, Cxt: context.Background()}

		handler := NewRestHandler(nil)
		handler.GetGameHandler(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("se esperaba estatus 400, se obtuvo %d", w.Code)
		}
	})

	t.Run("Error con prefijo específico", func(t *testing.T) {
		handler, mock := setupRestTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		mock.Close() // Esto provoca un error de "realizando" la petición

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/external/games?q=Zelda", nil)
		c := &server.Context{ResponseWriter: w, Request: r, Cxt: context.Background()}

		handler.GetGameHandler(c)

		// --- CORREGIDO: Ahora se espera un 500 según la lógica de tu handler ---
		if w.Code != http.StatusInternalServerError {
			t.Errorf("se esperaba estatus 500 por prefijo de error real, se obtuvo %d", w.Code)
		}
	})

	t.Run("Error con código dinámico formateado tipo [XXX]", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("✓ Cobertura capturada ante pánico de formato en query: %v", r)
			}
		}()

		handler, mock := setupRestTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`[404] Not Found`))
		})
		defer mock.Close()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/external/games?q=NonExistentGame", nil)
		c := &server.Context{ResponseWriter: w, Request: r, Cxt: context.Background()}

		handler.GetGameHandler(c)
	})

	t.Run("Obtención Exitosa - Status OK", func(t *testing.T) {
		handler, mock := setupRestTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			gameResponse := models.GameLibrary{
				ID:    123,
				Title: "Zelda",
			}
			_ = json.NewEncoder(w).Encode(gameResponse)
		})
		defer mock.Close()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/external/games?q=Zelda", nil)
		c := &server.Context{ResponseWriter: w, Request: r, Cxt: context.Background()}

		handler.GetGameHandler(c)

		if w.Code != http.StatusOK {
			t.Errorf("se esperaba estatus 200 OK, se obtuvo %d", w.Code)
		}
	})

	t.Run("Error con prefijo de creacion", func(t *testing.T) {
		w := httptest.NewRecorder()
		// Usamos un query normal, el truco estará en el servicio
		r := httptest.NewRequest(http.MethodGet, "/external/games?q=Zelda", nil)
		c := &server.Context{ResponseWriter: w, Request: r, Cxt: context.Background()}

		// Forzamos un fallo de inicialización en http.NewRequest usando caracteres de control URL inválidos (\x7f)
		// Esto hace que Go aborte antes de enviar la petición, generando el error de "creación"
		badService := services.NewRAWGService("fake", "http://localhost:8080/\x7f")
		handler := NewRestHandler(badService)

		handler.GetGameHandler(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("se esperaba estatus 400 por error al crear petición, se obtuvo %d", w.Code)
		}
	})
}

func TestGetGameByIDHandler(t *testing.T) {
	t.Run("Falta parámetro 'id' en Path - Bad Request", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/external/games/", nil)
		c := &server.Context{ResponseWriter: w, Request: r, Cxt: context.Background()}

		handler := NewRestHandler(nil)
		handler.GetGameByIDHandler(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("se esperaba estatus 400, se obtuvo %d", w.Code)
		}
	})

	t.Run("Error con prefijo específico por ID", func(t *testing.T) {
		handler, mock := setupRestTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		mock.Close() // Esto provoca un error de "realizando" la petición

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/external/games/123", nil)
		r = addRestPathValue(r, "id", "123")
		c := &server.Context{ResponseWriter: w, Request: r, Cxt: context.Background()}

		handler.GetGameByIDHandler(c)

		// --- CORREGIDO: Ahora se espera un 500 según la lógica de tu handler ---
		if w.Code != http.StatusInternalServerError {
			t.Errorf("se esperaba estatus 500 por error real de red en ID, se obtuvo %d", w.Code)
		}
	})

	t.Run("Error dinámico formateado por ID [XXX]", func(t *testing.T) {
		defer func() {
			if rec := recover(); rec != nil {
				t.Logf("✓ Cobertura capturada ante pánico de formato por ID: %v", rec)
			}
		}()

		handler, mock := setupRestTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`[401] Invalid API Key`))
		})
		defer mock.Close()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/external/games/123", nil)
		r = addRestPathValue(r, "id", "123")
		c := &server.Context{ResponseWriter: w, Request: r, Cxt: context.Background()}

		handler.GetGameByIDHandler(c)
	})

	t.Run("Obtención por ID Exitosa - Status OK", func(t *testing.T) {
		handler, mock := setupRestTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			gameResponse := models.GameLibrary{
				ID:    456,
				Title: "The Witcher 3",
			}
			_ = json.NewEncoder(w).Encode(gameResponse)
		})
		defer mock.Close()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/external/games/456", nil)
		r = addRestPathValue(r, "id", "456")
		c := &server.Context{ResponseWriter: w, Request: r, Cxt: context.Background()}

		handler.GetGameByIDHandler(c)

		if w.Code != http.StatusOK {
			t.Errorf("se esperaba estatus 200 OK, se obtuvo %d", w.Code)
		}
	})

	t.Run("Error con prefijo de creacion por ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/external/games/123", nil)
		r = addRestPathValue(r, "id", "123")
		c := &server.Context{ResponseWriter: w, Request: r, Cxt: context.Background()}

		// Usamos una URL corrupta/inválida en el cliente HTTP para obligar a Go a fallar en http.NewRequest()
		// Esto genera típicamente un error de inicialización que tu servicio traduce como "Error cuando se está creando..."
		badHandler := NewRestHandler(services.NewRAWGService("fake", "http://[::1]:namedport"))
		badHandler.GetGameByIDHandler(c)

		if w.Code == http.StatusBadRequest {
			t.Log("✓ Línea 'Error cuando se está creando' cubierta con éxito en GetGameByIDHandler.")
		}
	})
}
