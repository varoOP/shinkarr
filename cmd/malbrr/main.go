package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/autobrr/omegabrr/pkg/autobrr"
	"github.com/mitchellh/go-homedir"
	"github.com/nstratos/go-myanimelist/mal"
	"github.com/spf13/pflag"
	"github.com/varoOP/malbrr/internal/config"
	"github.com/varoOP/malbrr/internal/database"
	"github.com/varoOP/malbrr/internal/maloauth"
	"github.com/varoOP/malbrr/internal/sonarr"
)

func main() {
	var (
		configPath string
		dbPath     string
		seasonYear int
		season     string
	)

	d, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}

	pflag.StringVar(&dbPath, "shinkuro-db", filepath.Join(d, ".config/shinkuro/shinkuro.db"), "path to shinkuro.db")
	pflag.StringVar(&configPath, "config", filepath.Join(d, ".config/malbrr"), "path to malbrr configuration directory")
	pflag.IntVar(&seasonYear, "season-year", 0, "season year of anime")
	pflag.StringVar(&season, "season", "", "season of anime")
	pflag.Parse()

	if seasonYear == 0 || season == "" {
		log.Fatal("season-year or season not provided")
	}

	dsn := dbPath + "?_pragma=busy_timeout%3d1000"
	db := database.NewDB(dsn)
	cfg := config.NewConfig(configPath)
	oc := maloauth.NewOauth2Client(db)
	c := mal.NewClient(oc)

	a, _, err := c.Anime.Seasonal(
		context.Background(),
		seasonYear,
		mal.AnimeSeason(season),
		mal.Fields{
			"alternative_titles{en}",
			"my_list_status{status}",
		},
		mal.NSFW(true),
		mal.Limit(500),
		mal.SortSeasonalByAnimeNumListUsers,
	)
	if err != nil {
		log.Fatal(err)
	}

	aa := []mal.Anime{}
	malIds := []int{}
	for _, anime := range a {
		if anime.MyListStatus.Status == mal.AnimeStatusPlanToWatch || anime.MyListStatus.Status == mal.AnimeStatusWatching {
			aa = append(aa, anime)
			malIds = append(malIds, anime.ID)
		}
	}

	anime, err := db.GetTvdbID(malIds)
	if err != nil {
		log.Fatal(err)
	}

	s := sonarr.NewClient(cfg)
	for title, id := range anime {
		err := s.AddSeries(title, id)
		if err != nil {
			log.Printf("unable to add %v\n%v\n", title, err)
			continue
		}

		log.Println("added", title)
	}

	var shows []string
	for _, anime := range aa {
		if anime.AlternativeTitles.En != "" {
			shows = append(shows, "*"+anime.AlternativeTitles.En+"*")
		} else {
			shows = append(shows, "*"+anime.Title+"*")
		}
	}

	f := autobrr.UpdateFilter{
		Shows: strings.Join(shows, ","),
	}

	ac := autobrr.NewClient(cfg.AutobrrUrl.String(), cfg.K.MustString("autobrr.ApiKey"))
	err = ac.UpdateFilterByID(context.Background(), cfg.K.MustInt("autobrr.FilterID"), f)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(shows)
}
