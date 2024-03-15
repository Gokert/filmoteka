package psx_repo

import (
	"database/sql"
	"errors"
	"filmoteka/configs"
	"filmoteka/pkg/models"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/sirupsen/logrus"
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

func (repo *PsxRepo) GetAllFilms(start uint64, end uint64, order bool) (*[]models.FilmItem, error) {
	films := make([]models.FilmItem, 0, end-start)

	rows, err := repo.DB.Query("SELECT film.id, film.title, film.rating, film.info from film "+
		"ORDER BY rating ASC "+
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
