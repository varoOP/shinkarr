package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

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

func (db *DB) GetIDs(malids []int32, dbtype string) (map[string]int32, error) {
	var (
		notFound []string
		// found    []string
	)

	m := map[string]int32{}
	sqlstmt := fmt.Sprintf("SELECT title,%v_id from anime where mal_id=?", dbtype)
	tx, err := db.Handler.Begin()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()
	for _, malid := range malids {
		var (
			id    int32
			title string
		)

		row := tx.QueryRow(sqlstmt, malid)
		err := row.Scan(&title, &id)
		if err != nil {
			return nil, err
		}

		titleLink := fmt.Sprintf("%v (https://myanimelist.net/anime/%v)", title, malid)
		if id <= 0 {
			s, a, err := NewAnimeMaps()
			if err != nil {
				return nil, err
			}

			if dbtype == "tvdb" {
				id = int32(s.CheckMap(int(malid)))
			}

			if dbtype == "tmdb" {
				id = int32(a.CheckMap(int(malid)))
			}

			if id <= 0 {
				notFound = append(notFound, fmt.Sprintf("Title: %v\nLink: https://myanimelist.net/anime/%v\n", title, malid))
				continue
			}
		}
		// found = append(found, title)
		m[titleLink] = id
	}

	// if len(found) > 0 {
	// 	log.Printf("%vids for the following anime were found:\n%v", dbtype, strings.Join(found, "\n"))
	// }

	if len(notFound) > 0 {
		fmt.Printf("\n%vids for the following anime were not found (Total: %v):\n%v", dbtype, len(notFound), strings.Join(notFound, "\n"))
		fmt.Println()
	}

	tx.Commit()
	if dbtype == "tvdb" {
		fmt.Printf("Total number of anime series that can be added: %v\n", len(malids)-len(notFound))
	}

	return m, nil
}

func check(err error) {
	if err != nil {
		log.Fatalf("database error: %v", err)
	}
}
