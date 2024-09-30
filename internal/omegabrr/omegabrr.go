package omegabrr

import (
	"fmt"
	"os"
	"os/exec"

	"gopkg.in/yaml.v3"
)

type Omegabrr struct {
	Season     string
	SeasonYear int
	ConfigPath string `koanf:"ConfigPath"`
	ExecPath   string `koanf:"ExecPath"`
}

type Config struct {
	Server   ServerConfig  `yaml:"server"`
	Schedule string        `yaml:"schedule"`
	Clients  ClientsConfig `yaml:"clients"`
}

type ServerConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	APIToken string `yaml:"apiToken"`
}

type ClientsConfig struct {
	Autobrr AutobrrConfig `yaml:"autobrr"`
	Arr     []ArrConfig   `yaml:"arr"`
}

type AutobrrConfig struct {
	Host   string `yaml:"host"`
	APIKey string `yaml:"apikey"`
}

type ArrConfig struct {
	Name         string   `yaml:"name"`
	Type         string   `yaml:"type"`
	Host         string   `yaml:"host"`
	APIKey       string   `yaml:"apikey"`
	Filters      []int    `yaml:"filters"`
	TagsInclude  []string `yaml:"tagsInclude"`
	TagsExclude  []string `yaml:"tagsExclude"`
	MatchRelease bool     `yaml:"matchRelease"`
}

func NewOmegabrr(configPath, execPath, season string, seasonYear int) *Omegabrr {
	return &Omegabrr{
		Season:     season,
		SeasonYear: seasonYear,
		ConfigPath: configPath,
		ExecPath:   execPath,
	}
}

func (om *Omegabrr) getTag() string {
	return fmt.Sprintf("%v-%v", om.Season, om.SeasonYear)
}

func (om *Omegabrr) updateConfig() error {
	data, err := os.ReadFile(om.ConfigPath)
	if err != nil {
		return err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	for i := range config.Clients.Arr {
		config.Clients.Arr[i].TagsInclude = []string{om.getTag()}
	}

	updatedData, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}

	err = os.WriteFile(om.ConfigPath, updatedData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (om *Omegabrr) Run() error {
	if err := om.updateConfig(); err != nil {
		return err
	}

	cmd := exec.Command(om.ExecPath, "arr", "--config="+om.ConfigPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
