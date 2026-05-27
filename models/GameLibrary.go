package models

type GameLibrary struct {
	ID            uint   `json:"id"`
	RawgID        *uint  `json:"rawg_id"`
	Title         string `json:"title"`
	Genre         string `json:"genre"`
	Platform      string `json:"platform"`
	CoverURL      string `json:"cover_url"`
	PersonalNote  string `json:"personal_note"`
	PersonalScore *int   `json:"personal_score"`
	Status        string `json:"status"`
	AddedAt       string `json:"added_at"`
}
