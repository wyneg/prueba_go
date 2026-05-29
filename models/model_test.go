package models

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"
)

func ptr[T any](v T) *T {
	return &v
}

func TestGameLibrary_JSON(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	game := GameLibrary{
		ID:            1,
		RawgID:        ptr(uint(12345)),
		Title:         "Elden Ring",
		Genre:         "RPG",
		Platform:      "PC",
		CoverURL:      "https://example.com/cover.jpg",
		PersonalNote:  ptr("Juegazo absoluto"),
		PersonalScore: ptr(10),
		Status:        ptr("Jugando"),
		AddedAt:       now,
	}

	data, err := json.Marshal(game)
	if err != nil {
		t.Fatalf("Error al codificar GameLibrary a JSON: %v", err)
	}

	jsonStr := string(data)
	expectedKeys := []string{"\"id\"", "\"rawg_id\"", "\"title\"", "\"personal_note\"", "\"personal_score\"", "\"added_at\""}
	for _, key := range expectedKeys {
		if !strings.Contains(jsonStr, key) {
			t.Errorf("El JSON serializado no contiene la propiedad esperada: %s", key)
		}
	}

	var decodedGame GameLibrary
	if err := json.Unmarshal(data, &decodedGame); err != nil {
		t.Fatalf("Error al decodificar JSON a GameLibrary: %v", err)
	}

	if decodedGame.Title != game.Title || *decodedGame.PersonalScore != *game.PersonalScore {
		t.Errorf("Los datos decodificados no coinciden con los originales")
	}
}

func TestGameStatsResponse_JSON(t *testing.T) {
	stats := GameStatsResponse{
		Total: 10,
		ByStatus: StatusCounts{
			Completado: 5,
			Jugando:    2,
			Pendiente:  2,
			Abandonado: 1,
		},
		AverageScore: 8.5,
	}

	data, err := json.Marshal(stats)
	if err != nil {
		t.Fatalf("Error al codificar GameStatsResponse: %v", err)
	}

	var decodedStats GameStatsResponse
	if err := json.Unmarshal(data, &decodedStats); err != nil {
		t.Fatalf("Error al decodificar GameStatsResponse: %v", err)
	}

	if decodedStats.ByStatus.Completado != 5 || decodedStats.AverageScore != 8.5 {
		t.Errorf("Fallo en la integridad de la serialización de estadísticas")
	}
}

func TestCustomError(t *testing.T) {
	msg := "Recurso no encontrado"

	t.Run("Factory Constructors", func(t *testing.T) {
		errNotFound := NewNotFoundError(msg)
		if errNotFound.Code != "404" || errNotFound.ErrorMsg != msg {
			t.Errorf("NewNotFoundError configuró parámetros erróneos")
		}

		errBadRequest := NewBadRequestError(msg)
		if errBadRequest.Code != "400" {
			t.Errorf("NewBadRequestError configuró un código erróneo")
		}

		errInternal := NewInternalServerError(msg)
		if errInternal.Code != "500" {
			t.Errorf("NewInternalServerError configuró un código erróneo")
		}

		errConflict := NewConflictError(msg)
		if errConflict.Code != "409" {
			t.Errorf("NewConflictError configuró un código erróneo")
		}
	})

	t.Run("Method Error()", func(t *testing.T) {
		customErr := NewError(http.StatusForbidden, "Acceso denegado")
		expectedFormat := "[403] Acceso denegado"
		if customErr.Error() != expectedFormat {
			t.Errorf("Se esperaba el formato '%s', se obtuvo '%s'", expectedFormat, customErr.Error())
		}
	})

	t.Run("JSON Ignore Tag", func(t *testing.T) {
		customErr := &CustomError{
			Code:     "500",
			ErrorMsg: "Crash",
			Err:      errors.New("error nativo interno de sistema"),
		}

		data, _ := json.Marshal(customErr)
		jsonStr := string(data)

		if !strings.Contains(jsonStr, `"error"`) || strings.Contains(jsonStr, "error nativo") {
			t.Errorf("El tag '-' no ignoró el campo interno 'Err' de forma correcta: %s", jsonStr)
		}
	})
}

func TestRAWGResponse_JSON(t *testing.T) {
	rawgData := RAWGResponse{
		Count:    1,
		Next:     ptr("https://api.rawg.io/api/games?page=2"),
		Previous: nil,
		Results: []Game{
			{
				ID:       99,
				Slug:     "hollow-knight",
				Name:     "Hollow Knight",
				Playtime: 30,
				Platforms: []PlatformElement{
					{Platform: PlatformDetail{ID: 1, Name: "Nintendo Switch", Slug: "nintendo-switch"}},
				},
				Stores: []StoreElement{
					{Store: StoreDetail{ID: 3, Name: "Steam", Slug: "steam"}},
				},
				BackgroundImage: ptr("https://image.rawg.io/bg.jpg"),
				Rating:          4.8,
				AddedByStatus: &AddedByStatus{
					Playing: 120,
					Beaten:  4500,
				},
				Genres: []Genre{
					{ID: 4, Name: "Action", Slug: "action"},
				},
			},
		},
	}

	data, err := json.Marshal(rawgData)
	if err != nil {
		t.Fatalf("Error al codificar RAWGResponse: %v", err)
	}

	var decodedRawg RAWGResponse
	if err := json.Unmarshal(data, &decodedRawg); err != nil {
		t.Fatalf("Error al decodificar RAWGResponse: %v", err)
	}

	if len(decodedRawg.Results) != 1 || decodedRawg.Results[0].Name != "Hollow Knight" {
		t.Errorf("Fallo en la deserialización de los sub-objetos internos de RAWGResponse")
	}
}

func TestGameDetail_JSON(t *testing.T) {
	detail := GameDetail{
		ID:          444,
		Slug:        "cyberpunk-2077",
		Name:        "Cyberpunk 2077",
		Description: "Welcome to Night City",
		Metacritic:  ptr(86),
		AddedByStatus: map[string]int{
			"yet":    340,
			"owned":  1200,
			"beaten": 800,
		},
		ParentPlatforms: []ParentPlatformItem{
			{Platform: ParentPlatformInfo{ID: 2, Name: "PlayStation", Slug: "playstation"}},
		},
		Platforms: []PlatformItem{
			{
				Platform: PlatformDetails{
					ID:        1,
					Name:      "PS5",
					Slug:      "ps5",
					YearStart: ptr(2020),
				},
				ReleasedAt: "2020-12-10",
			},
		},
		Developers: []Developer{
			{ID: 12, Name: "CD PROJEKT RED", Slug: "cd-projekt-red", GamesCount: 15},
		},
		DescriptionRaw: "Welcome to Night City Raw",
	}

	data, err := json.Marshal(detail)
	if err != nil {
		t.Fatalf("Error al serializar GameDetail: %v", err)
	}

	var decodedDetail GameDetail
	if err := json.Unmarshal(data, &decodedDetail); err != nil {
		t.Fatalf("Error al deserializar GameDetail: %v", err)
	}

	if decodedDetail.ID != 444 || decodedDetail.Developers[0].Name != "CD PROJEKT RED" {
		t.Errorf("Fallo en la validación de campos complejos en GameDetail")
	}

	if decodedDetail.AddedByStatus["owned"] != 1200 {
		t.Errorf("El mapa dinámico added_by_status perdió consistencia en la conversión")
	}
}
