package repositories

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v3"
	"github.com/wyneg/prueba_go/models"
)

func ptr[T any](v T) *T {
	return &v
}

func TestCreate(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("error creando el mock: %v", err)
	}
	defer mock.Close(context.Background())

	repo := &RestRepository{db: mock}
	ctx := context.Background()

	game := &models.GameLibrary{
		RawgID:   ptr(uint(12345)),
		Title:    "The Witcher 3",
		Genre:    "RPG",
		Platform: "PC",
		CoverURL: "http://image.com",
	}

	expectedID := int64(1)
	expectedTime := time.Now()

	mock.ExpectQuery("INSERT INTO game_library").
		WithArgs(game.RawgID, game.Title, game.Genre, game.Platform, game.CoverURL).
		WillReturnRows(pgxmock.NewRows([]string{"id", "added_at"}).AddRow(expectedID, expectedTime))

	err = repo.Create(ctx, game)
	if err != nil {
		t.Errorf("no se esperaba error: %v", err)
	}

	if game.ID != uint(expectedID) {
		t.Errorf("se esperaba ID %d, se obtuvo %d", expectedID, game.ID)
	}

	if !game.AddedAt.Equal(expectedTime) {
		t.Errorf("se esperaba AddedAt %v, se obtuvo %v", expectedTime, game.AddedAt)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("quedaron expectativas sin cumplir: %v", err)
	}

	t.Run("Create con error de base de datos", func(t *testing.T) {
		mock.ExpectQuery("INSERT INTO game_library").
			WithArgs(game.RawgID, game.Title, game.Genre, game.Platform, game.CoverURL).
			WillReturnError(fmt.Errorf("db connection timeout"))

		err := repo.Create(ctx, game)
		if err == nil || !strings.Contains(err.Error(), "Error al crear juego") {
			t.Errorf("se esperaba un error de creación, se obtuvo: %v", err)
		}
	})
}

func TestSelect(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("error creando el mock: %v", err)
	}
	defer mock.Close(context.Background())

	repo := &RestRepository{db: mock}
	ctx := context.Background()

	columns := []string{"id", "rawg_id", "title", "genre", "platform", "cover_url", "personal_note", "personal_score", "status", "added_at"}
	now := time.Now()

	t.Run("Select con filtro de status", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM game_library WHERE status = \\$1").
			WithArgs("jugando").
			WillReturnRows(pgxmock.NewRows(columns).
				AddRow(uint(1), ptr(uint(12)), "Zelda", "Action", "Switch", "url", ptr("Nota"), ptr(9), ptr("jugando"), now))

		games, err := repo.Select(ctx, "jugando")
		if err != nil {
			t.Fatalf("error inesperado: %v", err)
		}

		if len(games) != 1 || *games[0].Status != "jugando" {
			t.Errorf("juegos devueltos incorrectos: %+v", games)
		}
	})

	t.Run("Select sin filtros", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM game_library").
			WillReturnRows(pgxmock.NewRows(columns).
				AddRow(uint(1), ptr(uint(12)), "Zelda", "Action", "Switch", "url", ptr("Nota"), ptr(9), ptr("jugando"), now).
				AddRow(uint(2), ptr(uint(13)), "Doom", "FPS", "PC", "url", ptr("Nota 2"), ptr(10), ptr("completado"), now))

		games, err := repo.Select(ctx, "")
		if err != nil {
			t.Fatalf("error inesperado: %v", err)
		}

		if len(games) != 2 {
			t.Errorf("se esperaban 2 juegos, se obtuvieron %d", len(games))
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("quedaron expectativas sin cumplir: %v", err)
	}

	t.Run("Select con error en Query", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM game_library").
			WillReturnError(fmt.Errorf("permiso denegado"))

		_, err := repo.Select(ctx, "")
		if err == nil || !strings.Contains(err.Error(), "Error al traer juego(s)") {
			t.Errorf("se esperaba error de query, se obtuvo: %v", err)
		}
	})

	t.Run("Select con error al escanear filas", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM game_library").
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(uint(1)))

		_, err := repo.Select(ctx, "")
		if err == nil || !strings.Contains(err.Error(), "Error al escanear juego(s)") {
			t.Errorf("se esperaba error de escaneo, se obtuvo: %v", err)
		}
	})
}

