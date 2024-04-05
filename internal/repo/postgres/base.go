package postgres

import (
	"avito_banners/internal/config"
	"avito_banners/internal/model"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"strconv"
)

var psgDb *sql.DB

func Init() error {

	host := config.Cfg.Postgres.Host
	port, err := strconv.Atoi(config.Cfg.Postgres.Port)
	if err != nil {
		return err
	}
	user := config.Cfg.Postgres.User
	password := config.Cfg.Postgres.Password
	dbname := config.Cfg.Postgres.Dbname

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	psgDb, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}

	if err = psgDb.Ping(); err != nil {
		return err
	}

	if err = createBannersTable(); err != nil {
		return err
	}

	return nil
}

func createBannersTable() error {
	var query = []string{
		`
		CREATE TABLE IF NOT EXISTS banners (
			id SERIAL PRIMARY KEY,
			title VARCHAR(100),
		    text VARCHAR(100),
			url VARCHAR(200)
		)
		`,
	}

	for _, q := range query {
		_, err := psgDb.Exec(q)
		if err != nil {
			return err
		}
	}

	return nil
}

func InsertBanner(banner model.Banner) error {

	query := `
		INSERT INTO banners (text, title, url)
		VALUES ($1, $2, $3)
	`
	_, err := psgDb.Exec(query, banner.Text, banner.Title, banner.Url)
	if err != nil {
		return err
	}

	return nil
}
