package psx

import (
	"context"
	"database/sql"
	"errors"
	"filmoteka/configs"
	utils "filmoteka/pkg"
	"filmoteka/pkg/models"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

type PsxRepo struct {
	db *sql.DB
}

func GetFilmRepo(config *configs.DbPsxConfig, log *logrus.Logger) (*PsxRepo, error) {
	dsn := fmt.Sprintf("user=%s dbname=%s password= %s host=%s port=%d sslmode=%s",
		config.User, config.Dbname, config.Password, config.Host, config.Port, config.Sslmode)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Errorf("sql open error: %s", err.Error())
		return nil, fmt.Errorf("get user repo err: %s", err.Error())
	}

	repo := &PsxRepo{db: db}

	errs := make(chan error)
	go func() {
		errs <- repo.pingDb(3, log)
	}()

	if err := <-errs; err != nil {
		log.Error(err.Error())
		return nil, err
	}
	db.SetMaxOpenConns(config.MaxOpenConns)

	log.Info("Postgres created successful on ", config.Port)
	return repo, nil
}
func (repo *PsxRepo) pingDb(timer uint32, log *logrus.Logger) error {
	var err error
	var retries int

	for retries < utils.MaxRetries {
		err = repo.db.Ping()
		if err == nil {
			return nil
		}

		retries++
		log.Errorf("sql ping error: %s", err.Error())
		time.Sleep(time.Duration(timer) * time.Second)
	}

	return fmt.Errorf("sql max pinging error: %s", err.Error())
}

