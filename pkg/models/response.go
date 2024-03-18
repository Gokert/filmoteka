package models

type Response struct {
	Status int `json:"status"`
	Body   any `json:"body"`
}

type FilmsResponse struct {
	Total int         `json:"total"`
	Films *[]FilmItem `json:"films"`
}

type SigninRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
type SignupRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthCheckResponse struct {
	Login string `json:"login"`
}

type FilmRequest struct {
	Id          uint64   `json:"id"`
	Title       string   `json:"title"`
	Info        string   `json:"info"`
	ReleaseDate string   `json:"release_date"`
	Rating      float32  `json:"rating"`
	Actors      []uint64 `json:"actors"`
}

type ActorResponse struct {
	Id       uint64     `json:"id"`
	Name     string     `json:"name"`
	Gender   string     `json:"gen"`
	Birthday string     `json:"birthday"`
	Films    []FilmItem `json:"films"`
}

type ActorRequest struct {
	Id       uint64   `json:"id"`
	Name     string   `json:"name"`
	Gender   string   `json:"gen"`
	Birthday string   `json:"birthday"`
	Films    []uint64 `json:"films"`
}

type FindFilmRequest struct {
	Title           string  `json:"title"`
	RatingFrom      float32 `json:"rating_from"`
	RatingTo        float32 `json:"rating_to"`
	ReleaseDateFrom string  `json:"release_date_from"`
	ReleaseDateTo   string  `json:"release_date_to"`
	Actor           string  `json:"actor"`
	Page            uint64  `json:"page"`
	PerPage         uint64  `json:"per_page"`
	Order           string  `json:"order"`
}
