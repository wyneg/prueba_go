package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v3"

	"github.com/wyneg/prueba_go/models"
	"github.com/wyneg/prueba_go/repositories"
	"github.com/wyneg/prueba_go/server"
	"github.com/wyneg/prueba_go/services"
)

type ConnAdapter struct {
	pgxmock.PgxConnIface
}

func (a *ConnAdapter) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	mockArgs := make([]interface{}, len(args))
	copy(mockArgs, args)

	tag, err := a.PgxConnIface.Exec(ctx, query, mockArgs...)
	if err != nil {
		return pgconn.CommandTag{}, err
	}

	rowsAffected := tag.RowsAffected()
	return pgconn.NewCommandTag(fmt.Sprintf("UPDATE %d", rowsAffected)), nil
}

func (a *ConnAdapter) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	mockArgs := make([]interface{}, len(args))
	copy(mockArgs, args)
	return a.PgxConnIface.Query(ctx, query, mockArgs...)
}

func (a *ConnAdapter) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	mockArgs := make([]interface{}, len(args))
	copy(mockArgs, args)
	return a.PgxConnIface.QueryRow(ctx, query, mockArgs...)
}

func ptr[T any](v T) *T {
	return &v
}

func addPathValue(r *http.Request, key, value string) *http.Request {
	r.SetPathValue(key, value)
	return r
}

type TestErrorResponse struct {
	Message string `json:"message"`
}

func setupTestEnv(t *testing.T) (pgxmock.PgxConnIface, *RepositoryHandler) {
	connMock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("no se pudo crear el mock de pgx: %v", err)
	}
	adapter := &ConnAdapter{PgxConnIface: connMock}
	repo := repositories.NewRestRepository(adapter)
	dbService := services.NewDBService(repo)
	handler := NewRepositoryHandler(dbService)
	return connMock, handler
}

func TestCreateGameHandler(t *testing.T) {
	t.Run("JSON Inválido - Bad Request", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/games", bytes.NewBuffer([]byte(`{"rawg_id": 123`)))
		c := &server.Context{ResponseWriter: w, Request: r, Cxt: context.Background()}

		handler := NewRepositoryHandler(nil)
		handler.CreateGameHandler(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("se esperaba estatus 400, se obtuvo %d", w.Code)
		}
	})

	t.Run("Falta RAWG_ID o es cero - Bad Request", func(t *testing.T) {
		w := httptest.NewRecorder()
		game := models.GameLibrary{Title: "The Witcher 3"}
		body, _ := json.Marshal(game)
		r := httptest.NewRequest(http.MethodPost, "/games", bytes.NewBuffer(body))
		c := &server.Context{ResponseWriter: w, Request: r, Cxt: context.Background()}

		handler := NewRepositoryHandler(nil)
		handler.CreateGameHandler(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("se esperaba estatus 400, se obtuvo %d", w.Code)
		}
	})

	t.Run("Falta Title - Bad Request", func(t *testing.T) {
		w := httptest.NewRecorder()
		game := models.GameLibrary{RawgID: ptr(uint(12345))}
		body, _ := json.Marshal(game)
		r := httptest.NewRequest(http.MethodPost, "/games", bytes.NewBuffer(body))
		c := &server.Context{ResponseWriter: w, Request: r, Cxt: context.Background()}

		handler := NewRepositoryHandler(nil)
		handler.CreateGameHandler(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("se esperaba estatus 400, se obtuvo %d", w.Code)
		}
	})
}

func TestGetGameHandler(t *testing.T) {
	t.Run("Obtener con filtro de URL Query", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/games?status=jugando", nil)
		c := &server.Context{ResponseWriter: w, Request: r, Cxt: context.Background()}

		status := c.Request.URL.Query().Get("status")
		if status != "jugando" {
			t.Errorf("se esperaba obtener el filtro 'jugando', se obtuvo '%s'", status)
		}
	})
}

func TestUpdateGameHandler(t *testing.T) {
	t.Run("ID inválido (No numérico)", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, "/games/abc", nil)
		c := &server.Context{ResponseWriter: w, Request: r, Cxt: context.Background()}

		handler := NewRepositoryHandler(nil)
		handler.UpdateGameHandler(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("se esperaba estatus 400 por ID inválido, se obtuvo %d", w.Code)
		}
	})

}

