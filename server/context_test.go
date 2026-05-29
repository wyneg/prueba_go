package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testPayload struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
}

func TestContext(t *testing.T) {

	t.Run("SendText - Escribe texto plano correctamente", func(t *testing.T) {
		w := httptest.NewRecorder()
		c := &Context{ResponseWriter: w}

		expectedText := "Hola Mundo desde Go"
		c.SendText(expectedText)

		if w.Body.String() != expectedText {
			t.Errorf("se esperaba el texto '%s', se obtuvo '%s'", expectedText, w.Body.String())
		}
	})

	t.Run("Status - Configura el código HTTP de respuesta", func(t *testing.T) {
		w := httptest.NewRecorder()
		c := &Context{ResponseWriter: w}

		c.Status(http.StatusTeapot)

		if w.Code != http.StatusTeapot {
			t.Errorf("se esperaba el código de estado %d, se obtuvo %d", http.StatusTeapot, w.Code)
		}
	})

	t.Run("JSON - Serializa correctamente objetos y cabeceras", func(t *testing.T) {
		w := httptest.NewRecorder()
		c := &Context{ResponseWriter: w}

		payload := testPayload{Message: "éxito", Status: true}
		err := c.JSON(http.StatusCreated, payload)

		if err != nil {
			t.Fatalf("no se esperaba un error al codificar JSON: %v", err)
		}

		contentType := w.Header().Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("se esperaba Content-Type 'application/json', se obtuvo '%s'", contentType)
		}

		if w.Code != http.StatusCreated {
			t.Errorf("se esperaba código %d, se obtuvo %d", http.StatusCreated, w.Code)
		}

		var response testPayload
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("el cuerpo de la respuesta está malformado: %v", err)
		}

		if response.Message != payload.Message || response.Status != payload.Status {
			t.Errorf("los datos serializados no coinciden con el objeto original")
		}
	})

	t.Run("BindJSON - Decodifica cuerpos JSON entrantes de peticiones", func(t *testing.T) {
		bodyBytes := []byte(`{"message": "petición entrante", "status": false}`)
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(bodyBytes))
		c := &Context{Request: req}

		var dest testPayload
		err := c.BindJSON(&dest)

		if err != nil {
			t.Fatalf("no se esperaba un error al decodificar: %v", err)
		}

		if dest.Message != "petición entrante" || dest.Status != false {
			t.Errorf("BindJSON no mapeó correctamente las propiedades del cuerpo de la solicitud")
		}
	})

	t.Run("BindJSON - Retorna error ante JSON inválido", func(t *testing.T) {
		bodyBytes := []byte(`{"message": "incompleto"`)
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(bodyBytes))
		c := &Context{Request: req}

		var dest testPayload
		err := c.BindJSON(&dest)

		if err == nil {
			t.Error("se esperaba un error de decodificación debido a un JSON malformado, se obtuvo nil")
		}
	})

	t.Run("SetUserID y GetUserID - Almacena y lee datos de sesión del contexto", func(t *testing.T) {
		c := &Context{}
		var expectedID uint = 42

		c.SetUserID(expectedID)
		actualID := c.GetUserID()

		if actualID != expectedID {
			t.Errorf("se esperaba el ID de usuario %d, se obtuvo %d", expectedID, actualID)
		}
	})

	t.Run("Context - Devuelve la instancia interna de context.Context", func(t *testing.T) {
		type ctxKey string
		key := ctxKey("mi-clave")

		nativeCtx := context.WithValue(context.Background(), key, "valor-secreto")
		c := &Context{Cxt: nativeCtx}

		returnedCtx := c.Context()
		val := returnedCtx.Value(key)

		if val != "valor-secreto" {
			t.Errorf("el contexto devuelto perdió la información del contexto original")
		}
	})
}
