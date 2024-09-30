package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/nstratos/go-myanimelist/mal"
	"github.com/spf13/pflag"
	"github.com/varoOP/shinkarr/internal/config"
	"github.com/varoOP/shinkarr/internal/database"
	"github.com/varoOP/shinkarr/internal/maloauth"
	"github.com/varoOP/shinkarr/internal/omegabrr"
	"github.com/varoOP/shinkarr/internal/radarr"
	"github.com/varoOP/shinkarr/internal/sonarr"
)

func main() {
	var (
		configPath string
		seasonYear int
		season     string
	)

	d, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}

	pflag.StringVar(&configPath, "config", filepath.Join(d, ".config/shinkarr"), "path to shinkarr configuration directory")
	pflag.IntVar(&seasonYear, "season-year", 0, "season year of anime")
	pflag.StringVar(&season, "season", "", "season of anime")
	pflag.Parse()

	cfg := config.NewConfig(configPath)
	if seasonYear == 0 || season == "" {
		log.Println("Season and SeasonYear not provided, detecting current season")
		season, seasonYear = config.DetectSeasonYear()
		log.Printf("Current season: %v-%v\n", season, seasonYear)
	}

	dsn := cfg.Shinkro.DBPath + "?_pragma=busy_timeout%3d1000"
	db := database.NewDB(dsn)
	oc := maloauth.NewOauth2Client(db)
	c := mal.NewClient(oc)

	a, _, err := c.Anime.Seasonal(
		context.Background(),
		seasonYear,
		mal.AnimeSeason(season),
		mal.Fields{
			"alternative_titles{en}",
			"my_list_status{status}",
			"media_type",
		},
		mal.NSFW(true),
		mal.Limit(500),
		mal.SortSeasonalByAnimeNumListUsers,
	)
	if err != nil {
		log.Fatal(err)
	}

	//aa := []mal.Anime{}
	malIdsSeries := []int32{}
	malIdsMovies := []int32{}
	for _, anime := range a {
		if anime.MyListStatus.Status == mal.AnimeStatusPlanToWatch || anime.MyListStatus.Status == mal.AnimeStatusWatching {
			//		aa = append(aa, anime)
			if anime.MediaType == "movie" {
				malIdsMovies = append(malIdsMovies, int32(anime.ID))
				continue
			}

			malIdsSeries = append(malIdsSeries, int32(anime.ID))
		}
	}
	fmt.Println("Total number of anime series we wish to add: ", len(malIdsSeries))
	fmt.Println("Total number of anime movies we wish to add: ", len(malIdsMovies))

	animeTv, err := db.GetIDs(malIdsSeries, "tvdb")
	if err != nil {
		log.Fatal(err)
	}

	animeMovie, err := db.GetIDs(malIdsMovies, "tmdb")
	if err != nil {
		log.Fatal(err)
	}

	s := sonarr.NewClient(cfg)
	tag := fmt.Sprintf("%v-%v", season, seasonYear)
	tagExists, tagId, err := s.TagExists(tag)
	if err != nil {
		log.Fatal(err)
	}

	if !tagExists {
		tagId, err = s.AddTag(tag)
		if err != nil {
			log.Fatal(err)
		}
	}

	seriesAdded := []string{}
	seriesNotAdded := []string{}

	for title, id := range animeTv {
		err := s.AddSeries(title, id, []int32{tagId})
		if err != nil {
			seriesNotAdded = append(seriesNotAdded, fmt.Sprintf("%v\nerror:%v\n", title, err))
			continue
		}

		seriesAdded = append(seriesAdded, title)
	}

	if len(seriesAdded) > 0 {
		fmt.Printf("\nFollowing series added (%v):\n", len(seriesAdded))
		for _, v := range seriesAdded {
			fmt.Println(v)
		}
	}

	if len(seriesNotAdded) > 0 {
		fmt.Printf("\nFollowing series not added (%v):\n", len(seriesNotAdded))
		for _, v := range seriesNotAdded {
			fmt.Println(v)
		}
	}

	m := radarr.NewClient(cfg)
	tag = fmt.Sprintf("%v-%v", season, seasonYear)
	tagExists, tagId, err = m.TagExists(tag)
	if err != nil {
		log.Fatal(err)
	}

	if !tagExists {
		tagId, err = m.AddTag(tag)
		if err != nil {
			log.Fatal(err)
		}
	}

	moviesAdded := []string{}
	moviesNotAdded := []string{}

	for title, id := range animeMovie {
		err := m.AddMovie(title, id, []int32{tagId})
		if err != nil {
			moviesNotAdded = append(moviesNotAdded, fmt.Sprintf("%v\nerror:%v\n", title, err))
			continue
		}

		moviesAdded = append(moviesAdded, title)
	}

	if len(moviesAdded) > 0 {
		fmt.Printf("\nFollowing movies added (%v):\n", len(moviesAdded))
		for _, v := range moviesAdded {
			fmt.Println(v)
		}
	}

	if len(moviesNotAdded) > 0 {
		fmt.Printf("\nFollowing movies not added (%v):\n", len(moviesNotAdded))
		for _, v := range moviesNotAdded {
			fmt.Println(v)
		}
	}

	fmt.Println("\nRunning omegabrr:")
	om := omegabrr.NewOmegabrr(cfg.Omegabrr.ConfigPath, cfg.Omegabrr.ExecPath, season, seasonYear)
	err = om.Run()
	if err != nil {
		log.Fatal(err)
	}
}