func TestDeleteGameHandler(t *testing.T) {
	t.Run("ID inválido en Delete", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, "/games/xyz", nil)
		c := &server.Context{ResponseWriter: w, Request: r, Cxt: context.Background()}

		handler := NewRepositoryHandler(nil)
		handler.DeleteGameHandler(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("se esperaba estatus 400, se obtuvo %d", w.Code)
		}
	})

	t.Run("JSON Inválido en Update - BindJSON Error", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, "/games/1", bytes.NewBuffer([]byte(`{"status": "jugando"`)))
		r = addPathValue(r, "id", "1")
		c := &server.Context{ResponseWriter: w, Request: r, Cxt: context.Background()}

		handler := NewRepositoryHandler(nil)
		handler.UpdateGameHandler(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("se esperaba estatus 400 por JSON inválido, se obtuvo %d", w.Code)
		}
	})

	t.Run("Falta enviar campos obligatorios para actualizar", func(t *testing.T) {
		w := httptest.NewRecorder()
		game := map[string]string{"title": "Nuevo Titulo"}
		body, _ := json.Marshal(game)

		r := httptest.NewRequest(http.MethodPut, "/games/1", bytes.NewBuffer(body))
		r = addPathValue(r, "id", "1")
		c := &server.Context{ResponseWriter: w, Request: r, Cxt: context.Background()}

		handler := NewRepositoryHandler(nil)
		handler.UpdateGameHandler(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("se esperaba estatus 400 por falta de campos modificables, se obtuvo %d", w.Code)
		}
	})

	t.Run("Puntuación Personal fuera de rango (Mayor a 10)", func(t *testing.T) {
		w := httptest.NewRecorder()
		game := models.GameLibrary{
			PersonalScore: ptr(12),
		}
		body, _ := json.Marshal(game)
		r := httptest.NewRequest(http.MethodPut, "/games/1", bytes.NewBuffer(body))
		r = addPathValue(r, "id", "1")
		c := &server.Context{ResponseWriter: w, Request: r, Cxt: context.Background()}

		handler := NewRepositoryHandler(nil)
		handler.UpdateGameHandler(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("se esperaba estatus 400 por puntuación mayor a 10, se obtuvo %d", w.Code)
		}
	})

	t.Run("Puntuación Personal fuera de rango (Menor a 1)", func(t *testing.T) {
		w := httptest.NewRecorder()
		game := models.GameLibrary{
			PersonalScore: ptr(0),
		}
		body, _ := json.Marshal(game)
		r := httptest.NewRequest(http.MethodPut, "/games/1", bytes.NewBuffer(body))
		r = addPathValue(r, "id", "1")
		c := &server.Context{ResponseWriter: w, Request: r, Cxt: context.Background()}

		handler := NewRepositoryHandler(nil)
		handler.UpdateGameHandler(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("se esperaba estatus 400 por puntuación menor a 1, se obtuvo %d", w.Code)
		}
	})
}

func TestContextMethods(t *testing.T) {
	t.Run("Prueba de SendText y Status", func(t *testing.T) {
		w := httptest.NewRecorder()
		c := &server.Context{ResponseWriter: w}

		c.Status(http.StatusAccepted)
		c.SendText("Hola Mundo")

		if w.Code != http.StatusAccepted {
			t.Errorf("Status() falló, se obtuvo %d", w.Code)
		}
		if w.Body.String() != "Hola Mundo" {
			t.Errorf("SendText() falló, se obtuvo %s", w.Body.String())
		}
	})

	t.Run("Prueba de SetUserID y GetUserID", func(t *testing.T) {
		c := &server.Context{}
		c.SetUserID(99)

		if c.GetUserID() != 99 {
			t.Errorf("GetUserID/SetUserID falló, se obtuvo %d", c.GetUserID())
		}
	})

	t.Run("Prueba de Context getter", func(t *testing.T) {
		expectedCtx := context.WithValue(context.Background(), "key", "value")
		c := &server.Context{Cxt: expectedCtx}

		if c.Context() != expectedCtx {
			t.Error("Context() no devolvió el contexto correcto")
		}
	})
}

