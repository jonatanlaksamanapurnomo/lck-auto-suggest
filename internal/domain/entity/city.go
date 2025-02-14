package entity

type City struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	ASCII      string  `json:"ascii"`
	AltNames   string  `json:"alt_names"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	Country    string  `json:"country"`
	Admin1     string  `json:"admin1"` // State/Province
	Population int64   `json:"population"`
	Timezone   string  `json:"timezone"`
}

type ScoredCity struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Score     float64 `json:"score"`
}