func TestUpdate(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("error creando el mock: %v", err)
	}
	defer mock.Close(context.Background())

	repo := &RestRepository{db: mock}
	ctx := context.Background()

	t.Run("Update exitoso", func(t *testing.T) {
		game := &models.GameLibrary{
			PersonalNote:  ptr("Me encantó"),
			PersonalScore: ptr(8),
			Status:        ptr("completado"),
		}

		mock.ExpectExec("UPDATE game_library SET personal_note = \\$1, personal_score = \\$2, status = \\$3 WHERE id = \\$4").
			WithArgs("Me encantó", 8, "completado", uint(1)).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err := repo.Update(ctx, 1, game)
		if err != nil {
			t.Errorf("no se esperaba error: %v", err)
		}
	})

	t.Run("Update ID no encontrado", func(t *testing.T) {
		game := &models.GameLibrary{
			Status: ptr("abandonado"),
		}

		mock.ExpectExec("UPDATE game_library SET status = \\$1 WHERE id = \\$2").
			WithArgs("abandonado", uint(99)).
			WillReturnResult(pgxmock.NewResult("UPDATE", 0))

		err := repo.Update(ctx, 99, game)
		if err == nil || err.Error() != "No se encontró el ID" {
			t.Errorf("se esperaba error 'No se encontró el ID', se obtuvo: %v", err)
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("quedaron expectativas sin cumplir: %v", err)
	}

	t.Run("Update con error de base de datos", func(t *testing.T) {
		game := &models.GameLibrary{Status: ptr("completado")}

		mock.ExpectExec("UPDATE game_library").
			WithArgs("completado", uint(1)).
			WillReturnError(fmt.Errorf("db crash"))

		err := repo.Update(ctx, 1, game)
		if err == nil || !strings.Contains(err.Error(), "Error al actualizar juego") {
			t.Errorf("se esperaba error de actualización, se obtuvo: %v", err)
		}
	})

}

func TestDelete(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("error creando el mock: %v", err)
	}
	defer mock.Close(context.Background())

	repo := &RestRepository{db: mock}
	ctx := context.Background()

	t.Run("Delete exitoso", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM game_library WHERE id = \\$1").
			WithArgs(uint(1)).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))

		err := repo.Delete(ctx, 1)
		if err != nil {
			t.Errorf("error inesperado: %v", err)
		}
	})

	t.Run("Delete no encontrado", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM game_library WHERE id = \\$1").
			WithArgs(uint(55)).
			WillReturnResult(pgxmock.NewResult("DELETE", 0))

		err := repo.Delete(ctx, 55)
		if err == nil || err.Error() != "No se encontró el juego" {
			t.Errorf("se esperaba error 'No se encontró el juego', se obtuvo: %v", err)
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("quedaron expectativas sin cumplir: %v", err)
	}

	t.Run("Delete con error de base de datos", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM game_library").
			WithArgs(uint(1)).
			WillReturnError(fmt.Errorf("foreign key constraint violation"))

		err := repo.Delete(ctx, 1)
		if err == nil || !strings.Contains(err.Error(), "Error al eliminar juego") {
			t.Errorf("se esperaba error de eliminación, se obtuvo: %v", err)
		}
	})
}

func TestStats(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("error creando el mock: %v", err)
	}
	defer mock.Close(context.Background())

	mock.MatchExpectationsInOrder(true)

	repo := &RestRepository{db: mock}
	ctx := context.Background()

	statusColumns := []string{"status", "count"}
	queryStatus := "SELECT status, COUNT(*) FROM game_library WHERE status IN ('completado', 'jugando', 'pendiente', 'abandonado') GROUP BY status"
	mock.ExpectQuery(regexp.QuoteMeta(queryStatus)).
		WillReturnRows(pgxmock.NewRows(statusColumns).
			AddRow("completado", uint(5)).
			AddRow("jugando", uint(2)).
			AddRow("pendiente", uint(1)).
			AddRow("abandonado", uint(1)))

	avgColumns := []string{"total", "promedio"}
	queryAverage := "SELECT COUNT(*) AS total, ROUND(COALESCE(AVG(personal_score), 0)::numeric, 1) AS promedio FROM game_library"
	mock.ExpectQuery(regexp.QuoteMeta(queryAverage)).
		WillReturnRows(pgxmock.NewRows(avgColumns).AddRow(uint(7), float64(8.5)))

	stats, err := repo.Stats(ctx)
	if err != nil {
		t.Fatalf("error inesperado en Stats: %v", err)
	}

	if stats.Total != 7 {
		t.Errorf("se esperaba un total de 7, se obtuvo %d", stats.Total)
	}

	if stats.AverageScore != 8.5 {
		t.Errorf("se esperaba un promedio de 8.5, se obtuvo %f", stats.AverageScore)
	}

	if stats.ByStatus.Completado != 5 || stats.ByStatus.Jugando != 2 {
		t.Errorf("los contadores por estatus están mal calculados: %+v", stats.ByStatus)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("quedaron expectativas sin cumplir: %v", err)
	}

	t.Run("Error al obtener estatus (QueryStatus)", func(t *testing.T) {
		queryStatus := "SELECT status, COUNT\\(\\*\\) FROM game_library"
		mock.ExpectQuery(queryStatus).WillReturnError(fmt.Errorf("error de red"))

		_, err := repo.Stats(ctx)
		if err == nil || !strings.Contains(err.Error(), "Error al obtener estatus") {
			t.Errorf("se esperaba error de obtención de estatus, se obtuvo: %v", err)
		}
	})

	t.Run("Error al escanear estatus (Scan)", func(t *testing.T) {
		queryStatus := "SELECT status, COUNT\\(\\*\\) FROM game_library"

		mock.ExpectQuery(queryStatus).
			WillReturnRows(pgxmock.NewRows([]string{"solo_una_columna"}).AddRow("completado"))

		_, err := repo.Stats(ctx)
		if err == nil || !strings.Contains(err.Error(), "Error al escanear estatus") {
			t.Errorf("se esperaba error de escaneo de estatus, se obtuvo: %v", err)
		}
	})

	t.Run("Error al obtener promedio (QueryAverage)", func(t *testing.T) {

		queryStatus := "SELECT status, COUNT\\(\\*\\) FROM game_library"
		mock.ExpectQuery(queryStatus).WillReturnRows(pgxmock.NewRows([]string{"status", "count"}).AddRow("completado", uint(1)))

		queryAverage := "SELECT COUNT\\(\\*\\) AS total"
		mock.ExpectQuery(queryAverage).WillReturnError(fmt.Errorf("tabla bloqueada"))

		_, err := repo.Stats(ctx)
		if err == nil || !strings.Contains(err.Error(), "Error al obtener promedio") {
			t.Errorf("se esperaba error de promedio, se obtuvo: %v", err)
		}
	})
}

func TestNewRestRepository(t *testing.T) {
	mock, _ := pgxmock.NewConn()
	repo := NewRestRepository(mock)

	if repo == nil {
		t.Fatal("se esperaba que el repositorio no fuera nil")
	}
	if repo.db != mock {
		t.Error("el mock de la base de datos no se asignó correctamente")
	}
}
