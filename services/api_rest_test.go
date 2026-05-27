package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/wyneg/prueba_go/models"
)

func TestGetGame(t *testing.T) {

	nextURL := "http://api.rawg.io/next_page"

	tests := []struct {
		name             string
		gameName         string
		mockStatus       int
		mockResponseBody interface{}
		expectedError    string
		validateResult   bool
	}{
		{
			name:       "Caso Exitoso - Retorna juegos",
			gameName:   "witcher",
			mockStatus: http.StatusOK,
			mockResponseBody: models.RAWGResponse{
				Count: 1,
				Next:  &nextURL,
				Results: []models.Game{
					{ID: 123, Name: "The Witcher 3"},
				},
			},
			expectedError:  "",
			validateResult: true,
		},
		{
			name:             "Error 404 - Juego no encontrado",
			gameName:         "juego_inexistente",
			mockStatus:       http.StatusNotFound,
			mockResponseBody: "Not Found",
			expectedError:    "[404] Juego no encontrado",
			validateResult:   false,
		},
		{
			name:             "Error 400 - Solicitud Incorrecta",
			gameName:         "invalid_query",
			mockStatus:       http.StatusBadRequest,
			mockResponseBody: "Bad Request",
			expectedError:    "[400] Solicitud incorrecta",
			validateResult:   false,
		},
		{
			name:             "Error de Decodificación - JSON Malformado",
			gameName:         "broken_json",
			mockStatus:       http.StatusOK,
			mockResponseBody: "{malformed json",
			expectedError:    "Error decodificando la respuesta:",
			validateResult:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if !strings.Contains(r.URL.Path, "/games") {
					t.Errorf("Se esperaba llamada a /games, se llamó a: %s", r.URL.Path)
				}

				w.WriteHeader(tt.mockStatus)

				if strBody, ok := tt.mockResponseBody.(string); ok {
					w.Write([]byte(strBody))
					return
				}

				json.NewEncoder(w).Encode(tt.mockResponseBody)
			}))

			defer server.Close()

			svc := &RAWGService{
				ApiKey:  "test_key",
				BaseURL: server.URL,
				httpClient: &http.Client{
					Timeout: 2 * time.Second,
				},
			}

			res, err := svc.GetGame(tt.gameName)

			if tt.expectedError != "" {
				if err == nil {
					t.Fatalf("Se esperaba un error que contuviera '%s', pero no se obtuvo ninguno", tt.expectedError)
				}
				if !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("Se obtuvo el error: '%s'\nSe esperaba que contuviera: '%s'", err.Error(), tt.expectedError)
				}
			} else {
				if err != nil {
					t.Fatalf("No se esperaba ningún error, pero ocurrió: %v", err)
				}
				if tt.validateResult && res == nil {
					t.Error("Se esperaba una respuesta válida, se obtuvo nil")
				}
			}
		})
	}
}

func TestGetGame_ErrorRealizandoSolicitud(t *testing.T) {
	svc := &RAWGService{
		ApiKey:  "test_key",
		BaseURL: "http://localhost:9999",
		httpClient: &http.Client{
			Timeout: 1 * time.Millisecond,
		},
	}

	_, err := svc.GetGame("test")

	if err == nil {
		t.Fatal("Se esperaba un error de red, pero la solicitud se completó con éxito")
	}

	expectedPrefix := "Error cuando se está realizando la solicitud:"
	if !strings.Contains(err.Error(), expectedPrefix) {
		t.Errorf("Se esperaba un error que contuviera '%s', pero se obtuvo: %v", expectedPrefix, err)
	}
}

func TestGetGame_ErrorCreandoSolicitud(t *testing.T) {
	urlCorrupta := "http://api.rawg.io\x7f/v1"

	svc := &RAWGService{
		ApiKey:  "test_key",
		BaseURL: urlCorrupta,
		httpClient: &http.Client{
			Timeout: http.DefaultClient.Timeout,
		},
	}

	_, err := svc.GetGame("any-game")

	if err == nil {
		t.Fatal("Se esperaba un error al crear la solicitud, pero err fue nil")
	}

	expectedText := "Error cuando se está creando la solicitud:"
	if !strings.Contains(err.Error(), expectedText) {
		t.Errorf("Se esperaba un error que contuviera:\n'%s'\nPero se obtuvo:\n'%s'", expectedText, err.Error())
	}
}

func TestNewRAWGService_Success(t *testing.T) {
	envContent := "ANY_KEY=mock_value"
	err := os.WriteFile(".env.test", []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("No se pudo crear el archivo .env de prueba: %v", err)
	}

	t.Cleanup(func() {
		os.Remove(".env.test")
	})

	apiKey := "mi-api-key"
	baseURL := "https://api.rawg.io/api"
	svc := NewRAWGService(apiKey, baseURL, ".env.test")

	if svc == nil {
		t.Fatal("Se esperaba un struct RAWGService, se obtuvo nil")
	}
	if svc.ApiKey != apiKey {
		t.Errorf("Se esperaba ApiKey '%s', se obtuvo '%s'", apiKey, svc.ApiKey)
	}
	if svc.BaseURL != baseURL {
		t.Errorf("Se esperaba BaseURL '%s', se obtuvo '%s'", baseURL, svc.BaseURL)
	}
	if svc.httpClient == nil {
		t.Error("El httpClient no fue inicializado")
	} else if svc.httpClient.Timeout != 10*time.Second {
		t.Errorf("Se esperaba un timeout de 10s, se obtuvo %v", svc.httpClient.Timeout)
	}
}

