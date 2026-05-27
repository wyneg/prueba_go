package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/wyneg/prueba_go/models"
)

type RestRepository struct {
	db *pgx.Conn
}

func NewRestRepository(db *pgx.Conn) *RestRepository {
	return &RestRepository{db: db}
}

func (p *RestRepository) Create(cxt context.Context, game *models.GameLibrary) error {
	query := "INSERT INTO game_library (rawg_id, title, genre, platform, cover_url) VALUES ($1, $2, $3, $4, $5) RETURNING id"

	var lastInsertID int64

	err := p.db.QueryRow(cxt, query, game.RawgID, game.Title, game.Genre, game.Platform, game.CoverURL).Scan(&lastInsertID)

	if err != nil {
		return fmt.Errorf("Error al crear juego: %w", err)
	}

	game.ID = uint(lastInsertID)
	return nil
}
