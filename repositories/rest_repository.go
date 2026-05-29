package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/wyneg/prueba_go/models"
)

type DBConnection interface {
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
}

type RestRepository struct {
	db DBConnection
}

func NewRestRepository(db DBConnection) *RestRepository {
	return &RestRepository{db: db}
}

// type RestRepository struct {
// 	db *pgx.Conn
// }

// func NewRestRepository(db *pgx.Conn) *RestRepository {
// 	return &RestRepository{db: db}
// }

func (p *RestRepository) Create(cxt context.Context, game *models.GameLibrary) error {
	query := "INSERT INTO game_library (rawg_id, title, genre, platform, cover_url) VALUES ($1, $2, $3, $4, $5) RETURNING id, added_at"

	var lastInsertID int64
	var addedAt time.Time

	err := p.db.QueryRow(cxt, query, game.RawgID, game.Title, game.Genre, game.Platform, game.CoverURL).Scan(&lastInsertID, &addedAt)

	if err != nil {
		return fmt.Errorf("Error al crear juego: %w", err)
	}

	game.ID = uint(lastInsertID)
	game.AddedAt = addedAt
	return nil
}

func (p *RestRepository) Select(cxt context.Context, status string) ([]models.GameLibrary, error) {
	query := "SELECT * FROM game_library"

	var args []interface{}

	if status != "" {
		query += " WHERE status = $1"
		args = append(args, status)
	}

	result, err := p.db.Query(cxt, query, args...)

	if err != nil {
		return nil, fmt.Errorf("Error al traer juego(s): %w", err)
	}

	defer result.Close()

	var games []models.GameLibrary

	for result.Next() {
		var game models.GameLibrary

		err := result.Scan(&game.ID, &game.RawgID, &game.Title, &game.Genre, &game.Platform, &game.CoverURL, &game.PersonalNote, &game.PersonalScore, &game.Status, &game.AddedAt)

		if err != nil {
			return nil, fmt.Errorf("Error al escanear juego(s): %w", err)
		}
		games = append(games, game)
	}

	return games, nil
}

func (p *RestRepository) Update(cxt context.Context, id uint, game *models.GameLibrary) error {

	var setClauses []string
	var args []interface{}
	argCount := 1

	if game.PersonalNote != nil {
		setClauses = append(setClauses, fmt.Sprintf("personal_note = $%d", argCount))
		args = append(args, *game.PersonalNote)
		argCount++
	}

	if game.PersonalScore != nil {
		setClauses = append(setClauses, fmt.Sprintf("personal_score = $%d", argCount))
		args = append(args, *game.PersonalScore)
		argCount++
	}

	if game.Status != nil {
		setClauses = append(setClauses, fmt.Sprintf("status = $%d", argCount))
		args = append(args, *game.Status)
		argCount++
	}

	query := fmt.Sprintf("UPDATE game_library SET %s WHERE id = $%d", strings.Join(setClauses, ", "), argCount)

	args = append(args, id)

	result, err := p.db.Exec(cxt, query, args...)

	if err != nil {
		return fmt.Errorf("Error al actualizar juego: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("No se encontró el ID")
	}

	return nil
}

func (p *RestRepository) Delete(cxt context.Context, id uint) error {
	query := "DELETE FROM game_library WHERE id = $1"

	result, err := p.db.Exec(cxt, query, id)

	if err != nil {
		return fmt.Errorf("Error al eliminar juego: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("No se encontró el juego")
	}

	return nil
}

func (p *RestRepository) Stats(cxt context.Context) (models.GameStatsResponse, error) {

	var stats models.GameStatsResponse

	queryStatus := "SELECT status, COUNT(*) FROM game_library WHERE status IN ('completado', 'jugando', 'pendiente', 'abandonado') GROUP BY status"

	queryAverage := "SELECT COUNT(*) AS total, ROUND(COALESCE(AVG(personal_score), 0)::numeric, 1) AS promedio FROM game_library"

	status, err := p.db.Query(cxt, queryStatus)

	if err != nil {
		return models.GameStatsResponse{}, fmt.Errorf("Error al obtener estatus: %w", err)
	}

	for status.Next() {
		var statusValue string
		var statusCount uint

		err := status.Scan(&statusValue, &statusCount)

		if err != nil {
			return models.GameStatsResponse{}, fmt.Errorf("Error al escanear estatus: %w", err)
		}

		switch statusValue {
		case "completado":
			stats.ByStatus.Completado = statusCount
		case "jugando":
			stats.ByStatus.Jugando = statusCount
		case "pendiente":
			stats.ByStatus.Pendiente = statusCount
		case "abandonado":
			stats.ByStatus.Abandonado = statusCount
		}
	}

	err = p.db.QueryRow(cxt, queryAverage).Scan(&stats.Total, &stats.AverageScore)

	if err != nil {
		return models.GameStatsResponse{}, fmt.Errorf("Error al obtener promedio: %w", err)
	}

	return stats, nil

}
