package sonarr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/varoOP/shinkarr/internal/config"
)

type AddSeriesOptions struct {
	IgnoreEpisodesWithFiles      bool         `json:"ignoreEpisodesWithFiles,omitempty"`
	IgnoreEpisodesWithoutFiles   bool         `json:"ignoreEpisodesWithoutFiles,omitempty"`
	Monitor                      MonitorTypes `json:"monitor,omitempty"`
	SearchForCutoffUnmetEpisodes bool         `json:"searchForCutoffUnmetEpisodes,omitempty"`
	SearchForMissingEpisodes     bool         `json:"searchForMissingEpisodes,omitempty"`
}

type MonitorTypes string

const (
	MonitorTypesAll               MonitorTypes = "all"
	MonitorTypesExisting          MonitorTypes = "existing"
	MonitorTypesFirstSeason       MonitorTypes = "firstSeason"
	MonitorTypesFuture            MonitorTypes = "future"
	MonitorTypesLatestSeason      MonitorTypes = "latestSeason"
	MonitorTypesMissing           MonitorTypes = "missing"
	MonitorTypesMonitorSpecials   MonitorTypes = "monitorSpecials"
	MonitorTypesNone              MonitorTypes = "none"
	MonitorTypesPilot             MonitorTypes = "pilot"
	MonitorTypesUnknown           MonitorTypes = "unknown"
	MonitorTypesUnmonitorSpecials MonitorTypes = "unmonitorSpecials"
)

type AlternateTitleResource struct {
	Comment           string `json:"comment"`
	SceneOrigin       string `json:"sceneOrigin"`
	SceneSeasonNumber int32  `json:"sceneSeasonNumber"`
	SeasonNumber      int32  `json:"seasonNumber"`
	Title             string `json:"title"`
}

type MediaCover struct {
	CoverType MediaCoverTypes `json:"coverType,omitempty"`
	RemoteUrl string          `json:"remoteUrl"`
	Url       string          `json:"url"`
}

type MediaCoverTypes string

const (
	MediaCoverTypesBanner     MediaCoverTypes = "banner"
	MediaCoverTypesClearlogo  MediaCoverTypes = "clearlogo"
	MediaCoverTypesFanart     MediaCoverTypes = "fanart"
	MediaCoverTypesHeadshot   MediaCoverTypes = "headshot"
	MediaCoverTypesPoster     MediaCoverTypes = "poster"
	MediaCoverTypesScreenshot MediaCoverTypes = "screenshot"
	MediaCoverTypesUnknown    MediaCoverTypes = "unknown"
)

type Language struct {
	Id   int32  `json:"id,omitempty"`
	Name string `json:"name"`
}

type Ratings struct {
	Value float64 `json:"value,omitempty"`
	Votes int32   `json:"votes,omitempty"`
}

type SeasonResource struct {
	Images       []MediaCover             `json:"images"`
	Monitored    bool                     `json:"monitored,omitempty"`
	SeasonNumber int32                    `json:"seasonNumber,omitempty"`
	Statistics   SeasonStatisticsResource `json:"statistics,omitempty"`
}

type SeasonStatisticsResource struct {
	EpisodeCount      int32     `json:"episodeCount,omitempty"`
	EpisodeFileCount  int32     `json:"episodeFileCount,omitempty"`
	NextAiring        time.Time `json:"nextAiring"`
	PercentOfEpisodes float64   `json:"percentOfEpisodes,omitempty"`
	PreviousAiring    time.Time `json:"previousAiring"`
	ReleaseGroups     []string  `json:"releaseGroups"`
	SizeOnDisk        int64     `json:"sizeOnDisk,omitempty"`
	TotalEpisodeCount int32     `json:"totalEpisodeCount,omitempty"`
}

type SeriesTypes string

const (
	Anime    SeriesTypes = "anime"
	Daily    SeriesTypes = "daily"
	Standard SeriesTypes = "standard"
)

type SeriesStatisticsResource struct {
	EpisodeCount      int32    `json:"episodeCount,omitempty"`
	EpisodeFileCount  int32    `json:"episodeFileCount,omitempty"`
	PercentOfEpisodes float64  `json:"percentOfEpisodes,omitempty"`
	ReleaseGroups     []string `json:"releaseGroups"`
	SeasonCount       int32    `json:"seasonCount,omitempty"`
	SizeOnDisk        int64    `json:"sizeOnDisk,omitempty"`
	TotalEpisodeCount int32    `json:"totalEpisodeCount,omitempty"`
}

type SeriesStatusType string

const (
	Continuing SeriesStatusType = "continuing"
	Deleted    SeriesStatusType = "deleted"
	Ended      SeriesStatusType = "ended"
	Upcoming   SeriesStatusType = "upcoming"
)

type SonarrErrorMessage string

const SonarrErrorMessageSeriesAlreadyAdded SonarrErrorMessage = "This series has already been added"