func (repo *PsxRepo) GetFilms(ctx context.Context, request *models.FindFilmRequest) (*[]models.FilmItem, error) {
	films := make([]models.FilmItem, 0, request.PerPage)
	var s strings.Builder
	var hasWhere bool
	paramNum := 1
	var params []interface{}

	s.WriteString("SELECT film.id ,film.title, film.rating, film.release_date, film.info FROM film " +
		"LEFT JOIN actor_in_film ON actor_in_film.id_film = film.id " +
		"LEFT JOIN actor ON actor_in_film.id_actor = actor.id ")

	if request.Title != "" {
		s.WriteString("WHERE ")
		hasWhere = true
		s.WriteString("fts @@ to_tsquery($" + strconv.Itoa(paramNum) + ") ")
		paramNum++
		params = append(params, request.Title)
	}
	if request.ReleaseDateFrom != "" {
		if !hasWhere {
			s.WriteString("WHERE ")
			hasWhere = true
		} else {
			s.WriteString("AND ")
		}
		s.WriteString("release_date >= $" + strconv.Itoa(paramNum) + " ")
		paramNum++
		params = append(params, request.ReleaseDateFrom)
	}
	if request.ReleaseDateTo != "" {
		if !hasWhere {
			s.WriteString("WHERE ")
			hasWhere = true
		} else {
			s.WriteString("AND ")
		}
		s.WriteString("release_date <= $" + strconv.Itoa(paramNum) + " ")
		paramNum++
	}

	if !hasWhere {
		s.WriteString("WHERE ")
	} else {
		s.WriteString("AND ")
	}
	s.WriteString("rating >= $" + strconv.Itoa(paramNum) + " AND " + "rating <= $" + strconv.Itoa(paramNum+1) + " ")
	paramNum += 2
	params = append(params, request.RatingFrom, request.RatingTo)

	s.WriteString("GROUP BY film.rating, film.id, film.title, film.release_date, film.info ")

	switch request.Order {
	case "title":
		s.WriteString("ORDER BY film.title DESC ")
	case "release_date":
		s.WriteString("ORDER BY film.release_date DESC ")
	case "rating":
		s.WriteString("ORDER BY film.rating DESC ")
	default:
		s.WriteString("ORDER BY film.rating DESC ")
	}

	s.WriteString("OFFSET $" + strconv.Itoa(paramNum) + " LIMIT $" + strconv.Itoa(paramNum+1))
	params = append(params, request.Page, request.PerPage)

	rows, err := repo.db.QueryContext(ctx, s.String(), params...)
	if err != nil {
		return nil, fmt.Errorf("find film err: %s", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		post := models.FilmItem{}

		err := rows.Scan(&post.Id, &post.Title, &post.Rating, &post.Info, &post.ReleaseDate)
		if err != nil {
			return nil, fmt.Errorf("find film scan err: %s", err.Error())
		}

		films = append(films, post)
	}

	return &films, nil
}

func (repo *PsxRepo) SearchFilms(ctx context.Context, titleFilm string, nameActor string, page uint64, perPage uint64) ([]models.FilmItem, error) {
	films := make([]models.FilmItem, 0, perPage)
	var s strings.Builder
	haveSelect := false
	var params []interface{}
	count := 0

	s.WriteString("SELECT film.id ,film.title, film.info, film.rating, film.release_date FROM film " +
		"LEFT JOIN actor_in_film ON actor_in_film.id_film = film.id " +
		"LEFT JOIN actor ON actor_in_film.id_actor = actor.id ")

	if titleFilm != "" {
		haveSelect = true
		params = append(params, titleFilm)
		count++
		s.WriteString("WHERE film.title LIKE '%' || $" + strconv.Itoa(count) + " || '%'")
	}

	if nameActor != "" {
		params = append(params, nameActor)
		count++
		if haveSelect {
			s.WriteString("AND actor.name LIKE '%' || $" + strconv.Itoa(count) + " || '%'")
		} else {
			s.WriteString("WHERE actor.name LIKE '%' || $" + strconv.Itoa(count) + " || '%'")
		}
	}

	s.WriteString("ORDER BY film.rating DESC ")
	s.WriteString("OFFSET $" + strconv.Itoa(count+1) + " LIMIT $" + strconv.Itoa(count+2) + " ")
	params = append(params, page, perPage)

	rows, err := repo.db.QueryContext(ctx, s.String(), params...)
	if err != nil {
		return nil, fmt.Errorf("find film error: %s", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		post := models.FilmItem{}
		err := rows.Scan(&post.Id, &post.Title, &post.Info, &post.Rating, &post.ReleaseDate)
		if err != nil {
			return nil, fmt.Errorf("find film scan err: %s", err.Error())
		}
		films = append(films, post)
	}

	return films, nil
}

func (repo *PsxRepo) FindActors(ctx context.Context, page uint64, perPage uint64) ([]models.ActorResponse, error) {
	rows, err := repo.db.QueryContext(ctx, `
		SELECT
			actor.id,
			actor.name,
			actor.gen,
			actor.birthdate,
			film.id,
			film.title,
			film.info,
			film.release_date,
			film.rating
		FROM (
				 SELECT * FROM actor
				 OFFSET $1 LIMIT $2
			 ) AS actor
				 LEFT JOIN actor_in_film ON actor.id = actor_in_film.id_actor
				 LEFT JOIN film ON actor_in_film.id_film = film.id`, page, perPage)
	if err != nil {
		return nil, fmt.Errorf("sql query error: %s", err.Error())
	}
	defer rows.Close()

	actorsMap := make(map[uint64]*models.ActorResponse)
	for rows.Next() {
		var actorID uint64
		var actorName, actorGender, actorBirthday string
		var filmID sql.NullInt64
		var filmTitle, filmInfo, filmReleaseDate string
		var filmRating float64

		err := rows.Scan(&actorID, &actorName, &actorGender, &actorBirthday, &filmID, &filmTitle, &filmInfo, &filmReleaseDate, &filmRating)
		if err != nil && filmID.Valid {
			return nil, fmt.Errorf("sql Scan error: %s", err.Error())
		}
		actor, ok := actorsMap[actorID]
		if !ok {
			actor = &models.ActorResponse{
				Id:       actorID,
				Name:     actorName,
				Gender:   actorGender,
				Birthday: actorBirthday,
			}
			actorsMap[actorID] = actor
		}

		if filmID.Valid {
			actor.Films = append(actor.Films, models.FilmItem{
				Id:          uint64(filmID.Int64),
				Title:       filmTitle,
				Info:        filmInfo,
				ReleaseDate: filmReleaseDate,
				Rating:      filmRating,
			})
		}
	}

	var actors []models.ActorResponse
	for _, actor := range actorsMap {
		actors = append(actors, *actor)
	}

	return actors, nil
}

func (repo *PsxRepo) FindFilmsByActor(ctx context.Context, actorId uint64) ([]models.FilmItem, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT film.id, film.title, film,info, film.release_date FROM film LEFT JOIN actor_in_film ON actor_in_film.id_film = film.id LEFT JOIN actor ON actor_in_film.id_actor = actor.id WHERE actor.id = $1", actorId)
	if err != nil {
		return nil, fmt.Errorf("sql request error: %s", err.Error())
	}

	var response []models.FilmItem

	for rows.Next() {
		var film models.FilmItem

		err := rows.Scan(&film.Id, &film.Title, &film.Info, &film.ReleaseDate)
		if err != nil {
			return nil, fmt.Errorf("sql Scan error: %s", err.Error())
		}
		response = append(response, film)
	}

	return response, nil
}

func (repo *PsxRepo) GetRelationByFilmId(ctx context.Context, filmId uint64) ([]uint64, error) {
	var ids []uint64

	rows, err := repo.db.QueryContext(ctx, `SELECT actor_in_film.id_actor FROM actor_in_film WHERE actor_in_film.id_film=$1`, filmId)
	if err != nil {
		return nil, fmt.Errorf("sql request find relation actors error: %s", err.Error())
	}

	for rows.Next() {
		var id uint64

		err := rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("sql scan actors in films error: %s", err.Error())
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (repo *PsxRepo) GetRelationByActorId(ctx context.Context, actorId uint64) ([]uint64, error) {
	var ids []uint64

	rows, err := repo.db.QueryContext(ctx, `SELECT actor_in_film.id_film FROM actor_in_film WHERE actor_in_film.id_actor=$1`, actorId)
	if err != nil {
		return nil, fmt.Errorf("sql request find relation films error: %s", err.Error())
	}

	for rows.Next() {
		var id uint64

		err := rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("sql scan actors in films error: %s", err.Error())
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (repo *PsxRepo) DeleteRelation(ctx context.Context, filmId uint64, actorId uint64) error {
	_, err := repo.db.QueryContext(ctx, `DELETE FROM actor_in_film WHERE id_actor=$1 AND id_film=$2`, actorId, filmId)
	if err != nil {
		return fmt.Errorf("sql delete relation error: %s", err.Error())
	}

	return err
}

func (repo *PsxRepo) InsertRelation(ctx context.Context, filmId uint64, actorId uint64) error {
	result, err := repo.db.ExecContext(ctx, `INSERT INTO actor_in_film (id_actor, id_film) VALUES ($1, $2)`, actorId, filmId)
	if err != nil && result != nil {
		return fmt.Errorf("sql insert relation error: %s", err.Error())
	}

	return err
}

func (repo *PsxRepo) UpdateFilm(ctx context.Context, film *models.FilmRequest) error {
	if film.Id == 0 {
		return fmt.Errorf("film id missing")
	}

	var s strings.Builder
	var haveEq bool = false
	params := make([]any, 0, 5)
	count := 0

	s.WriteString("UPDATE film SET ")
	if film.Title != "" {
		if haveEq {
			s.WriteString(",")
		} else {
			haveEq = true
		}

		params = append(params, film.Title)
		count++
		s.WriteString("title = $" + strconv.Itoa(count))
	}
	if film.Info != "" {
		if haveEq {
			s.WriteString(",")
		} else {
			haveEq = true
		}

		params = append(params, film.Info)
		count++
		s.WriteString(" info = $" + strconv.Itoa(count))
	}
	if film.ReleaseDate != "" {
		if haveEq {
			s.WriteString(",")
		} else {
			haveEq = true
		}
		params = append(params, film.ReleaseDate)
		count++
		s.WriteString(" release_date = $" + strconv.Itoa(count))
	}
	if film.Rating != 0 {
		if haveEq {
			s.WriteString(",")
		} else {
			haveEq = true
		}

		params = append(params, film.Rating)
		count++
		s.WriteString(" rating = $" + strconv.Itoa(count))
	}

	params = append(params, film.Id)
	count++
	s.WriteString(" WHERE film.id = $" + strconv.Itoa(count))

	if count < 2 {
		return fmt.Errorf("not have params")
	}

	_, err := repo.db.QueryContext(ctx, s.String(), params...)
	if err != nil {
		return fmt.Errorf("update film error: %s", err.Error())
	}

	for range film.Actors {
		existingActorIds, err := repo.GetRelationByFilmId(ctx, film.Id)
		if err != nil {
			return err
		}

		for _, existingActorId := range existingActorIds {
			found := false
			for _, actorId := range film.Actors {
				if existingActorId == actorId {
					found = true
					break
				}
			}
			if !found {
				err := repo.DeleteRelation(ctx, existingActorId, film.Id)
				if err != nil {
					return err
				}
			}
		}

		for _, actorId := range film.Actors {
			found := false
			for _, existingActorId := range existingActorIds {
				if existingActorId == actorId {
					found = true
					break
				}
			}
			if !found {
				err := repo.InsertRelation(ctx, film.Id, actorId)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (repo *PsxRepo) DeleteFilm(ctx context.Context, filmId uint64) (bool, error) {
	_, err := repo.db.ExecContext(ctx, "DELETE FROM film "+
		"WHERE film.id = $1", filmId)
	if err != nil {
		return false, fmt.Errorf("remove favorite film error: %s", err.Error())
	}

	return true, nil
}

func (repo *PsxRepo) DeleteActor(ctx context.Context, actorId uint64) error {
	_, err := repo.db.ExecContext(ctx, "DELETE FROM actor WHERE actor.id = $1", actorId)
	if err != nil {
		return fmt.Errorf("sql exec error: %s", err.Error())
	}

	return nil
}

func (repo *PsxRepo) AddFilm(ctx context.Context, film *models.FilmRequest) (uint64, error) {
	err := repo.db.QueryRowContext(ctx, "INSERT INTO film(title, info, release_date, rating) VALUES($1, $2, $3, $4) RETURNING id",
		film.Title, film.Info, film.ReleaseDate, film.Rating).Scan(&film.Id)
	if err != nil {
		return 0, fmt.Errorf("insert film err: %s", err.Error())
	}

	return film.Id, nil
}

func (repo *PsxRepo) AddActor(ctx context.Context, actor *models.ActorItem) (uint64, error) {
	err := repo.db.QueryRowContext(ctx, "INSERT INTO actor(name, gen, birthdate) VALUES($1, $2, $3) RETURNING id", actor.Name, actor.Gender, actor.Birthday).Scan(&actor.Id)
	if err != nil {
		return 0, fmt.Errorf("add actor error: %s", err.Error())
	}

	return actor.Id, nil
}

func (repo *PsxRepo) UpdateActor(ctx context.Context, actor *models.ActorRequest) error {
	if actor.Id == 0 {
		return fmt.Errorf("actor id missing")
	}

	var s strings.Builder
	var haveEq bool = false
	params := make([]any, 0, 4)
	count := 0

	s.WriteString("UPDATE actor SET ")
	if actor.Name != "" {
		if haveEq {
			s.WriteString(",")
		} else {
			haveEq = true
		}

		params = append(params, actor.Name)
		count++
		s.WriteString("name = $" + strconv.Itoa(count))
	}
	if actor.Birthday != "" {
		if haveEq {
			s.WriteString(",")
		} else {
			haveEq = true
		}

		params = append(params, actor.Birthday)
		count++
		s.WriteString(" birthdate = $" + strconv.Itoa(count))
	}
	if actor.Gender != "" {
		if haveEq {
			s.WriteString(",")
		} else {
			haveEq = true
		}
		params = append(params, actor.Gender)
		count++
		s.WriteString(" gen = $" + strconv.Itoa(count))
	}

	params = append(params, actor.Id)
	count++
	s.WriteString(" WHERE actor.id = $" + strconv.Itoa(count))

	if count < 2 {
		return fmt.Errorf("not have params")
	}

	_, err := repo.db.Query(s.String(), params...)
	if err != nil {
		return fmt.Errorf("update actor error: %s", err.Error())
	}

	for range actor.Films {
		existingFilmIds, err := repo.GetRelationByActorId(ctx, actor.Id)
		if err != nil {
			return err
		}

		for _, existingFilmId := range existingFilmIds {
			found := false
			for _, filmId := range actor.Films {
				if existingFilmId == filmId {
					found = true
					break
				}
			}
			if !found {
				err := repo.DeleteRelation(ctx, existingFilmId, actor.Id)
				if err != nil {
					return err
				}
			}
		}

		for _, filmId := range actor.Films {
			found := false
			for _, existingFilmId := range existingFilmIds {
				if existingFilmId == filmId {
					found = true
					break
				}
			}
			if !found {
				err := repo.InsertRelation(ctx, filmId, actor.Id)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (repo *PsxRepo) AddActorsForFilm(ctx context.Context, filmId uint64, actors []uint64) error {
	if len(actors) == 0 {
		return nil
	}

	var s strings.Builder
	var params []interface{}
	params = append(params, filmId)

	s.WriteString("INSERT INTO actor_in_film(id_film, id_actor) VALUES")
	for i, actor := range actors {
		if i != 0 {
			s.WriteString(",")
		}
		s.WriteString("($1, $" + strconv.Itoa(i+2) + ")")
		params = append(params, actor)
	}

	reasult, err := repo.db.ExecContext(ctx, s.String(), params...)
	if err != nil && reasult != nil {
		return fmt.Errorf("add film actors error: %w", err)
	}
	return nil
}

func (repo *PsxRepo) GetUser(ctx context.Context, login string, password []byte) (*models.UserItem, bool, error) {
	post := &models.UserItem{}

	err := repo.db.QueryRowContext(ctx, "SELECT profile.id, profile.login, profile.role FROM profile "+
		"WHERE profile.login = $1 AND profile.password = $2 ", login, password).Scan(&post.Id, &post.Login, &post.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("get query user error: %s", err.Error())
	}

	return post, true, nil
}

func (repo *PsxRepo) FindUser(ctx context.Context, login string) (bool, error) {
	post := &models.UserItem{}

	err := repo.db.QueryRowContext(ctx,
		"SELECT login FROM profile "+
			"WHERE login = $1", login).Scan(&post.Login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("find user query error: %s", err.Error())
	}

	return true, nil
}

func (repo *PsxRepo) CreateUser(ctx context.Context, login string, password []byte) error {
	var userID uint64
	err := repo.db.QueryRowContext(ctx, "INSERT INTO profile(login, role, password) VALUES($1, $2, $3) RETURNING id", login, "user", password).Scan(&userID)
	if err != nil {
		return fmt.Errorf("create user error: %s", err.Error())
	}

	return nil
}

func (repo *PsxRepo) GetUserId(ctx context.Context, login string) (uint64, error) {
	var userID uint64

	err := repo.db.QueryRowContext(ctx,
		"SELECT profile.id FROM profile WHERE profile.login = $1", login).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("user not found for login: %s", login)
		}
		return 0, fmt.Errorf("get userpro file id error: %s", err.Error())
	}

	return userID, nil
}

func (repo *PsxRepo) GetRole(ctx context.Context, userId uint64) (string, error) {
	var roleName string

	err := repo.db.QueryRowContext(ctx, "SELECT profile.role FROM profile  WHERE profile.id = $1", userId).Scan(&roleName)
	if err != nil {
		return "", fmt.Errorf("get user role err: %s", err.Error())
	}

	return roleName, nil
}
