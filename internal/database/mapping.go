package database

import (
	"io"
	"net/http"

	"gopkg.in/yaml.v3"
)

const (
	communityMapTVDB = "https://github.com/varoOP/shinkro-mapping/raw/main/tvdb-mal.yaml"
	communityMapTMDB = "https://github.com/varoOP/shinkro-mapping/raw/main/tmdb-mal.yaml"
)

type AnimeTVDBMap struct {
	Anime []Anime `yaml:"AnimeMap" json:"AnimeMap"`
}

type Anime struct {
	Malid        int            `yaml:"malid" json:"malid"`
	Title        string         `yaml:"title" json:"title"`
	Type         string         `yaml:"type" json:"type"`
	Tvdbid       int            `yaml:"tvdbid" json:"tvdbid"`
	TvdbSeason   int            `yaml:"tvdbseason" json:"tvdbseason"`
	Start        int            `yaml:"start" json:"start"`
	UseMapping   bool           `yaml:"useMapping" json:"useMapping"`
	AnimeMapping []AnimeMapping `yaml:"animeMapping" json:"animeMapping"`
}

type AnimeMapping struct {
	TvdbSeason int `yaml:"tvdbseason" json:"tvdbseason"`
	Start      int `yaml:"start" json:"start"`
}

type AnimeMovies struct {
	AnimeMovie []AnimeMovie `yaml:"animeMovies" json:"animeMovies"`
}

type AnimeMovie struct {
	MainTitle string `yaml:"mainTitle" json:"mainTitle"`
	TMDBID    int    `yaml:"tmdbid" json:"tmdbid"`
	MALID     int    `yaml:"malid" json:"malid"`
}

func NewAnimeMaps() (*AnimeTVDBMap, *AnimeMovies, error) {
	return loadCommunityMaps()
}

func (s *AnimeTVDBMap) CheckMap(malid int) int {
	for _, anime := range s.Anime {
		if anime.Malid == malid {
			return anime.Tvdbid
		}
	}

	return 0
}

func (am *AnimeMovies) CheckMap(malid int) int {
	for _, animeMovie := range am.AnimeMovie {
		if animeMovie.MALID == malid {
			return animeMovie.TMDBID
		}
	}

	return 0
}

func loadCommunityMaps() (*AnimeTVDBMap, *AnimeMovies, error) {
	s := &AnimeTVDBMap{}
	respTVDB, err := http.Get(communityMapTVDB)
	if err != nil {
		return nil, nil, err
	}

	err = readYamlHTTP(respTVDB, s)
	if err != nil {
		return nil, nil, err
	}

	am := &AnimeMovies{}
	respTMDB, err := http.Get(communityMapTMDB)
	if err != nil {
		return nil, nil, err
	}

	err = readYamlHTTP(respTMDB, am)
	if err != nil {
		return nil, nil, err
	}

	return s, am, nil
}

func readYamlHTTP(resp *http.Response, mapping interface{}) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	err = yaml.Unmarshal(body, mapping)
	if err != nil {
		return err
	}

	return nil
}
