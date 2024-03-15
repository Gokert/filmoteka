package models

type FilmItem struct {
	Id          uint64  `json:"id"`
	Title       string  `json:"title"`
	Info        string  `json:"info"`
	Rating      float64 `json:"rating"`
	ReleaseDate string  `json:"release_date"`
}
