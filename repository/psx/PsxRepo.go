package psx

import (
	"database/sql"
	"errors"
	"filmoteka/configs"
	"filmoteka/pkg/models"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

type PsxRepo struct {
	DB *sql.DB
}

func GetFilmRepo(config *configs.DbPsxConfig, log *logrus.Logger) (*PsxRepo, error) {
	dsn := fmt.Sprintf("user=%s dbname=%s password= %s host=%s port=%d sslmode=%s",
		config.User, config.Dbname, config.Password, config.Host, config.Port, config.Sslmode)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Error("sql open error", "err", err.Error())
		return nil, fmt.Errorf("get user repo err: %w", err)
	}
	err = db.Ping()
	if err != nil {
		log.Error("sql ping error", "err", err.Error())
		return nil, fmt.Errorf("get user repo err: %w", err)
	}
	db.SetMaxOpenConns(config.MaxOpenConns)

	log.Info("Psx created successful")
	return &PsxRepo{DB: db}, nil
}

func (repo *PsxRepo) GetFilms(request *models.FindFilmRequest) (*[]models.FilmItem, error) {
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
	//if request.Actor != "" {
	//	if !hasWhere {
	//		s.WriteString("WHERE ")
	//	} else {
	//		s.WriteString("AND ")
	//	}
	//	s.WriteString("(CASE WHEN array_length($" + strconv.Itoa(paramNum) + "::varchar[], 1)> 0 " +
	//		"THEN crew.name = ANY ($" + strconv.Itoa(paramNum) + "::varchar[]) ELSE TRUE END) ")
	//	paramNum++
	//	params = append(params, pq.Array(request.Actor))
	//}
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

	rows, err := repo.DB.Query(s.String(), params...)

	if err != nil {
		return nil, fmt.Errorf("find film err: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		post := models.FilmItem{}

		err := rows.Scan(&post.Id, &post.Title, &post.Rating, &post.Info, &post.ReleaseDate)
		if err != nil {
			return nil, fmt.Errorf("find film scan err: %w", err)
		}

		films = append(films, post)
	}

	return &films, err
}

func (repo *PsxRepo) AddFilm(film *models.FilmRequest) (uint64, error) {
	result, err := repo.DB.Exec("INSERT INTO film(title, info, release_date, rating ) VALUES($1, $2, $3, $4)", film.Title, film.Info, film.ReleaseDate, film.Rating)
	if err != nil {
		return 0, fmt.Errorf("AddFilm err: %w", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("getting LastInsertId err: %w", err)
	}

	return uint64(lastID), nil
}

func (repo *PsxRepo) AddActor(actor *models.ActorItem) (uint64, error) {
	result, err := repo.DB.Exec("INSERT INTO actor(name, gen, birthdate) VALUES($1, $2, $3)", actor.Name, actor.Gender, actor.Birthday)
	if err != nil {
		return 0, fmt.Errorf("AddActor err: %w", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("getting LastInsertId err: %w", err)
	}

	return uint64(lastID), nil
}

func (repo *PsxRepo) AddActorsForFilm(filmId uint64, actors []uint64) error {
	var s strings.Builder
	var params []interface{}
	params = append(params, filmId)

	s.WriteString("INSERT INTO actor_in_film(id_film, id_actor) VALUES")
	for i, actor := range actors {
		if i != 0 {
			s.WriteString(",")
		}
		s.WriteString("($1, $" + strconv.Itoa(i) + ")")
		params = append(params, actor)
	}

	_, err := repo.DB.Exec(s.String(), params...)
	if err != nil {
		return fmt.Errorf("add film actors error: %w", err)
	}
	return nil
}

func (repo *PsxRepo) GetUser(login string, password string) (*models.UserItem, bool, error) {
	post := &models.UserItem{}

	err := repo.DB.QueryRow("SELECT profile.id, profile.login, profile.password, profile.role FROM profile "+
		"WHERE login = $1 AND password = $2", login, password).Scan(&post.Id, &post.Login, &post.Password, &post.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("GetUser err: %w", err)
	}

	return post, true, nil
}

func (repo *PsxRepo) FindUser(login string) (bool, error) {
	post := &models.UserItem{}

	err := repo.DB.QueryRow(
		"SELECT login FROM profile "+
			"WHERE login = $1", login).Scan(&post.Login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("FindUser err: %w", err)
	}

	return true, nil
}

func (repo *PsxRepo) CreateUser(login string, password string) error {
	_, err := repo.DB.Exec(
		"INSERT INTO profile(login, password) "+
			"VALUES($1, $2)", login, password)
	if err != nil {
		return fmt.Errorf("CreateUser err: %w", err)
	}

	return nil
}
