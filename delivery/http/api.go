package delivery

import (
	"encoding/json"
	"filmoteka/pkg/models"
	"filmoteka/usecase"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
)

type Api struct {
	log  *logrus.Logger
	mx   *http.ServeMux
	core usecase.ICore
}

func GetApi(core *usecase.Core, log *logrus.Logger) *Api {
	api := &Api{
		core: core,
		log:  log,
		mx:   http.NewServeMux(),
	}

	api.mx.HandleFunc("/api/v1/films", api.FindFilms)
	api.mx.HandleFunc("/signin", api.Signin)
	api.mx.HandleFunc("/authcheck", api.AuthAccept)
	//api.mx.HandleFunc("/films/find", api.FindFilms)

	return api
}

func (a *Api) ListenAndServe(port string) error {
	err := http.ListenAndServe(":"+port, a.mx)
	if err != nil {
		a.log.Error("ListenAndServer error: ", err.Error())
		return err
	}

	return nil
}

func (a *Api) SendResponse(w http.ResponseWriter, r *http.Request, response *models.Response) {
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		a.log.Error("Send response error: ", err)
		response.Status = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResponse)
	if err != nil {
		a.log.Error("Failed to send response: ", err.Error())
	}
}

func (a *Api) Signin(w http.ResponseWriter, r *http.Request) {
	response := models.Response{Status: http.StatusOK, Body: nil}
	a.log.Info(r.Host, r.URL)

	if r.Method != http.MethodPost {
		response.Status = http.StatusMethodNotAllowed
		a.SendResponse(w, r, &response)
		return
	}

	var request models.SigninRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.Status = http.StatusBadRequest
		a.SendResponse(w, r, &response)
		return
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		response.Status = http.StatusBadRequest
		a.SendResponse(w, r, &response)
		return
	}

	//_, err := a.core.

	a.SendResponse(w, r, &response)
}

func (a *Api) FindFilms(w http.ResponseWriter, r *http.Request) {
	response := models.Response{Status: http.StatusOK, Body: nil}
	a.log.Info(r.Host, r.URL)

	if r.Method != http.MethodGet {
		response.Status = http.StatusMethodNotAllowed
		a.SendResponse(w, r, &response)
		return
	}

	title := r.URL.Query().Get("title")
	actor := r.URL.Query().Get("actor")
	releaseDataFrom := r.URL.Query().Get("release_date_from")
	releaseDataTo := r.URL.Query().Get("release_date_to")
	order := r.URL.Query().Get("order")

	RatingFrom, err := strconv.ParseFloat(r.URL.Query().Get("rating_from"), 32)
	if err != nil {
		RatingFrom = 0
	}

	RatingTo, err := strconv.ParseFloat(r.URL.Query().Get("rating_to"), 32)
	if err != nil {
		RatingTo = 10
	}

	page, err := strconv.ParseUint(r.URL.Query().Get("page"), 10, 64)
	if err != nil {
		page = 0
	}

	pageSize, err := strconv.ParseUint(r.URL.Query().Get("per_page"), 10, 64)
	if err != nil {
		pageSize = 8
	}

	request := &models.FindFilmRequest{
		Title:           title,
		RatingFrom:      float32(RatingFrom),
		RatingTo:        float32(RatingTo),
		ReleaseDateFrom: releaseDataFrom,
		ReleaseDateTo:   releaseDataTo,
		Actor:           actor,
		Page:            page,
		PerPage:         pageSize,
		Order:           order,
	}

	films, err := a.core.GetFilms(request)
	if err != nil {
		return
	}

	response.Body = &models.FilmsResponse{
		Total: len(*films),
		Films: films,
	}

	a.SendResponse(w, r, &response)
}

func (a *Api) AuthAccept(w http.ResponseWriter, r *http.Request) {
	response := models.Response{Status: http.StatusOK, Body: nil}
	var authorized bool

	if r.Method != http.MethodGet {
		response.Status = http.StatusMethodNotAllowed
		a.SendResponse(w, r, &response)
		return
	}

	session, err := r.Cookie("session_id")
	if err == nil && session != nil {
		authorized, _ = a.core.FindActiveSession(r.Context(), session.Value)
	}

	if !authorized {
		response.Status = http.StatusUnauthorized
		a.SendResponse(w, r, &response)
		return
	}

	login, err := a.core.GetUserName(r.Context(), session.Value)
	if err != nil {
		a.log.Error("auth accept error: ", err.Error())
		response.Status = http.StatusInternalServerError
		a.SendResponse(w, r, &response)
		return
	}

	response.Body = models.AuthCheckResponse{
		Login: login,
		Role:  "hui",
	}

	a.SendResponse(w, r, &response)
}
