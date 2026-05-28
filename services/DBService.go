package services

import (
	"context"
	"fmt"

	"github.com/wyneg/prueba_go/models"
)

type GameRepository interface {
	Create(ctx context.Context, game *models.GameLibrary) error
	Select(ctx context.Context, status string) ([]models.GameLibrary, error)
	Update(cxt context.Context, id uint, game *models.GameLibrary) error
	Delete(cxt context.Context, id uint) error
	Stats(ctx context.Context) (models.GameStatsResponse, error)
}

type DBService struct {
	repo GameRepository
}

func NewDBService(repo GameRepository) *DBService {
	return &DBService{repo: repo}
}

func (db *DBService) CreateGame(ctx context.Context, game *models.GameLibrary) error {

	if game.RawgID == nil {
		return fmt.Errorf("Game no se pudo crear, no envió RawgID")
	}

	err := db.repo.Create(ctx, game)

	if err != nil {
		return fmt.Errorf("Game no se pudo crear: %s", err)
	}

	return nil
}

func (db *DBService) GetGame(ctx context.Context, status string) (interface{}, error) {

	result, err := db.repo.Select(ctx, status)

	if err != nil {
		return nil, fmt.Errorf("Error: %s", err)
	}

	return result, nil
}

func (db *DBService) UpdateGame(ctx context.Context, id uint, game *models.GameLibrary) error {

	err := db.repo.Update(ctx, id, game)

	if err != nil {
		return fmt.Errorf("Error: %s", err)
	}

	return nil
}

func (db *DBService) DeleteGame(ctx context.Context, id uint) error {
	err := db.repo.Delete(ctx, id)

	if err != nil {
		return fmt.Errorf("Error: %s", err)
	}
	return nil
}

func (db *DBService) StatsGames(ctx context.Context) (models.GameStatsResponse, error) {
	stats, err := db.repo.Stats(ctx)

	if err != nil {
		return models.GameStatsResponse{}, fmt.Errorf("Error: %s", err)
	}
	return stats, nil
}