type SonarrError struct {
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

type Series struct {
	AddOptions        AddSeriesOptions         `json:"addOptions,omitempty"`
	Added             time.Time                `json:"added,omitempty"`
	AirTime           string                   `json:"airTime"`
	AlternateTitles   []AlternateTitleResource `json:"alternateTitles"`
	Certification     string                   `json:"certification"`
	CleanTitle        string                   `json:"cleanTitle"`
	Ended             bool                     `json:"ended,omitempty"`
	EpisodesChanged   bool                     `json:"episodesChanged"`
	FirstAired        time.Time                `json:"firstAired"`
	Folder            string                   `json:"folder"`
	Genres            []string                 `json:"genres"`
	Id                int32                    `json:"id,omitempty"`
	Images            []MediaCover             `json:"images"`
	ImdbId            string                   `json:"imdbId"`
	Monitored         bool                     `json:"monitored,omitempty"`
	Network           string                   `json:"network"`
	NextAiring        time.Time                `json:"nextAiring"`
	OriginalLanguage  Language                 `json:"originalLanguage,omitempty"`
	Overview          string                   `json:"overview"`
	Path              string                   `json:"path"`
	PreviousAiring    time.Time                `json:"previousAiring"`
	ProfileName       string                   `json:"profileName"`
	QualityProfileId  int32                    `json:"qualityProfileId,omitempty"`
	Ratings           Ratings                  `json:"ratings,omitempty"`
	RemotePoster      string                   `json:"remotePoster"`
	RootFolderPath    string                   `json:"rootFolderPath"`
	Runtime           int32                    `json:"runtime,omitempty"`
	SeasonFolder      bool                     `json:"seasonFolder,omitempty"`
	Seasons           []SeasonResource         `json:"seasons"`
	SeriesType        SeriesTypes              `json:"seriesType,omitempty"`
	SortTitle         string                   `json:"sortTitle"`
	Statistics        SeriesStatisticsResource `json:"statistics,omitempty"`
	Status            SeriesStatusType         `json:"status,omitempty"`
	Tags              []int32                  `json:"tags"`
	Title             string                   `json:"title"`
	TitleSlug         string                   `json:"titleSlug"`
	TvMazeId          int32                    `json:"tvMazeId,omitempty"`
	TvRageId          int32                    `json:"tvRageId,omitempty"`
	TvdbId            int32                    `json:"tvdbId,omitempty"`
	UseSceneNumbering bool                     `json:"useSceneNumbering,omitempty"`
	Year              int32                    `json:"year,omitempty"`
}

type Tag struct {
	Id    int32  `json:"id"`
	Label string `json:"label"`
}

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

func (c *Client) AddSeries(title string, tvdbid int32, tags []int32) error {
	s := Series{
		Title:            title,
		QualityProfileId: c.config.Sonarr.QualityProfileID,
		SeasonFolder:     c.config.Sonarr.SeasonFolder,
		Monitored:        c.config.Sonarr.Monitored,
		TvdbId:           tvdbid,
		SeriesType:       SeriesTypes("anime"),
		RootFolderPath:   c.config.Sonarr.RootFolderPath,
		AddOptions: AddSeriesOptions{
			IgnoreEpisodesWithFiles:      false,
			IgnoreEpisodesWithoutFiles:   false,
			Monitor:                      MonitorTypes(c.config.Sonarr.MonitorType),
			SearchForMissingEpisodes:     true,
			SearchForCutoffUnmetEpisodes: false,
		},
		Tags: tags,
	}

	p, err := json.Marshal(s)
	if err != nil {
		return err
	}

	_, se, err := c.SendPostRequest(c.config.Sonarr.Url.JoinPath("/api/v3/series").String(), p)
	if err != nil {
		return err
	}

	if se != nil {
		if se.ErrorMessage == string(SonarrErrorMessageSeriesAlreadyAdded) {
			ss, err := c.GetSeries(s.TvdbId)
			if err != nil {
				return err
			}

			if ss[0].HaveTag(tags[0]) {
				return nil
			} else {
				err = c.PutSeries(s.Id, ss, tags)
				if err != nil {
					return err
				}

				return nil
			}
		}
	}

	return nil
}

func (c *Client) PutSeries(id int32, s []Series, tags []int32) error {
	u := c.config.Sonarr.Url.JoinPath(fmt.Sprintf("/api/v3/series/%v", id))
	s[0].Tags = append(s[0].Tags, tags[0])
	body, err := json.MarshalIndent(s[0], "", "  ")
	if err != nil {
		return err
	}

	_, err = c.SendPutRequest(u.String(), body)
	if err != nil {
		return err
	}

	return nil
}

func (s *Series) HaveTag(id int32) bool {
	for _, v := range s.Tags {
		if v == id {
			return true
		}
	}

	return false
}

func (c *Client) GetSeries(tvdbid int32) ([]Series, error) {
	s := []Series{}
	u := c.config.Sonarr.Url.JoinPath("/api/v3/series")
	params := u.Query()
	params.Add("tvdbId", fmt.Sprintf("%v", tvdbid))
	u.RawQuery = params.Encode()

	data, err := c.SendGetRequest(u.String())
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (c *Client) AddTag(label string) (int32, error) {
	t := Tag{
		Label: label,
	}

	p, err := json.Marshal(t)
	if err != nil {
		return -1, err
	}

	url := c.config.Sonarr.Url.JoinPath("/api/v3/tag").String()
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

func (c *Client) TagExists(label string) (bool, int32, error) {
	var tags []Tag
	url := c.config.Sonarr.Url.JoinPath("/api/v3/tag").String()
	data, err := c.SendGetRequest(url)
	if err != nil {
		return false, -1, err
	}

	err = json.Unmarshal(data, &tags)
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

func (c *Client) SendPostRequest(url string, body []byte) ([]byte, *SonarrError, error) {
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
		se := []SonarrError{}
		err = json.Unmarshal(rb, &se)
		if err != nil {
			return nil, nil, fmt.Errorf("error: %v responsePost from Sonarr:\n%v", err, string(rb))
		}

		return nil, &se[0], nil
	}

	return rb, nil, nil
}

func (c *Client) SendGetRequest(url string) ([]byte, error) {
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	rb, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("responseGet from Sonarr:\n%v", string(rb))
	}

	//fmt.Printf("Get got this:\n%v", string(rb))
	return rb, nil
}

func (c *Client) SendPutRequest(url string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	rb, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("\nrequest:\n%v\nresponsePut from Sonarr:\n%v", string(body), string(rb))
	}

	return rb, nil
}
