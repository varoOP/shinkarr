package database

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

type DB struct {
	Handler *sql.DB
}

func NewDB(DSN string) *DB {
	db := &DB{}
	var err error
	db.Handler, err = sql.Open("sqlite", DSN)
	check(err)
	if _, err = db.Handler.Exec(`PRAGMA journal_mode = wal;`); err != nil {
		check(err)
	}

	return db
}

func (db *DB) GetMalCreds() map[string]string {
	var (
		client_id     string
		client_secret string
		access_token  string
	)

	sqlstmt := "SELECT * from malauth;"

	row := db.Handler.QueryRow(sqlstmt)
	err := row.Scan(&client_id, &client_secret, &access_token)
	if err != nil {
		check(err)
	}

	return map[string]string{
		"client_id":     client_id,
		"client_secret": client_secret,
		"access_token":  access_token,
	}
}

func (db *DB) GetTvdbID(malid []int) (map[string]int, error) {
	var notFound []int
	m := map[string]int{}
	sqlstmt := "SELECT title,tvdb_id from anime where mal_id=?"
	tx, err := db.Handler.Begin()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()
	for _, id := range malid {
		var (
			tvdbId int
			title  string
		)

		row := tx.QueryRow(sqlstmt, id)
		err := row.Scan(&title, &tvdbId)
		if err != nil {
			return nil, err
		}

		if tvdbId == 0 {
			notFound = append(notFound, id)
			continue
		}

		m[title] = tvdbId
	}

	if len(notFound) > 0 {
		log.Println("tvdbIds for the following malIds were not found:", notFound)
	}

	tx.Commit()
	return m, nil
}

func check(err error) {
	if err != nil {
		log.Fatalf("database error: %v", err)
	}
}
