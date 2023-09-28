package radarr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/varoOP/malbrr/internal/config"
)

type Movie struct {
	Title               string          `json:"title"`
	QualityProfileID    int             `json:"qualityProfileId"`
	MinimumAvailability MovieStatusType `json:"minimumAvailability,omitempty"`
	Monitored           bool            `json:"monitored"`
	TmdbID              int             `json:"tmdbId"`
	Tags                []int           `json:"tags,omitempty"`
	RootFolderPath      string          `json:"rootFolderPath"`
	AddOptions          AddMovieOptions `json:"addOptions"`
}

type MovieStatusType string

const (
	MovieStatusTypeTBA       MovieStatusType = "tba"
	MovieStatusTypeAnnounced MovieStatusType = "announced"
	MovieStatusTypeInCinemas MovieStatusType = "inCinemas"
	MovieStatusTypeReleased  MovieStatusType = "released"
	MovieStatusTypeDeleted   MovieStatusType = "deleted"
)

type AddMovieOptions struct {
	IgnoreEpisodesWithFiles    bool           `json:"ignoreEpisodesWithFiles"`
	IgnoreEpisodesWithoutFiles bool           `json:"ignoreEpisodesWithoutFiles"`
	Monitor                    MonitorTypes   `json:"monitor"`
	SearchForMovie             bool           `json:"searchForMovie"`
	AddMethod                  AddMovieMethod `json:"addMethod"`
}

type MonitorTypes string

const (
	MonitorTypeMovieOnly          MonitorTypes = "movieOnly"
	MonitorTypeMovieAndCollection MonitorTypes = "movieAndCollection"
	MonitorTypeNone               MonitorTypes = "none"
)

type AddMovieMethod string

const (
	AddMovieMethodManual     AddMovieMethod = "manual"
	AddMovieMethodList       AddMovieMethod = "list"
	AddMovieMethodCollection AddMovieMethod = "collection"
)

type Tag struct {
	Id    int    `json:"id"`
	Label string `json:"label"`
}

type RadarrError struct {
	PropertyName                      string `json:"propertyName"`
	ErrorMessage                      string `json:"errorMessage"`
	AttemptedValue                    int    `json:"attemptedValue"`
	Severity                          string `json:"severity"`
	ErrorCode                         string `json:"errorCode"`
	FormattedMessageArguments         []any  `json:"formattedMessageArguments"`
	FormattedMessagePlaceholderValues struct {
		PropertyName  string `json:"propertyName"`
		PropertyValue int    `json:"propertyValue"`
	} `json:"formattedMessagePlaceholderValues"`
}

// [
//     {
//         "propertyName": "TmdbId",
//         "errorMessage": "This movie has already been added",
//         "attemptedValue": 1056803,
//         "severity": "error",
//         "errorCode": "MovieExistsValidator",
//         "formattedMessageArguments": [],
//         "formattedMessagePlaceholderValues": {
//             "propertyName": "Tmdb Id",
//             "propertyValue": 1056803
//         }
//     }
// ]

type Client struct {
	client *http.Client
	config *config.Config
}

func NewClient(cfg *config.Config) *Client {
	c := &http.Client{
		Transport: &config.ApiKeyTransport{ApiKey: cfg.Sonarr.ApiKey},
	}

	return &Client{
		config: cfg,
		client: c,
	}
}

// {
// 	"title": "yada",
// 	"qualityProfileId": 4,
// 	"minimumAvailability": "announced",
// 	"monitored": true,
// 	"tmdbId": 1056803,
// 	"rootFolderPath": "/home/varo/gdrive/Seedbox/Anime-Movies",
// 	"addOptions": {
// 		"monitor": "movieOnly",
// 		"searchForMovie": false,
// 		"addMethod": "manual"
// 	}
//   }

func (c *Client) AddMovie(title string, tmdbid int, tags []int) error {
	s := Movie{
		Title:               title,
		MinimumAvailability: MovieStatusType(c.config.Radarr.MinimumAvailability),
		QualityProfileID:    c.config.Radarr.QualityProfileID,
		Monitored:           c.config.Radarr.Monitored,
		TmdbID:              tmdbid,
		RootFolderPath:      c.config.Radarr.RootFolderPath,
		AddOptions: AddMovieOptions{
			IgnoreEpisodesWithFiles:    false,
			IgnoreEpisodesWithoutFiles: false,
			Monitor:                    MonitorTypes(c.config.Radarr.MonitorType),
			SearchForMovie:             c.config.Radarr.SearchForMovie,
			AddMethod:                  AddMovieMethod(c.config.Radarr.AddMethod),
		},
		Tags: tags,
	}

	p, err := json.Marshal(s)
	if err != nil {
		return err
	}

	_, radarrError, err := c.SendPostRequest(c.config.Radarr.Url.JoinPath("/api/v3/movie").String(), p)
	if err != nil {
		return err
	}

	if radarrError.ErrorMessage == "This movie has already been added" {
		return nil
	}

	return nil
}

func (c *Client) AddTag(label string) (int, error) {
	t := Tag{
		Label: label,
	}

	p, err := json.Marshal(t)
	if err != nil {
		return -1, err
	}

	url := c.config.Radarr.Url.JoinPath("/api/v3/tag").String()

	resp, _, err := c.SendPostRequest(url, p)
	if err != nil {
		return -1, err
	}

	err = json.Unmarshal(resp, &t)
	if err != nil {
		return -1, err
	}

	return t.Id, nil
}

func (c *Client) TagExists(label string) (bool, int, error) {
	var tags []Tag
	url := c.config.Radarr.Url.JoinPath("/api/v3/tag").String()
	resp, err := c.client.Get(url)
	if err != nil {
		return false, -1, err
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&tags)
	if err != nil {
		return false, -1, err
	}

	for _, v := range tags {
		if v.Label == label {
			return true, v.Id, nil
		}
	}

	return false, -1, nil
}

func (c *Client) SendPostRequest(url string, body []byte) ([]byte, *RadarrError, error) {
	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()
	rb, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		radarrError := &RadarrError{}
		err = json.Unmarshal(rb, radarrError)
		if err != nil {
			return nil, nil, fmt.Errorf("response from Sonarr:\n%v", string(rb))
		}

		return nil, radarrError, nil
	}

	return rb, nil, nil
}
