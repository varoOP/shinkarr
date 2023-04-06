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
	SonarrUrl  url.URL
	AutobrrUrl url.URL
	K          *koanf.Koanf
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

	c := &Config{
		K: k,
	}

	targets := []string{"sonarr", "autobrr"}
	for _, target := range targets {
		c.buildUrl(target, k)
	}

	return c
}

func (c *Config) buildUrl(target string, k *koanf.Koanf) {
	var url url.URL
	if k.Bool(target + ".TLS") {
		url.Scheme = "https"
	} else {
		url.Scheme = "http"
	}

	url.Host = k.MustString(target+".Host") + ":" + strconv.Itoa(k.MustInt(target+".Port"))
	url = *url.JoinPath(k.String(target + ".BaseUrl"))
	if target == "sonarr" {
		c.SonarrUrl = url
	}

	if target == "autobrr" {
		c.AutobrrUrl = url
	}
}
