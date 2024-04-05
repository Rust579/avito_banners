package postgres

import (
	"avito_banners/internal/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"strconv"
)

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

	psgDb, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}

	if err = psgDb.Ping(); err != nil {
		return err
	}

	if err = createBannersTable(psgDb); err != nil {
		return err
	}

	return nil
}

func createBannersTable(db *sql.DB) error {
	var query = []string{
		`
		CREATE TABLE IF NOT EXISTS banners (
			id INT PRIMARY KEY,
			title VARCHAR(100),
		    text VARCHAR(100),
			url VARCHAR(200)
		)
		`,
	}

	for _, q := range query {
		_, err := db.Exec(q)
		if err != nil {
			return err
		}
	}

	return nil
}