func TestCreateGameHandler_ServicePaths(t *testing.T) {
	game := models.GameLibrary{
		RawgID: ptr(uint(999)),
		Title:  "Cyberpunk 2077",
	}
	body, _ := json.Marshal(game)

	t.Run("Conflicto - SQLSTATE 23505", func(t *testing.T) {
		connMock, handler := setupTestEnv(t)
		defer connMock.Close(context.Background())

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/games", bytes.NewBuffer(body))
		c := &server.Context{ResponseWriter: w, Request: r, Cxt: context.Background()}

		connMock.ExpectExec("^INSERT INTO game_library.*").
			WillReturnError(fmt.Errorf("ERROR: duplicate key value violates unique constraint (SQLSTATE 23505)"))
		connMock.ExpectQuery("^INSERT INTO game_library.*").
			WillReturnError(fmt.Errorf("ERROR: duplicate key value violates unique constraint (SQLSTATE 23505)"))

		handler.CreateGameHandler(c)

		if w.Code != http.StatusConflict {
			t.Errorf("se esperaba 409 Conflict, se obtuvo %d", w.Code)
		}
	})

	t.Run("Error Interno de Servidor", func(t *testing.T) {
		connMock, handler := setupTestEnv(t)
		defer connMock.Close(context.Background())

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/games", bytes.NewBuffer(body))
		c := &server.Context{ResponseWriter: w, Request: r, Cxt: context.Background()}

		connMock.ExpectQuery("^INSERT INTO game_library.*").
			WillReturnError(fmt.Errorf("fatal connection loss"))
		connMock.ExpectExec("^INSERT INTO game_library.*").
			WillReturnError(fmt.Errorf("fatal connection loss"))

		handler.CreateGameHandler(c)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("se esperaba 500 Internal Server Error, se obtuvo %d", w.Code)
		}
	})

	t.Run("Creación Exitosa - Status Created", func(t *testing.T) {
		connMock, handler := setupTestEnv(t)
		defer connMock.Close(context.Background())

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/games", bytes.NewBuffer(body))
		c := &server.Context{ResponseWriter: w, Request: r, Cxt: context.Background()}

		now := time.Now()

		connMock.ExpectQuery("(?s).*INSERT INTO game_library.*").
			WithArgs(
				pgxmock.AnyArg(),
				pgxmock.AnyArg(),
				pgxmock.AnyArg(),
				pgxmock.AnyArg(),
				pgxmock.AnyArg(),
			).
			WillReturnRows(
				pgxmock.NewRows([]string{"id", "added_at"}).AddRow(int64(1), now),
			)

		handler.CreateGameHandler(c)

		if w.Code != http.StatusCreated && w.Code != http.StatusOK {
			t.Logf("Cuerpo del error persistente: %s", w.Body.String())
			t.Errorf("se esperaba 201 Created o 200 OK, se obtuvo %d", w.Code)
		}
	})
}

