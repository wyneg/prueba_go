package handlers

import (
	"net/http"
	"strconv"
	"strings"

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
		c.JSON(http.StatusBadRequest, models.NewBadRequestError("Datos inválidos"))
		return
	}

	if request.RawgID == nil || *request.RawgID == 0 {
		c.JSON(http.StatusBadRequest, models.NewBadRequestError("El campo RAWG_ID es obligatorio"))
		return
	}

	if request.Title == "" {
		c.JSON(http.StatusBadRequest, models.NewBadRequestError("El campo Title es obligatorio"))
		return
	}

	err := r.dbService.CreateGame(c.Context(), &request)

	if err != nil {
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			c.JSON(http.StatusConflict, models.NewConflictError("El RAWG_ID de juego ya existe"))
			return
		}

		c.JSON(http.StatusInternalServerError, models.NewInternalServerError("Error al crear el juego"))
		return
	}

	c.JSON(http.StatusCreated, request)
}

func (r *RepositoryHandler) GetGameHandler(c *server.Context) {

	status := c.Request.URL.Query().Get("status")

	result, err := r.dbService.GetGame(c.Context(), status)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewInternalServerError("Error al obtener los juegos"))
		return

	}

	if gamesSlice, ok := result.([]models.GameLibrary); ok {
		if len(gamesSlice) == 0 {
			c.JSON(http.StatusNotFound, models.NewNotFoundError("Juegos no encontrados en BD"))
			return
		}
	}

	c.JSON(http.StatusOK, result)
}

func (r *RepositoryHandler) UpdateGameHandler(c *server.Context) {
	id := c.Request.PathValue("id")

	parsedUint, err := strconv.ParseUint(id, 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewBadRequestError("ID inválido"))
		return
	}

	var request models.GameLibrary

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.NewBadRequestError("Datos inválidos"))
		return
	}

	if request.PersonalScore == nil && request.PersonalNote == nil && request.Status == nil {
		c.JSON(http.StatusBadRequest, models.NewBadRequestError("Debe enviar al menos uno o más campos para actualizar"))
		return
	}

	if request.PersonalScore != nil && (*request.PersonalScore < 1 || *request.PersonalScore > 10) {
		c.JSON(http.StatusBadRequest, models.NewBadRequestError("La puntuación personal debe estar entre 1 y 10"))
		return
	}

	err = r.dbService.UpdateGame(c.Context(), uint(parsedUint), &request)

	if err != nil {
		if strings.Contains(err.Error(), "No se encontró el ID") {
			c.JSON(http.StatusNotFound, models.NewNotFoundError("ID no encontrado"))
			return
		}

		c.JSON(http.StatusInternalServerError, models.NewInternalServerError("Error al actualizar el juego"))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (r *RepositoryHandler) DeleteGameHandler(c *server.Context) {
	id := c.Request.PathValue("id")

	parsedUint, err := strconv.ParseUint(id, 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewBadRequestError("ID inválido"))
		return
	}

	err = r.dbService.DeleteGame(c.Context(), uint(parsedUint))

	if err != nil {
		if strings.Contains(err.Error(), "No se encontró el juego") {
			c.JSON(http.StatusNotFound, models.NewNotFoundError("Juego no encontrado"))
			return
		}

		c.JSON(http.StatusInternalServerError, models.NewInternalServerError("Error al eliminar el juego"))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (r *RepositoryHandler) StatsGameHandler(c *server.Context) {
	stats, err := r.dbService.StatsGames(c.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewInternalServerError("Error al obtener estadísticas"))
		return
	}

	c.JSON(http.StatusOK, stats)
}
