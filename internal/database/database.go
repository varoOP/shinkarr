package database

import (
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"encoding/hex"
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

func (db *DB) GetMalCreds(encryptionKey string) map[string]string {
	var (
		client_id     []byte
		client_secret []byte
		access_token  []byte
		token_iv      []byte
	)

	sqlstmt := "SELECT client_id, client_secret, access_token, token_iv from malauth;"

	row := db.Handler.QueryRow(sqlstmt)
	err := row.Scan(&client_id, &client_secret, &access_token, &token_iv)
	if err != nil {
		check(err)
	}

	// Decrypt all fields using the encryption key and IV
	var decryptedClientID, decryptedClientSecret, decryptedAccessToken string

	if encryptionKey != "" {
		// Decrypt client_id
		if len(client_id) > 0 {
			decrypted, err := decrypt(client_id, token_iv, encryptionKey)
			if err != nil {
				log.Fatalf("failed to decrypt client_id: %v", err)
			}
			decryptedClientID = string(decrypted)
		}

		// Decrypt client_secret
		if len(client_secret) > 0 {
			decrypted, err := decrypt(client_secret, token_iv, encryptionKey)
			if err != nil {
				log.Fatalf("failed to decrypt client_secret: %v", err)
			}
			decryptedClientSecret = string(decrypted)
		}

		// Decrypt access_token
		if len(access_token) > 0 {
			decrypted, err := decrypt(access_token, token_iv, encryptionKey)
			if err != nil {
				log.Fatalf("failed to decrypt access_token: %v", err)
			}
			decryptedAccessToken = string(decrypted)
		}
	} else {
		// Fallback to plain text (backward compatibility)
		decryptedClientID = string(client_id)
		decryptedClientSecret = string(client_secret)
		decryptedAccessToken = string(access_token)
	}

	return map[string]string{
		"client_id":     decryptedClientID,
		"client_secret": decryptedClientSecret,
		"access_token":  decryptedAccessToken,
	}
}

// decrypt matches shinkro's Decrypt function signature
// It takes ciphertext and iv as separate byte slices, matching shinkro's implementation
func decrypt(ciphertext, iv []byte, encryptionKey string) ([]byte, error) {
	// Get and decode the encryption key
	key, err := getEncryptionKey(encryptionKey)
	if err != nil {
		return nil, err
	}

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %v", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM mode: %v", err)
	}

	// Decrypt using the provided IV
	plaintext, err := gcm.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %v", err)
	}

	return plaintext, nil
}

// getEncryptionKey matches shinkro's getEncryptionKey function
// It decodes the encryption key from hex format and requires exactly 32 bytes
func getEncryptionKey(encryptionKey string) ([]byte, error) {
	key, err := hex.DecodeString(encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("invalid hex encryption key: %v", err)
	}

	if len(key) != 32 {
		return nil, fmt.Errorf("encryption key must be 32 bytes, got %d bytes", len(key))
	}

	return key, nil
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

	stmt, err := tx.Prepare(sqlstmt)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	for _, malid := range malids {
		var (
			id    int32
			title string
		)

		row := stmt.QueryRow(malid)
		err := row.Scan(&title, &id)
		if err != nil {
			if err == sql.ErrNoRows {
				notFound = append(notFound, fmt.Sprintf("Title: %v\nLink: https://myanimelist.net/anime/%v\n", "No Title", malid))
				continue
			}
			return nil, err
		}

		titleLink := fmt.Sprintf("%v (https://myanimelist.net/anime/%v)", title, malid)
		// if id > 0 {
		// 	log.Printf("%v tvdbid found in db: %v\n", titleLink, id)
		// }
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
		fmt.Printf("\nTotal number of anime series that can be added: %v\n", len(malids)-len(notFound))
	}

	if dbtype == "tmdb" {
		fmt.Printf("Total number of anime movies that can be added: %v\n", len(malids)-len(notFound))
	}

	return m, nil
}

func check(err error) {
	if err != nil {
		log.Fatalf("database error: %v", err)
	}
}
