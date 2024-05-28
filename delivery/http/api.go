package delivery

import (
	"encoding/json"
	_ "filmoteka/docs"
	"filmoteka/pkg/middleware"
	"filmoteka/pkg/models"
	httpResponse "filmoteka/pkg/response"
	"filmoteka/usecase"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
	"time"
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

	api.mx.HandleFunc("/signin", api.Signin)
	api.mx.HandleFunc("/signup", api.Signup)
	api.mx.HandleFunc("/logout", api.Logout)
	api.mx.HandleFunc("/authcheck", api.AuthAccept)

	api.mx.HandleFunc("/api/v1/actors", api.FindActors)
	api.mx.Handle("/api/v1/actors/add", middleware.AuthCheck(middleware.CheckRole(http.HandlerFunc(api.AddActor), core, log), core, log))
	api.mx.Handle("/api/v1/actors/update", middleware.AuthCheck(middleware.CheckRole(http.HandlerFunc(api.UpdateActor), core, log), core, log))
	api.mx.Handle("/api/v1/actors/delete", middleware.AuthCheck(middleware.CheckRole(http.HandlerFunc(api.DeleteActor), core, log), core, log))

	api.mx.HandleFunc("/api/v1/films", api.FindFilms)
	api.mx.HandleFunc("/api/v1/films/search", api.SearchFilms)
	api.mx.Handle("/api/v1/films/add", middleware.AuthCheck(middleware.CheckRole(http.HandlerFunc(api.AddFilm), core, log), core, log))
	api.mx.Handle("/api/v1/films/update", middleware.AuthCheck(middleware.CheckRole(http.HandlerFunc(api.UpdateFilm), core, log), core, log))
	api.mx.Handle("/api/v1/films/delete", middleware.AuthCheck(middleware.CheckRole(http.HandlerFunc(api.DeleteFilm), core, log), core, log))

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

// @Summary signIn
// @Tags Auth
// @Description authenticate user by providing login and password credentials
// @ID authenticate-user
// @Accept json
// @Produce json
// @Param input body models.SigninRequest false "login and password"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 405 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /signin [post]
func (a *Api) Signin(w http.ResponseWriter, r *http.Request) {
	response := models.Response{Status: http.StatusOK, Body: nil}

	if r.Method != http.MethodPost {
		response.Status = http.StatusMethodNotAllowed
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	var request models.SigninRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		a.log.Error("Signin error: ", err.Error())
		response.Status = http.StatusBadRequest
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		a.log.Error("Signin error: ", err.Error())
		response.Status = http.StatusBadRequest
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	_, found, err := a.core.FindUserAccount(r.Context(), request.Login, request.Password)
	if err != nil {
		a.log.Error("Signin error: ", err.Error())
		response.Status = http.StatusInternalServerError
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	if !found {
		response.Status = http.StatusUnauthorized
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	session, _ := a.core.CreateSession(r.Context(), request.Login)
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    session.SID,
		Path:     "/",
		Expires:  session.ExpiresAt,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	httpResponse.SendResponse(w, r, &response, a.log)
}

// @Summary signUp
// @Tags Auth
// @Desription create account
// @ID create-account
// @Accept json
// @Produce json
// @Param input body models.SignupRequest false "account information"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 405 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /signup [post]
func (a *Api) Signup(w http.ResponseWriter, r *http.Request) {
	response := models.Response{Status: http.StatusOK, Body: nil}

	if r.Method != http.MethodPost {
		response.Status = http.StatusMethodNotAllowed
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	var request models.SignupRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.Status = http.StatusBadRequest
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		response.Status = http.StatusInternalServerError
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	found, err := a.core.FindUserByLogin(r.Context(), request.Login)
	if err != nil {
		a.log.Error("Signup error: ", err.Error())
		response.Status = http.StatusInternalServerError
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	if found {
		response.Status = http.StatusConflict
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	err = a.core.CreateUserAccount(r.Context(), request.Login, request.Password)
	if err != nil {
		a.log.Error("create user error: ", err.Error())
		response.Status = http.StatusBadRequest
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	httpResponse.SendResponse(w, r, &response, a.log)
}

// @Summary add a new film
// @Description add a new film along with associated actors
// @Tags Film
// @Accept json
// @Produce json
// @Param session_id header string false "Session ID"
// @Param input body models.FilmRequest true "Film details and actors"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 405 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/films/add [post]
func (a *Api) AddFilm(w http.ResponseWriter, r *http.Request) {
	response := models.Response{Status: http.StatusOK, Body: nil}

	if r.Method != http.MethodPost {
		response.Status = http.StatusMethodNotAllowed
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	var request models.FilmRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.Status = http.StatusBadRequest
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		response.Status = http.StatusInternalServerError
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	_, err = a.core.AddFilm(r.Context(), &request, request.Actors)
	if err != nil {
		response.Status = http.StatusInternalServerError
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	httpResponse.SendResponse(w, r, &response, a.log)
}

// @Summary add a new actor
// @Description add a new actor
// @Tags Actor
// @Accept json
// @Produce json
// @Param session_id header string false "Session ID"
// @Param input body models.ActorItem true "Actor details"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 405 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/actors/add [post]
func (a *Api) AddActor(w http.ResponseWriter, r *http.Request) {
	response := models.Response{Status: http.StatusOK, Body: nil}

	if r.Method != http.MethodPost {
		response.Status = http.StatusMethodNotAllowed
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	var request models.ActorItem

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.Status = http.StatusBadRequest
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		response.Status = http.StatusInternalServerError
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	_, err = a.core.AddActor(r.Context(), &request)
	if err != nil {
		response.Status = http.StatusInternalServerError
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	httpResponse.SendResponse(w, r, &response, a.log)
}

// @Summary search for films by title and actor name
// @Description search for films by title and actor name, optionally specify page number and size
// @Tags Film
// @Accept json
// @Produce json
// @Param title_film query string false "Movie title fragment"
// @Param name_actor query string false "Actor name fragment"
// @Param page query uint64 false "Page number (optional)" Enums(0)
// @Param per_page query uint64 false "Number of results per page (optional)" Enums(8)
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 405 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/films/search [get]
func (a *Api) SearchFilms(w http.ResponseWriter, r *http.Request) {
	response := models.Response{Status: http.StatusOK, Body: nil}

	if r.Method != http.MethodGet {
		response.Status = http.StatusMethodNotAllowed
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	titleFilm := r.URL.Query().Get("title_film")
	nameActor := r.URL.Query().Get("name_actor")

	page, err := strconv.ParseUint(r.URL.Query().Get("page"), 10, 64)
	if err != nil {
		page = 0
	}

	pageSize, err := strconv.ParseUint(r.URL.Query().Get("per_page"), 10, 64)
	if err != nil {
		pageSize = 8
	}

	films, err := a.core.SearchFilms(r.Context(), titleFilm, nameActor, page, pageSize)
	if err != nil {
		response.Status = http.StatusInternalServerError
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	response.Body = films

	httpResponse.SendResponse(w, r, &response, a.log)
}

// @Summary find films based on various criteria
// @Description get a list of films based on title, actor, release date, rating, and order
// @Tags Film
// @Accept json
// @Produce json
// @Param title query string false "Film title" example:"The Shawshank Redemption"
// @Param actor query string false "Actor name" example:"Tim Robbins"
// @Param release_date_from query string false "Release date from" format="date" example:"1994-01-01"
// @Param release_date_to query string false "Release date to" format="date" example:"1995-12-31"
// @Param rating_from query number false "Minimum rating" example:"7.0" minimum="0" maximum="10"
// @Param rating_to query number false "Maximum rating" example:"8.5" minimum="0" maximum="10"
// @Param order query string false "Sorting order" enum:"asc,desc" default:"desc"
// @Param page query integer false "Page number" example:"1" minimum="1"
// @Param per_page query integer false "Number of items per page" example:"20" minimum="1" maximum="100"
// @Success 200 {object} models.FilmsResponse "Successful response"
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 405 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/films [get]
func (a *Api) FindFilms(w http.ResponseWriter, r *http.Request) {
	response := models.Response{Status: http.StatusOK, Body: nil}

	if r.Method != http.MethodGet {
		response.Status = http.StatusMethodNotAllowed
		httpResponse.SendResponse(w, r, &response, a.log)
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

	films, err := a.core.GetFilms(r.Context(), request)
	if err != nil {
		return
	}

	response.Body = &models.FilmsResponse{
		Total: len(*films),
		Films: films,
	}

	httpResponse.SendResponse(w, r, &response, a.log)
}

// @Summary delete a film by ID
// @Description deletes a film with the given ID
// @Tags Film
// @Accept json
// @Produce json
// @Param film_id query integer true "Film ID"
// @Param session_id header string false "Session ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 405 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/films/delete [delete]
func (a *Api) DeleteFilm(w http.ResponseWriter, r *http.Request) {
	response := models.Response{Status: http.StatusOK, Body: nil}

	if r.Method != http.MethodDelete {
		response.Status = http.StatusMethodNotAllowed
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	filmId, err := strconv.ParseUint(r.URL.Query().Get("film_id"), 10, 64)
	if err != nil {
		response.Status = http.StatusBadRequest
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	_, err = a.core.DeleteFilm(r.Context(), filmId)
	if err != nil {
		response.Status = http.StatusInternalServerError
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	httpResponse.SendResponse(w, r, &response, a.log)
}

// @Summary update film information
// @Tags Film
// @ID update-film
// @Produce json
// @Consume json
// @Param session_id header string false "Session ID"
// @Param Film body models.FilmRequest true "Updated Film Information"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 405 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/films/update [patch]
func (a *Api) UpdateFilm(w http.ResponseWriter, r *http.Request) {
	response := models.Response{Status: http.StatusOK, Body: nil}

	if r.Method != http.MethodPatch {
		response.Status = http.StatusMethodNotAllowed
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	var request models.FilmRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.Status = http.StatusBadRequest
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		response.Status = http.StatusInternalServerError
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	err = a.core.UpdateFilm(r.Context(), &request)
	if err != nil {
		response.Status = http.StatusInternalServerError
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	httpResponse.SendResponse(w, r, &response, a.log)
}

// @Summary get list of actors with pagination
// @Tags Actor
// @ID find-actors
// @Produce json
// @Param page query uint64 false "Page number, starting from 0 (optional)"
// @Param per_size query uint64 false "Number of items per page, defaults to 8 (optional)"
// @Success 200 {array} models.ActorItem
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 405 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/actors [get]
func (a *Api) FindActors(w http.ResponseWriter, r *http.Request) {
	response := models.Response{Status: http.StatusOK, Body: nil}

	if r.Method != http.MethodGet {
		response.Status = http.StatusMethodNotAllowed
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	page, err := strconv.ParseUint(r.URL.Query().Get("page"), 10, 64)
	if err != nil {
		page = 0
	}

	perSize, err := strconv.ParseUint(r.URL.Query().Get("per_size"), 10, 64)
	if err != nil {
		perSize = 8
	}

	actors, err := a.core.FindActors(r.Context(), page, perSize)
	if err != nil {
		response.Status = http.StatusInternalServerError
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	response.Body = actors

	httpResponse.SendResponse(w, r, &response, a.log)
}

// @Summary delete actor by ID
// @Tags Actor
// @ID delete-actor
// @Produce json
// @Param actor_id query uint64 true "Actor ID"
// @Param session_id header string false "Session ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 405 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/actors/delete [delete]
func (a *Api) DeleteActor(w http.ResponseWriter, r *http.Request) {
	response := models.Response{Status: http.StatusOK, Body: nil}

	if r.Method != http.MethodDelete {
		response.Status = http.StatusMethodNotAllowed
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	actorId, err := strconv.ParseUint(r.URL.Query().Get("actor_id"), 10, 64)
	if err != nil {
		response.Status = http.StatusBadRequest
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	err = a.core.DeleteActor(r.Context(), actorId)
	if err != nil {
		response.Status = http.StatusInternalServerError
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	httpResponse.SendResponse(w, r, &response, a.log)
}

// @Summary update actor information
// @Tags Actor
// @ID update-actor
// @Produce json
// @Consume json
// @Param session_id header string false "Session ID"
// @Param Actor body models.ActorRequest true "Updated Actor Information"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 405 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/actors/update [patch]
func (a *Api) UpdateActor(w http.ResponseWriter, r *http.Request) {
	response := models.Response{Status: http.StatusOK, Body: nil}

	if r.Method != http.MethodPatch {
		response.Status = http.StatusMethodNotAllowed
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	var request models.ActorRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.Status = http.StatusBadRequest
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		response.Status = http.StatusInternalServerError
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	err = a.core.UpdateActor(r.Context(), &request)
	if err != nil {
		response.Status = http.StatusInternalServerError
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	httpResponse.SendResponse(w, r, &response, a.log)
}

// @Summary end current user session
// @Tags Auth
// @ID logout
// @Produce json
// @Param session_id header string false "Session ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 405 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /logout [delete]
func (a *Api) Logout(w http.ResponseWriter, r *http.Request) {
	response := models.Response{Status: http.StatusOK, Body: nil}

	if r.Method != http.MethodDelete {
		response.Status = http.StatusMethodNotAllowed
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		response.Status = http.StatusBadRequest
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	err = a.core.KillSession(r.Context(), cookie.Value)
	if err != nil {
		response.Status = http.StatusInternalServerError
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	cookie.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, cookie)

	httpResponse.SendResponse(w, r, &response, a.log)
}

// @summary check authentication status and return user info
// @description returns user info if they are currently logged in
// @Tags Auth
// @produce application/json
// @Param session_id header string false "Session ID"
// @success 200 {object} models.AuthCheckResponse
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 405 {object} models.Response
// @Failure 500 {object} models.Response
// @router /authcheck [get]
func (a *Api) AuthAccept(w http.ResponseWriter, r *http.Request) {
	response := models.Response{Status: http.StatusOK, Body: nil}
	var authorized bool

	if r.Method != http.MethodGet {
		response.Status = http.StatusMethodNotAllowed
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	session, err := r.Cookie("session_id")
	if err == nil && session != nil {
		authorized, _ = a.core.FindActiveSession(r.Context(), session.Value)
	}
	a.log.Warning("API", authorized)
	if !authorized {
		response.Status = http.StatusUnauthorized
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	login, err := a.core.GetUserName(r.Context(), session.Value)
	if err != nil {
		a.log.Error("auth accept error: ", err.Error())
		response.Status = http.StatusInternalServerError
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	response.Body = models.AuthCheckResponse{
		Login: login,
	}

	httpResponse.SendResponse(w, r, &response, a.log)
}
