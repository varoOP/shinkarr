package config

import (
	"log"
	"net/url"
	"path/filepath"
	"strconv"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
)

type Config struct {
	Sonarr *SonarrConfig
	Radarr *RadarrConfig
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
	QualityProfileID int32    `koanf:"QualityProfileID"`
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
	AddMethod           string `koanf:"AddMethod"`
	MinimumAvailability string `konaf:"MinimumAvailability"`
	QualityProfileID    int    `koanf:"QualityProfileID"`
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
	k.Unmarshal("sonarr", &s)
	k.Unmarshal("radarr", &r)
	s.BuildUrl()
	r.BuildUrl()

	return &Config{
		Sonarr: &s,
		Radarr: &r,
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