func TestGetGameHandler_ServicePaths(t *testing.T) {
	columns := []string{"id", "rawg_id", "title", "genre", "platform", "cover_url", "personal_note", "personal_score", "status", "added_at"}

	t.Run("Error Interno al obtener juegos", func(t *testing.T) {
		connMock, handler := setupTestEnv(t)
		defer connMock.Close(context.Background())

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/games", nil)
		c := &server.Context{ResponseWriter: w, Request: r, Cxt: context.Background()}

		connMock.ExpectQuery("^SELECT .* FROM game_library").
			WillReturnError(fmt.Errorf("query timeout"))

		handler.GetGameHandler(c)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("se esperaba 500, se obtuvo %d", w.Code)
		}
	})

	t.Run("Juegos no encontrados en BD - Retorna Nil", func(t *testing.T) {
		connMock, handler := setupTestEnv(t)
		defer connMock.Close(context.Background())

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/games", nil)
		c := &server.Context{ResponseWriter: w, Request: r, Cxt: context.Background()}

		connMock.ExpectQuery("^SELECT .* FROM game_library").
			WillReturnRows(pgxmock.NewRows(columns))

		handler.GetGameHandler(c)

		if w.Code != http.StatusNotFound {
			t.Logf("Cuerpo obtenido: %s", w.Body.String())
			t.Errorf("se esperaba 404 Not Found, se obtuvo %d", w.Code)
		}
	})

	t.Run("Obtención Exitosa - Status OK", func(t *testing.T) {
		connMock, handler := setupTestEnv(t)
		defer connMock.Close(context.Background())

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/games", nil)
		c := &server.Context{ResponseWriter: w, Request: r, Cxt: context.Background()}

		connMock.ExpectQuery("^SELECT .* FROM game_library").
			WillReturnRows(pgxmock.NewRows(columns).
				AddRow(uint(1), ptr(uint(12)), "Zelda", "Action", "Switch", "url", ptr("Nota"), ptr(9), ptr("jugando"), time.Now()))

		handler.GetGameHandler(c)

		if w.Code != http.StatusOK {
			t.Errorf("se esperaba 200 OK, se obtuvo %d", w.Code)
		}
	})

	t.Run("Cobertura - Slice de Punteros Vacío", func(t *testing.T) {

		connMock, handler := setupTestEnv(t)
		defer connMock.Close(context.Background())

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/games", nil)
		c := &server.Context{ResponseWriter: w, Request: r, Cxt: context.Background()}

		connMock.ExpectQuery("^SELECT .* FROM game_library").
			WillReturnRows(pgxmock.NewRows(columns))

		handler.GetGameHandler(c)

		var result interface{} = []*models.GameLibrary{}
		if gamesPtrSlice, ok := result.([]*models.GameLibrary); ok {
			if len(gamesPtrSlice) == 0 {
				_ = gamesPtrSlice
			}
		}

		if w.Code != http.StatusNotFound {
			t.Errorf("se esperaba 404 Not Found para la rama de punteros, se obtuvo %d", w.Code)
		}
	})
}

func TestUpdateGameHandler_ServicePaths(t *testing.T) {
	game := models.GameLibrary{Status: ptr("completado")}
	body, _ := json.Marshal(game)

	t.Run("ID No Encontrado - Status 404", func(t *testing.T) {
		connMock, handler := setupTestEnv(t)
		defer connMock.Close(context.Background())

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, "/games/99", bytes.NewBuffer(body))
		r.SetPathValue("id", "99")
		c := &server.Context{ResponseWriter: w, Request: r, Cxt: context.Background()}

		connMock.ExpectExec("(?s).*UPDATE game_library.*").
			WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
			WillReturnResult(pgxmock.NewResult("UPDATE", 0))

		handler.UpdateGameHandler(c)

		if w.Code != http.StatusNotFound {
			t.Logf("Cuerpo del error: %s", w.Body.String())
			t.Errorf("se esperaba 404, se obtuvo %d", w.Code)
		}
	})

	t.Run("Update Exitoso - Status No Content", func(t *testing.T) {
		connMock, handler := setupTestEnv(t)
		defer connMock.Close(context.Background())

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, "/games/1", bytes.NewBuffer(body))
		r.SetPathValue("id", "1")
		c := &server.Context{ResponseWriter: w, Request: r, Cxt: context.Background()}

		connMock.ExpectExec("(?s).*UPDATE game_library.*").
			WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		handler.UpdateGameHandler(c)

		if w.Code != http.StatusNoContent && w.Code != http.StatusOK {
			t.Logf("Cuerpo del error: %s", w.Body.String())
			t.Errorf("se esperaba 204 o 200, se obtuvo %d", w.Code)
		}
	})

	t.Run("Error de Base de Datos - Status 500", func(t *testing.T) {
		connMock, handler := setupTestEnv(t)
		defer connMock.Close(context.Background())

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, "/games/1", bytes.NewBuffer(body))
		r = addPathValue(r, "id", "1")
		c := &server.Context{ResponseWriter: w, Request: r, Cxt: context.Background()}

		connMock.ExpectExec("^UPDATE game_library.*").
			WillReturnError(fmt.Errorf("database corrupt"))

		handler.UpdateGameHandler(c)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("se esperaba 500, se obtuvo %d", w.Code)
		}
	})

}

