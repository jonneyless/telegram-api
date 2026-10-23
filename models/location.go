package models

type Location struct {
	Latitude           float64 `json:"latitude"`
	Longitude          float64 `json:"longitude"`
	HorizontalAccuracy float64 `json:"horizontal_accuracy"`
}

type Venue struct {
	Location Location `json:"location"`
	Title    string   `json:"title"`
	Address  string   `json:"address"`
}