func TestNewRAWGService_MissingEnv_ShouldFatal(t *testing.T) {
	cmd := exec.Command(os.Args[0], "-test.run=TestRenderFatalHelper")
	cmd.Env = append(os.Environ(), "BE_CRASHING_TEST=1")

	err := cmd.Run()

	if err == nil {
		t.Fatal("Se esperaba que la función causara un log.Fatal (exit status 1), pero el proceso terminó con éxito")
	}

	if _, ok := err.(*exec.ExitError); !ok {
		t.Errorf("Se esperaba un ExitError debido al log.Fatal, se obtuvo: %v", err)
	}
}

func TestRenderFatalHelper(t *testing.T) {
	if os.Getenv("BE_CRASHING_TEST") != "1" {
		return
	}

	NewRAWGService("key", "url", ".env.fallar")
}

func TestGetGameByID(t *testing.T) {
	tests := []struct {
		name             string
		gameID           string
		mockStatus       int
		mockResponseBody interface{}
		expectedError    string
		validateResult   bool
	}{
		{
			name:       "Caso Exitoso - Retorna detalles del juego",
			gameID:     "479836",
			mockStatus: http.StatusOK,
			mockResponseBody: models.GameDetail{
				ID:           479836,
				Slug:         "zelda-2",
				Name:         "ZELDA (Raul Fernandes)",
				NameOriginal: "ZELDA (Raul Fernandes)",
				Released:     "2020-08-13",
			},
			expectedError:  "",
			validateResult: true,
		},
		{
			name:             "Error 404 - Juego por ID no encontrado",
			gameID:           "9999999",
			mockStatus:       http.StatusNotFound,
			mockResponseBody: "Not Found",
			expectedError:    "[404] Juego no encontrado",
			validateResult:   false,
		},
		{
			name:             "Error 400 - Solicitud Incorrecta por ID",
			gameID:           "invalid_id",
			mockStatus:       http.StatusBadRequest,
			mockResponseBody: "Bad Request",
			expectedError:    "[400] Solicitud incorrecta",
			validateResult:   false,
		},
		{
			name:             "Error de Decodificación - JSON Malformado",
			gameID:           "479836",
			mockStatus:       http.StatusOK,
			mockResponseBody: "{malformed json",
			expectedError:    "Error decodificando la respuesta:",
			validateResult:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				expectedPath := "/games/" + tt.gameID
				if !strings.Contains(r.URL.Path, expectedPath) {
					t.Errorf("Se esperaba llamada a %s, se llamó a: %s", expectedPath, r.URL.Path)
				}

				w.WriteHeader(tt.mockStatus)

				if strBody, ok := tt.mockResponseBody.(string); ok {
					w.Write([]byte(strBody))
					return
				}

				json.NewEncoder(w).Encode(tt.mockResponseBody)
			}))
			defer server.Close()

			svc := &RAWGService{
				ApiKey:  "test_key",
				BaseURL: server.URL,
				httpClient: &http.Client{
					Timeout: 2 * time.Second,
				},
			}

			res, err := svc.GetGameByID(tt.gameID)

			if tt.expectedError != "" {
				if err == nil {
					t.Fatalf("Se esperaba un error que contuviera '%s', pero err fue nil", tt.expectedError)
				}
				if !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("Se obtuvo el error: '%s'\nSe esperaba que contuviera: '%s'", err.Error(), tt.expectedError)
				}
			} else {
				if err != nil {
					t.Fatalf("No se esperaba ningún error, pero ocurrió: %v", err)
				}
				if tt.validateResult && res == nil {
					t.Error("Se esperaba una respuesta válida (models.GameDetail), se obtuvo nil")
				}
				if tt.validateResult && res.ID != 479836 {
					t.Errorf("El mapeo falló, se esperaba ID 479836, se obtuvo %d", res.ID)
				}
			}
		})
	}
}

func TestGetGameByID_ErrorRealizandoSolicitud(t *testing.T) {
	svc := &RAWGService{
		ApiKey:  "test_key",
		BaseURL: "http://localhost:9999",
		httpClient: &http.Client{
			Timeout: 1 * time.Millisecond,
		},
	}

	_, err := svc.GetGameByID("479836")

	if err == nil {
		t.Fatal("Se esperaba un error de red, pero la solicitud se completó con éxito")
	}

	expectedPrefix := "Error cuando se está realizando la solicitud:"
	if !strings.Contains(err.Error(), expectedPrefix) {
		t.Errorf("Se esperaba un error que contuviera '%s', pero se obtuvo: %v", expectedPrefix, err)
	}
}

func TestGetGameByID_ErrorCreandoSolicitud(t *testing.T) {
	urlCorrupta := "http://api.rawg.io\x7f/v1"

	svc := &RAWGService{
		ApiKey:  "test_key",
		BaseURL: urlCorrupta,
		httpClient: &http.Client{
			Timeout: http.DefaultClient.Timeout,
		},
	}

	_, err := svc.GetGameByID("479836")

	if err == nil {
		t.Fatal("Se esperaba un error al crear la solicitud, pero err fue nil")
	}

	expectedText := "Error cuando se está creando la solicitud:"
	if !strings.Contains(err.Error(), expectedText) {
		t.Errorf("Se esperaba un error que contuviera:\n'%s'\nPero se obtuvo:\n'%s'", expectedText, err.Error())
	}
}
