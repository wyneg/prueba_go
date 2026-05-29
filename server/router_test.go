package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

type AppTest struct {
	mux *http.ServeMux
}

func TestAppTest_HttpMethods(t *testing.T) {
	app := &App{
		mux: http.NewServeMux(),
	}

	handlerEjecutado := false
	var ctxRecibido *Context

	method := http.MethodPost
	path := "/api/v1/games"

	app.HttpMethods(method, path, func(c *Context) {
		handlerEjecutado = true
		ctxRecibido = c

		c.SendText("Ruta alcanzada con éxito")
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(method, path, nil)

	type ctxKey string
	r = r.WithContext(context.WithValue(r.Context(), ctxKey("test-key"), "test-value"))

	app.mux.ServeHTTP(w, r)

	if !handlerEjecutado {
		t.Errorf("El handler registrado mediante HttpMethods no fue ejecutado por el multiplexor")
	}

	if ctxRecibido == nil {
		t.Fatalf("El handler se ejecutó pero recibió un *Context de valor nil")
	}

	if ctxRecibido.ResponseWriter != w {
		t.Errorf("El Context no retuvo el ResponseWriter original de la petición")
	}

	if ctxRecibido.Request != r {
		t.Errorf("El Context no retuvo el Request original de la petición")
	}

	val := ctxRecibido.Cxt.Value(ctxKey("test-key"))
	if val != "test-value" {
		t.Errorf("El campo 'Cxt' del Context no heredó los valores de r.Context(). Se obtuvo: %v", val)
	}

	if w.Body.String() != "Ruta alcanzada con éxito" {
		t.Errorf("Se esperaba la respuesta del handler 'Ruta alcanzada con éxito', se obtuvo: '%s'", w.Body.String())
	}
}
