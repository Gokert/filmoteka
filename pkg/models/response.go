package models

type Response struct {
	Status int `json:"status"`
	Body   any `json:"body"`
}

type FilmsResponse struct {
	Page     uint64      `json:"current_page"`
	PageSize uint64      `json:"page_size"`
	Total    uint64      `json:"total"`
	Films    *[]FilmItem `json:"films"`
}
