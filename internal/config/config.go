package config

import (
	"log"
	"net/url"
	"path/filepath"
	"strconv"
	"time"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/varoOP/shinkarr/internal/omegabrr"
)

type Config struct {
	Sonarr   *SonarrConfig
	Radarr   *RadarrConfig
	Omegabrr *omegabrr.Omegabrr
	Shinkro  *ShinkroConfig
}

type ShinkroConfig struct {
	DBPath     string `koanf:"DBPath"`
	ConfigPath string `koanf:"ConfigPath"`
	EncryptionKey string
}

type SonarrConfig struct {
	Url              *url.URL
	Host             string `koanf:"Host"`
	Port             int    `koanf:"Port"`
	BaseUrl          string `koanf:"BaseUrl"`
	TLS              bool   `koanf:"TLS"`
	ApiKey           string `koanf:"ApiKey"`
	RootFolderPath   string `koanf:"RootFolderPath"`
	SeasonFolder     bool   `koanf:"SeasonFolder"`
	Monitored        bool   `koanf:"Monitored"`
	MonitorType      string `koanf:"MonitorType"`
	QualityProfileID int32  `koanf:"QualityProfileID"`
}

type RadarrConfig struct {
	Url                 *url.URL
	Host                string `koanf:"Host"`
	Port                int    `koanf:"Port"`
	BaseUrl             string `koanf:"BaseUrl"`
	TLS                 bool   `koanf:"TLS"`
	ApiKey              string `koanf:"ApiKey"`
	RootFolderPath      string `koanf:"RootFolderPath"`
	Monitored           bool   `koanf:"Monitored"`
	MonitorType         string `koanf:"MonitorType"`
	SearchForMovie      bool   `koanf:"SearchForMovie"`
	MinimumAvailability string `konaf:"MinimumAvailability"`
	QualityProfileID    int32  `koanf:"QualityProfileID"`
}

func NewConfig(dir string) *Config {
	if dir == "" {
		log.Fatal("config location not found")
	}

	configPath := filepath.Join(dir, "config.toml")
	k := koanf.New(".")
	if err := k.Load(file.Provider(configPath), toml.Parser()); err != nil {
		log.Fatal(err)
	}

	s := SonarrConfig{}
	r := RadarrConfig{}
	sh := ShinkroConfig{}
	om := omegabrr.Omegabrr{}
	k.Unmarshal("sonarr", &s)
	k.Unmarshal("radarr", &r)
	k.Unmarshal("shinkro", &sh)
	k.Unmarshal("omegabrr", &om)
	s.BuildUrl()
	r.BuildUrl()
	sh.LoadEncryptionKey()

	return &Config{
		Sonarr:   &s,
		Radarr:   &r,
		Shinkro:  &sh,
		Omegabrr: &om,
	}
}

func (s *SonarrConfig) BuildUrl() {
	scheme := "http"
	if s.TLS {
		scheme = "https"
	}

	url := url.URL{
		Scheme: scheme,
		Host:   s.Host + ":" + strconv.Itoa(s.Port),
	}

	s.Url = url.JoinPath(s.BaseUrl)
}

func (r *RadarrConfig) BuildUrl() {
	scheme := "http"
	if r.TLS {
		scheme = "https"
	}

	url := url.URL{
		Scheme: scheme,
		Host:   r.Host + ":" + strconv.Itoa(r.Port),
	}

	r.Url = url.JoinPath(r.BaseUrl)
}

func (sh *ShinkroConfig) LoadEncryptionKey() {
	if sh.ConfigPath == "" {
		log.Fatal("shinkro ConfigPath not set")
	}

	shinkroConfigPath := filepath.Join(sh.ConfigPath, "config.toml")
	k := koanf.New(".")
	if err := k.Load(file.Provider(shinkroConfigPath), toml.Parser()); err != nil {
		log.Fatalf("failed to load shinkro config: %v", err)
	}

	sh.EncryptionKey = k.String("EncryptionKey")
	if sh.EncryptionKey == "" {
		log.Fatal("EncryptionKey not found in shinkro config.toml")
	}
}

func DetectSeasonYear() (string, int) {
	var (
		season     string
		seasonYear int
	)
	now := time.Now()
	month := now.Month()
	year := now.Year()

	switch month {
	case time.January, time.February, time.March:
		season = "winter"
	case time.April, time.May, time.June:
		season = "spring"
	case time.July, time.August:
		season = "summer"
	case time.September, time.October, time.November, time.December:
		season = "fall"
	default:
		season = "unknown"
	}

	seasonYear = year

	return season, seasonYear
}
