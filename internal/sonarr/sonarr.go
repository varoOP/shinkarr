package sonarr

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/varoOP/malbrr/internal/config"
)

const SeriesAlreadyAdded = "This series has already been added"

type Series struct {
	Title            string           `json:"title"`
	QualityProfileID int              `json:"qualityProfileId"`
	SeasonFolder     bool             `json:"seasonFolder"`
	Monitored        bool             `json:"monitored"`
	TvdbID           int              `json:"tvdbId"`
	SeriesType       string           `json:"seriesType"`
	RootFolderPath   string           `json:"rootFolderPath"`
	AddOptions       AddSeriesOptions `json:"addOptions"`
}

type AddSeriesOptions struct {
	IgnoreEpisodesWithFiles       bool         `json:"ignoreEpisodesWithFiles"`
	IgnoreEpisodesWithoutFiles    bool         `json:"ignoreEpisodesWithoutFiles"`
	Monitor                       MonitorTypes `json:"monitor"`
	SearchForMissingEpisodes      bool         `json:"searchForMissingEpisodes"`
	SearchForCutoffUnmentEpisodes bool         `json:"searchForCutoffUnmetEpisodes"`
}

type MonitorTypes string

const (
	MonitorTypeAll               MonitorTypes = "all"
	MonitorTypeFuture            MonitorTypes = "future"
	MonitorTypeUnknown           MonitorTypes = "unknown"
	MonitorTypeMissing           MonitorTypes = "missing"
	MonitorTypeExisting          MonitorTypes = "existing"
	MonitorTypeFirstSeason       MonitorTypes = "firstSeason"
	MonitorTypeLatestSeason      MonitorTypes = "latestSeason"
	MonitorTypePilot             MonitorTypes = "pilot"
	MonitorTypeMonitorSpecials   MonitorTypes = "monitorSpecials"
	MonitorTypeUnmonitorSpecials MonitorTypes = "UnmonitorSpecials"
	MonitorTypeNone              MonitorTypes = "none"
)

type apiKeyTransport struct {
	Transport http.RoundTripper
	ApiKey    string
}

func (c *apiKeyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if c.Transport == nil {
		c.Transport = http.DefaultTransport
	}

	req.Header.Add("X-Api-Key", c.ApiKey)
	return c.Transport.RoundTrip(req)
}

type Client struct {
	client *http.Client
	config *config.Config
}

func NewClient(cfg *config.Config) *Client {
	c := &http.Client{
		Transport: &apiKeyTransport{ApiKey: cfg.K.MustString("sonarr.ApiKey")},
	}

	return &Client{
		config: cfg,
		client: c,
	}
}

func (c *Client) AddSeries(title string, tvdbid int) error {
	s := Series{
		Title:            title,
		QualityProfileID: c.config.K.MustInt("sonarr.QualityProfileID"),
		SeasonFolder:     c.config.K.Bool("sonarr.SeasonFolder"),
		Monitored:        c.config.K.Bool("sonarr.Monitored"),
		TvdbID:           tvdbid,
		SeriesType:       "anime",
		RootFolderPath:   c.config.K.MustString("sonarr.RootFolderPath"),
		AddOptions: AddSeriesOptions{
			IgnoreEpisodesWithFiles:       false,
			IgnoreEpisodesWithoutFiles:    false,
			Monitor:                       MonitorTypes(c.config.K.MustString("sonarr.MonitorType")),
			SearchForMissingEpisodes:      true,
			SearchForCutoffUnmentEpisodes: false,
		},
	}

	p, err := json.Marshal(s)
	if err != nil {
		return err
	}

	url := c.config.SonarrUrl.JoinPath("/api/v3/series").String()
	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(p))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		var i []interface{}
		err := json.NewDecoder(resp.Body).Decode(&i)
		if err != nil {
			return err
		}

		b, ok := i[0].(map[string]interface{})
		if !ok {
			return errors.New("failed to decode response from sonarr")
		}

		if b["errorMessage"] == SeriesAlreadyAdded {
			return nil
		}

		s, err := json.Marshal(i)
		if err != nil {
			return err
		}

		return errors.New(string(s))
	}

	return nil
}
