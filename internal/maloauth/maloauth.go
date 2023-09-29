package maloauth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/varoOP/shinkarr/internal/database"
	"golang.org/x/oauth2"
)

func NewOauth2Client(db *database.DB) *http.Client {
	ctx := context.Background()
	creds := db.GetMalCreds()
	cfg := &oauth2.Config{
		ClientID:     creds["client_id"],
		ClientSecret: creds["client_secret"],
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://myanimelist.net/v1/oauth2/authorize",
			TokenURL:  "https://myanimelist.net/v1/oauth2/token",
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}

	t := &oauth2.Token{}
	err := json.Unmarshal([]byte(creds["access_token"]), t)
	if err != nil {
		log.Fatalln(err)
	}

	fresh_token, err := cfg.TokenSource(ctx, t).Token()
	if err != nil {
		log.Fatal(err)
	}

	client := cfg.Client(ctx, fresh_token)
	return client
}
