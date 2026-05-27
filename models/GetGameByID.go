package models

type GameDetail struct {
	ID                        int                  `json:"id"`
	Slug                      string               `json:"slug"`
	Name                      string               `json:"name"`
	NameOriginal              string               `json:"name_original"`
	Description               string               `json:"description"`
	Metacritic                *int                 `json:"metacritic"`
	MetacriticPlatforms       []interface{}        `json:"metacritic_platforms"`
	Released                  string               `json:"released"`
	Tba                       bool                 `json:"tba"`
	Updated                   string               `json:"updated"`
	BackgroundImage           *string              `json:"background_image"`
	BackgroundImageAdditional *string              `json:"background_image_additional"`
	Website                   string               `json:"website"`
	Rating                    float64              `json:"rating"`
	RatingTop                 int                  `json:"rating_top"`
	Ratings                   []interface{}        `json:"ratings"`
	Reactions                 interface{}          `json:"reactions"`
	Added                     int                  `json:"added"`
	AddedByStatus             map[string]int       `json:"added_by_status"`
	Playtime                  int                  `json:"playtime"`
	ScreenshotsCount          int                  `json:"screenshots_count"`
	MoviesCount               int                  `json:"movies_count"`
	CreatorsCount             int                  `json:"creators_count"`
	AchievementsCount         int                  `json:"achievements_count"`
	ParentAchievementsCount   int                  `json:"parent_achievements_count"`
	RedditURL                 string               `json:"reddit_url"`
	RedditName                string               `json:"reddit_name"`
	RedditDescription         string               `json:"reddit_description"`
	RedditLogo                string               `json:"reddit_logo"`
	RedditCount               int                  `json:"reddit_count"`
	TwitchCount               int                  `json:"twitch_count"`
	YoutubeCount              int                  `json:"youtube_count"`
	ReviewsTextCount          int                  `json:"reviews_text_count"`
	RatingsCount              int                  `json:"ratings_count"`
	SuggestionsCount          int                  `json:"suggestions_count"`
	AlternativeNames          []string             `json:"alternative_names"`
	MetacriticURL             string               `json:"metacritic_url"`
	ParentsCount              int                  `json:"parents_count"`
	AdditionsCount            int                  `json:"additions_count"`
	GameSeriesCount           int                  `json:"game_series_count"`
	UserGame                  interface{}          `json:"user_game"`
	ReviewsCount              int                  `json:"reviews_count"`
	CommunityRating           int                  `json:"community_rating"`
	SaturatedColor            string               `json:"saturated_color"`
	DominantColor             string               `json:"dominant_color"`
	ParentPlatforms           []ParentPlatformItem `json:"parent_platforms"`
	Platforms                 []PlatformItem       `json:"platforms"`
	Stores                    []StoreItem          `json:"stores"`
	Developers                []Developer          `json:"developers"`
	Genres                    []Genres             `json:"genres"`
	Tags                      []Tags               `json:"tags"`
	Publishers                []interface{}        `json:"publishers"`
	EsrbRating                interface{}          `json:"esrb_rating"`
	Clip                      interface{}          `json:"clip"`
	DescriptionRaw            string               `json:"description_raw"`
}

type ParentPlatformItem struct {
	Platform ParentPlatformInfo `json:"platform"`
}

type ParentPlatformInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type PlatformItem struct {
	Platform     PlatformDetails `json:"platform"`
	ReleasedAt   string          `json:"released_at"`
	Requirements interface{}     `json:"requirements"`
}

type PlatformDetails struct {
	ID              int     `json:"id"`
	Name            string  `json:"name"`
	Slug            string  `json:"slug"`
	Image           *string `json:"image"`
	YearEnd         *int    `json:"year_end"`
	YearStart       *int    `json:"year_start"`
	GamesCount      int     `json:"games_count"`
	ImageBackground string  `json:"image_background"`
}

type StoreItem struct {
	ID    int       `json:"id"`
	URL   string    `json:"url"`
	Store StoreInfo `json:"store"`
}

type StoreInfo struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Slug            string `json:"slug"`
	Domain          string `json:"domain"`
	GamesCount      int    `json:"games_count"`
	ImageBackground string `json:"image_background"`
}

type Developer struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Slug            string `json:"slug"`
	GamesCount      int    `json:"games_count"`
	ImageBackground string `json:"image_background"`
}

type Genres struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Slug            string `json:"slug"`
	GamesCount      int    `json:"games_count"`
	ImageBackground string `json:"image_background"`
}

type Tags struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Slug            string `json:"slug"`
	Language        string `json:"language"`
	GamesCount      int    `json:"games_count"`
	ImageBackground string `json:"image_background"`
}
