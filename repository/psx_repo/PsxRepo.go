package psx_repo

import (
	"database/sql"
	"errors"
	"filmoteka/configs"
	"filmoteka/pkg/models"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/lib/pq"
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

func (repo *PsxRepo) GetAllFilms(start uint64, end uint64) (*[]models.FilmItem, error) {
	films := make([]models.FilmItem, 0, end-start)

	rows, err := repo.DB.Query("SELECT film.id, film.title, film.rating, film.info from film "+
		"ORDER BY rating DESC "+
		"OFFSET $1 LIMIT $2",
		start, end)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("GetFilms error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		post := models.FilmItem{}
		err = rows.Scan(&post.Id, &post.Info, &post.Rating, &post.ReleaseDate)
		if err != nil {
			return nil, fmt.Errorf("GetFilms scan error: %w", err)
		}

		films = append(films, post)
	}

	return &films, nil
}

func (repo *PsxRepo) GetFilms(request *models.FindFilmRequest) (*[]models.FilmItem, error) {
	//films := make([]models.FilmItem, 0, request.PerPage)
	var s strings.Builder
	var hasWhere bool
	paramNum := 1
	var params []interface{}

	s.WriteString("SELECT DISTINCT film.title, film.id, film.rating, film.release_date FROM film " +
		"JOIN actor_in_film ON film.id = actor_in_film.id_film " +
		"JOIN actor ON actor.id = actor_in_film.id_actor ")

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
		params = append(params, request.ReleaseDateTo)
	}
	if request.Actor != "" {
		if !hasWhere {
			s.WriteString("WHERE ")
		} else {
			s.WriteString("AND ")
		}
		s.WriteString("(CASE WHEN array_length($" + strconv.Itoa(paramNum) + "::varchar[], 1)> 0 " +
			"THEN crew.name = ANY ($" + strconv.Itoa(paramNum) + "::varchar[]) ELSE TRUE END) ")
		paramNum++
		params = append(params, pq.Array(request.Actor))
	}
	s.WriteString("ORDER BY film.rating " +
		"LIMIT $" + strconv.Itoa(paramNum) + " OFFSET $" + strconv.Itoa(paramNum+1))
	params = append(params, request.RatingFrom, request.RatingTo, request.Page, request.PerPage)
	rows, err := repo.DB.Query(s.String(), params...)

	fmt.Println(s.String(), params, request)

	if err != nil {
		return nil, fmt.Errorf("find film err: %w", err)
	}
	defer rows.Close()

	fmt.Println(rows)

	return nil, err
}
