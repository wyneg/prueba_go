package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/wyneg/prueba_go/models"

	"github.com/joho/godotenv"
)

type RAWGService struct {
	ApiKey     string
	BaseURL    string
	httpClient *http.Client
}

func NewRAWGService(apiKey string, baseURL string, filenames ...string) *RAWGService {

	err := godotenv.Load(filenames...)

	if err != nil {
		log.Fatal("Error cargando archivo .env")
	}

	return &RAWGService{
		ApiKey:  apiKey,
		BaseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *RAWGService) GetGame(gameName string) (*models.RAWGResponse, error) {

	url := s.BaseURL + "/games?key=" + s.ApiKey + "&search=" + gameName

	//Formular request
	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, fmt.Errorf("Error cuando se está creando la solicitud: %v", err)
	}

	//Formuular response
	response, err := s.httpClient.Do(request)

	if err != nil {
		return nil, fmt.Errorf("Error cuando se está realizando la solicitud: %v", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusNotFound {
			return nil, models.NewNotFoundError("Juego no encontrado")
		}
		if response.StatusCode == http.StatusBadRequest {
			return nil, models.NewBadRequestError("Solicitud incorrecta")
		}

	}

	var rawgResponse models.RAWGResponse

	err = json.NewDecoder(response.Body).Decode(&rawgResponse)

	if err != nil {
		return nil, fmt.Errorf("Error decodificando la respuesta: %v", err)
	}

	return &rawgResponse, nil
}

func (s *RAWGService) GetGameByID(gameName string) (*models.GameDetail, error) {
	url := s.BaseURL + "/games/" + gameName + "?key=" + s.ApiKey

	//Formular request
	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, fmt.Errorf("Error cuando se está creando la solicitud: %v", err)
	}

	//Formuular response
	response, err := s.httpClient.Do(request)

	if err != nil {
		return nil, fmt.Errorf("Error cuando se está realizando la solicitud: %v", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusNotFound {
			return nil, models.NewNotFoundError("Juego no encontrado")
		}
		if response.StatusCode == http.StatusBadRequest {
			return nil, models.NewBadRequestError("Solicitud incorrecta")
		}

	}

	var gameDetail models.GameDetail

	err = json.NewDecoder(response.Body).Decode(&gameDetail)

	if err != nil {
		return nil, fmt.Errorf("Error decodificando la respuesta: %v", err)
	}

	return &gameDetail, nil
}
