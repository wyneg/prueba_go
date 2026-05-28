package models

type StatusCounts struct {
	Completado uint `json:"completado"`
	Jugando    uint `json:"jugando"`
	Pendiente  uint `json:"pendiente"`
	Abandonado uint `json:"abandonado"`
}

type GameStatsResponse struct {
	Total        uint         `json:"total"`
	ByStatus     StatusCounts `json:"by_status"`
	AverageScore float64      `json:"average_score"`
}