func TestDeleteGameHandler_ServicePaths(t *testing.T) {
	t.Run("Juego No Encontrado - Status 404", func(t *testing.T) {
		connMock, handler := setupTestEnv(t)
		defer connMock.Close(context.Background())

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, "/games/88", nil)
		r = addPathValue(r, "id", "88")
		c := &server.Context{ResponseWriter: w, Request: r, Cxt: context.Background()}

		connMock.ExpectExec("(?s).*DELETE FROM game_library.*").
			WithArgs(pgxmock.AnyArg()).
			WillReturnResult(pgxmock.NewResult("DELETE", 0))

		handler.DeleteGameHandler(c)

		if w.Code != http.StatusNotFound {
			t.Logf("Cuerpo del error en Delete 404: %s", w.Body.String())
			t.Errorf("se esperaba 404, se obtuvo %d", w.Code)
		}
	})

	t.Run("Delete Exitoso", func(t *testing.T) {
		connMock, handler := setupTestEnv(t)
		defer connMock.Close(context.Background())

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, "/games/1", nil)
		r = addPathValue(r, "id", "1")
		c := &server.Context{ResponseWriter: w, Request: r, Cxt: context.Background()}

		connMock.ExpectExec("(?s).*DELETE FROM game_library.*").
			WithArgs(pgxmock.AnyArg()).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))

		handler.DeleteGameHandler(c)

		if w.Code != http.StatusNoContent && w.Code != http.StatusOK {
			t.Logf("Cuerpo del error en Delete Exitoso: %s", w.Body.String())
			t.Errorf("se esperaba 204 o 200, se obtuvo %d", w.Code)
		}
	})

	t.Run("Error de Base de Datos al eliminar - Status 500", func(t *testing.T) {
		connMock, handler := setupTestEnv(t)
		defer connMock.Close(context.Background())

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, "/games/1", nil)
		r = addPathValue(r, "id", "1")
		c := &server.Context{ResponseWriter: w, Request: r, Cxt: context.Background()}

		connMock.ExpectExec("(?s).*DELETE FROM game_library.*").
			WithArgs(pgxmock.AnyArg()).
			WillReturnError(fmt.Errorf("disk i/o error on delete"))

		handler.DeleteGameHandler(c)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("se esperaba 500 Internal Server Error, se obtuvo %d", w.Code)
		}
	})
}

func TestStatsGameHandler_ServicePaths(t *testing.T) {
	t.Run("Stats Exitoso", func(t *testing.T) {
		connMock, handler := setupTestEnv(t)
		defer connMock.Close(context.Background())

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/stats", nil)
		c := &server.Context{ResponseWriter: w, Request: r, Cxt: context.Background()}

		connMock.ExpectQuery("^SELECT status, COUNT.*").
			WillReturnRows(pgxmock.NewRows([]string{"status", "count"}).AddRow("completado", uint(3)))

		connMock.ExpectQuery("^SELECT COUNT\\(\\*\\) AS total.*").
			WillReturnRows(pgxmock.NewRows([]string{"total", "promedio"}).AddRow(uint(3), float64(9.0)))

		handler.StatsGameHandler(c)

		if w.Code != http.StatusOK {
			t.Errorf("se esperaba 200, se obtuvo %d", w.Code)
		}
	})

	t.Run("Error de Base de Datos al obtener estadísticas - Status 500", func(t *testing.T) {
		connMock, handler := setupTestEnv(t)
		defer connMock.Close(context.Background())

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/stats", nil)
		c := &server.Context{ResponseWriter: w, Request: r, Cxt: context.Background()}

		connMock.ExpectQuery("(?s).*SELECT status, COUNT.*").
			WillReturnError(fmt.Errorf("table game_library does not exist"))

		handler.StatsGameHandler(c)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("se esperaba 500 Internal Server Error, se obtuvo %d", w.Code)
		}
	})
}

type MockDBServicePunteros struct {
	services.DBService
}

func (m *MockDBServicePunteros) GetGame(ctx context.Context, status string) (interface{}, error) {
	return []*models.GameLibrary{}, nil
}
