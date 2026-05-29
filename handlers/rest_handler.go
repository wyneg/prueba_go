package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/wyneg/prueba_go/models"
	"github.com/wyneg/prueba_go/server"
	"github.com/wyneg/prueba_go/services"
)

type RestHandler struct {
	rawgService *services.RAWGService
}

func NewRestHandler(rawgService *services.RAWGService) *RestHandler {
	return &RestHandler{rawgService: rawgService}
}

func (r *RestHandler) GetGameHandler(c *server.Context) {

	gameName := c.Request.URL.Query().Get("q")

	if gameName == "" {
		c.JSON(http.StatusBadRequest, models.NewBadRequestError("El parámetro de consulta 'q' es requerido"))
		return
	}

	game, err := r.rawgService.GetGame(gameName)

	if err != nil {
		if strings.HasPrefix(err.Error(), "Error cuando se está") {
			c.JSON(http.StatusBadRequest, models.NewBadRequestError(err.Error()))
			return
		}

		codeError := err.Error()[1:4]
		descriptionError := err.Error()[6:]

		code, _ := strconv.Atoi(codeError)

		c.JSON(code, models.NewError(code, descriptionError))
		return
	}

	c.JSON(http.StatusOK, game)

}

func (r *RestHandler) GetGameByIDHandler(c *server.Context) {
	id := c.Request.PathValue("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, models.NewBadRequestError("El parámetro de consulta 'id' es requerido"))
		return
	}

	game, err := r.rawgService.GetGameByID(id)

	if err != nil {
		if strings.HasPrefix(err.Error(), "Error cuando se está") {
			c.JSON(http.StatusBadRequest, models.NewBadRequestError(err.Error()))
			return
		}

		codeError := err.Error()[1:4]
		descriptionError := err.Error()[6:]

		code, _ := strconv.Atoi(codeError)

		c.JSON(code, models.NewError(code, descriptionError))
		return
	}

	c.JSON(http.StatusOK, game)
}
