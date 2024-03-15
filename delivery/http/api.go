package delivery

import (
	"encoding/json"
	"filmoteka/pkg/models"
	"filmoteka/usecase"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type Api struct {
	log  *logrus.Logger
	mx   *http.ServeMux
	core *usecase.Core
}

func GetApi(core *usecase.Core, log *logrus.Logger) *Api {
	api := &Api{
		core: core,
		log:  log,
		mx:   http.NewServeMux(),
	}

	api.mx.HandleFunc("/films", api.GetFilms)
	api.mx.HandleFunc("/films/find", api.FindFilms)

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
		page = 1
	}

	pageSize, err := strconv.ParseUint(r.URL.Query().Get("page_size"), 10, 64)
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
	}

	_, err = a.core.GetFilms(request)
	if err != nil {
		return
	}

	//a.log.Info(request)

}

func (a *Api) GetFilms(w http.ResponseWriter, r *http.Request) {
	response := models.Response{Status: http.StatusOK, Body: nil}
	a.log.Info(r.Host, r.URL)

	if r.Method != http.MethodGet {
		response.Status = http.StatusMethodNotAllowed
		a.SendResponse(w, r, &response)
		return
	}

	page, err := strconv.ParseUint(r.URL.Query().Get("page"), 10, 64)
	if err != nil {
		page = 1
	}

	pageSize, err := strconv.ParseUint(r.URL.Query().Get("page_size"), 10, 64)
	if err != nil {
		pageSize = 8
	}

	films, err := a.core.GetAll((page-1)*pageSize, pageSize)
	if err != nil {
		a.log.Error("get films error: ", err.Error())
		response.Status = http.StatusInternalServerError
		a.SendResponse(w, r, &response)
		return
	}

	response.Body = models.FilmsResponse{
		Page:     page,
		PageSize: pageSize,
		Total:    uint64(len(*films)),
		Films:    films,
	}

	a.SendResponse(w, r, &response)
}
