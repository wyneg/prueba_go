package services

import (
	"context"
	"fmt"

	"github.com/wyneg/prueba_go/models"
	"github.com/wyneg/prueba_go/repositories"
)

type DBService struct {
	repo *repositories.RestRepository
}

func NewDBService(repo *repositories.RestRepository) *DBService {
	return &DBService{repo: repo}
}

func (db *DBService) CreateGame(ctx context.Context, game *models.GameLibrary) error {

	if game.RawgID == nil {
		return fmt.Errorf("Game no se pudo crear, no envío RawgID")
	}

	err := db.repo.Create(ctx, game)

	if err != nil {
		return fmt.Errorf("Game no se pudo crear")
	}

	return nil
}
