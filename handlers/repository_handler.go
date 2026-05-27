package handlers

import (
	"net/http"
	"strconv"

	"github.com/wyneg/prueba_go/models"
	"github.com/wyneg/prueba_go/server"
	"github.com/wyneg/prueba_go/services"
)

type RepositoryHandler struct {
	dbService *services.DBService
}

func NewRepositoryHandler(dbService *services.DBService) *RepositoryHandler {
	return &RepositoryHandler{dbService: dbService}
}

func (r *RepositoryHandler) CreateGameHandler(c *server.Context) {

	var request models.GameLibrary

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"code":  strconv.Itoa(http.StatusBadRequest),
			"error": "Datos inválidos",
		})
		return

	}

	err := r.dbService.CreateGame(c.Context(), &request)

	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"code":  strconv.Itoa(http.StatusInternalServerError),
			"error": "Error al crear el juego en handler",
		})
		return
	}
}
