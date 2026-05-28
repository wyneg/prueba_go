package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wyneg/prueba_go/models"
)

type MockGameRepository struct {
	mock.Mock
}

func (m *MockGameRepository) Create(ctx context.Context, game *models.GameLibrary) error {
	args := m.Called(ctx, game)
	return args.Error(0)
}

func (m *MockGameRepository) Select(ctx context.Context, status string) ([]models.GameLibrary, error) {
	args := m.Called(ctx, status)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]models.GameLibrary), args.Error(1)
}

func (m *MockGameRepository) Update(ctx context.Context, id uint, game *models.GameLibrary) error {
	args := m.Called(ctx, id, game)
	return args.Error(0)
}

func (m *MockGameRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockGameRepository) Stats(ctx context.Context) (models.GameStatsResponse, error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return models.GameStatsResponse{}, args.Error(1)
	}

	return args.Get(0).(models.GameStatsResponse), args.Error(1)
}

func TestCreateGame(t *testing.T) {

	var idValido uint = 123

	tests := []struct {
		name          string
		inputGame     *models.GameLibrary
		setupMock     func(m *MockGameRepository)
		expectedError string
	}{
		{
			name: "Error cuando RawgID es nil",
			inputGame: &models.GameLibrary{
				RawgID: nil,
			},
			setupMock: func(m *MockGameRepository) {
			},
			expectedError: "Game no se pudo crear, no envió RawgID",
		},
		{
			name: "Error cuando el repositorio falla",
			inputGame: &models.GameLibrary{
				RawgID: &idValido,
			},
			setupMock: func(m *MockGameRepository) {
				m.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			expectedError: "Game no se pudo crear: db error",
		},
		{
			name: "Caso exitoso",
			inputGame: &models.GameLibrary{
				RawgID: &idValido,
			},
			setupMock: func(m *MockGameRepository) {
				m.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockGameRepository)
			tt.setupMock(mockRepo)

			service := NewDBService(mockRepo)

			err := service.CreateGame(context.Background(), tt.inputGame)

			if tt.expectedError != "" {
				assert.NotNil(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
			} else {
				assert.Nil(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetGameFromDB(t *testing.T) {

	var rawgMock uint = 999
	idMock := 1
	tituloMock := "The Legend of Zelda"
	fechaMock := time.Now()

	juegosSimulados := []models.GameLibrary{
		{
			ID:      uint(idMock),
			RawgID:  &rawgMock,
			Title:   tituloMock,
			AddedAt: fechaMock,
		},
	}

	tests := []struct {
		name          string
		setupMock     func(m *MockGameRepository)
		expectedError string
		expectedData  interface{}
	}{
		{
			name: "Error cuando el repositorio falla",
			setupMock: func(m *MockGameRepository) {
				m.On("Select", mock.Anything, mock.Anything).Return(nil, errors.New("db connection failure"))
			},
			expectedError: "Error: db connection failure",
			expectedData:  nil,
		},
		{
			name: "Caso exitoso con datos o vacío",
			setupMock: func(m *MockGameRepository) {
				m.On("Select", mock.Anything, mock.Anything).Return(juegosSimulados, nil)
			},
			expectedError: "",
			expectedData:  juegosSimulados,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockRepo := new(MockGameRepository)
			tt.setupMock(mockRepo)

			service := NewDBService(mockRepo)

			result, err := service.GetGame(context.Background(), "")

			if tt.expectedError != "" {
				assert.NotNil(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				assert.Nil(t, result)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.expectedData, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateGame(t *testing.T) {
	var idInput uint = 42
	nota := "Un análisis de prueba"
	score := 8
	status := "jugando"

	gameInput := &models.GameLibrary{
		PersonalNote:  &nota,
		PersonalScore: &score,
		Status:        &status,
	}

	tests := []struct {
		name          string
		setupMock     func(m *MockGameRepository)
		expectedError string
	}{
		{
			name: "Error cuando el repositorio falla al actualizar",
			setupMock: func(m *MockGameRepository) {
				m.On("Update", mock.Anything, idInput, gameInput).
					Return(errors.New("db write failure"))
			},
			expectedError: "Error: db write failure",
		},
		{
			name: "Actualización exitosa",
			setupMock: func(m *MockGameRepository) {
				m.On("Update", mock.Anything, idInput, gameInput).
					Return(nil)
			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockRepo := new(MockGameRepository)
			tt.setupMock(mockRepo)

			service := NewDBService(mockRepo)

			err := service.UpdateGame(context.Background(), idInput, gameInput)

			if tt.expectedError != "" {
				assert.NotNil(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
			} else {
				assert.Nil(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteGame(t *testing.T) {
	var idInput uint = 99

	tests := []struct {
		name          string
		setupMock     func(m *MockGameRepository)
		expectedError string
	}{
		{
			name: "Error cuando el repositorio falla al eliminar",
			setupMock: func(m *MockGameRepository) {
				m.On("Delete", mock.Anything, idInput).
					Return(errors.New("db delete error"))
			},
			expectedError: "Error: db delete error",
		},
		{
			name: "Eliminación exitosa",
			setupMock: func(m *MockGameRepository) {
				m.On("Delete", mock.Anything, idInput).
					Return(nil)
			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockRepo := new(MockGameRepository)
			tt.setupMock(mockRepo)

			service := NewDBService(mockRepo)

			err := service.DeleteGame(context.Background(), idInput)

			if tt.expectedError != "" {
				assert.NotNil(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
			} else {
				assert.Nil(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestStatsGames(t *testing.T) {

	mockStatsData := models.GameStatsResponse{
		Total: 12,
		ByStatus: models.StatusCounts{
			Completado: 5,
			Jugando:    3,
			Pendiente:  4,
			Abandonado: 0,
		},
		AverageScore: 7.8,
	}

	tests := []struct {
		name          string
		setupMock     func(m *MockGameRepository)
		expectedError string
		expectedData  models.GameStatsResponse
	}{
		{
			name: "Error cuando el repositorio falla al obtener estadísticas",
			setupMock: func(m *MockGameRepository) {

				m.On("Stats", mock.Anything).
					Return(models.GameStatsResponse{}, errors.New("db query failure"))
			},
			expectedError: "Error: db query failure",
			expectedData:  models.GameStatsResponse{},
		},
		{
			name: "Estadísticas obtenidas con éxito",
			setupMock: func(m *MockGameRepository) {

				m.On("Stats", mock.Anything).
					Return(mockStatsData, nil)
			},
			expectedError: "",
			expectedData:  mockStatsData,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockRepo := new(MockGameRepository)
			tt.setupMock(mockRepo)

			service := NewDBService(mockRepo)

			result, err := service.StatsGames(context.Background())

			if tt.expectedError != "" {
				assert.NotNil(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				assert.Equal(t, tt.expectedData, result)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.expectedData, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
