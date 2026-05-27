package models

type RAWGResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []Game  `json:"results"`
}

type Game struct {
	ID               int               `json:"id"`
	Slug             string            `json:"slug"`
	Name             string            `json:"name"`
	Playtime         int               `json:"playtime"`
	Platforms        []PlatformElement `json:"platforms"`
	Stores           []StoreElement    `json:"stores"`
	Released         string            `json:"released"`
	Tba              bool              `json:"tba"`
	BackgroundImage  *string           `json:"background_image"`
	Rating           float64           `json:"rating"`
	RatingTop        int               `json:"rating_top"`
	Ratings          []Rating          `json:"ratings"`
	RatingsCount     int               `json:"ratings_count"`
	ReviewsTextCount int               `json:"reviews_text_count"`
	Added            int               `json:"added"`
	AddedByStatus    *AddedByStatus    `json:"added_by_status"`
	Metacritic       *int              `json:"metacritic"`
	SuggestionsCount int               `json:"suggestions_count"`
	Updated          string            `json:"updated"`
	Score            string            `json:"score"`
	Clip             *string           `json:"clip"`
	Tags             []Tag             `json:"tags"`
	EsrbRating       *EsrbRating       `json:"esrb_rating"`
	UserGame         interface{}       `json:"user_game"`
	ReviewsCount     int               `json:"reviews_count"`
	CommunityRating  int               `json:"community_rating,omitempty"`
	SaturatedColor   string            `json:"saturated_color"`
	DominantColor    string            `json:"dominant_color"`
	ShortScreenshots []ShortScreenshot `json:"short_screenshots"`
	ParentPlatforms  []ParentPlatform  `json:"parent_platforms"`
	Genres           []Genre           `json:"genres"`
}

type PlatformElement struct {
	Platform PlatformDetail `json:"platform"`
}

type PlatformDetail struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type StoreElement struct {
	Store StoreDetail `json:"store"`
}

type StoreDetail struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type Rating struct {
	ID      int     `json:"id"`
	Title   string  `json:"title"`
	Count   int     `json:"count"`
	Percent float64 `json:"percent"`
}

type AddedByStatus struct {
	Yet     int `json:"yet"`
	Owned   int `json:"owned"`
	Beaten  int `json:"beaten"`
	Dropped int `json:"dropped"`
	Playing int `json:"playing"`
	Toplay  int `json:"toplay"`
}

type Tag struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Slug            string `json:"slug"`
	Language        string `json:"language"`
	GamesCount      int    `json:"games_count"`
	ImageBackground string `json:"image_background"`
}

type EsrbRating struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Slug   string `json:"slug"`
	NameEn string `json:"name_en"`
	NameRu string `json:"name_ru"`
}

type ShortScreenshot struct {
	ID    int    `json:"id"`
	Image string `json:"image"`
}

type ParentPlatform struct {
	Platform PlatformDetail `json:"platform"`
}

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}
