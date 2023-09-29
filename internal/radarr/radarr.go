package radarr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/varoOP/malbrr/internal/config"
)

type MovieTranslation struct {
	CleanTitle      string   `json:"cleanTitle"`
	Id              int32    `json:"id,omitempty"`
	Language        Language `json:"language,omitempty"`
	MovieMetadataId int32    `json:"movieMetadataId,omitempty"`
	Overview        string   `json:"overview"`
	Title           string   `json:"title"`
}

type RatingType string

const (
	RatingTypeCritic RatingType = "critic"
	RatingTypeUser   RatingType = "user"
)

type RatingChild struct {
	Type  RatingType `json:"type,omitempty"`
	Value float64    `json:"value,omitempty"`
	Votes int32      `json:"votes,omitempty"`
}

type Ratings struct {
	Imdb           RatingChild `json:"imdb,omitempty"`
	Metacritic     RatingChild `json:"metacritic,omitempty"`
	RottenTomatoes RatingChild `json:"rottenTomatoes,omitempty"`
	Tmdb           RatingChild `json:"tmdb,omitempty"`
}

type Language struct {
	Id   int32  `json:"id,omitempty"`
	Name string `json:"name"`
}

type AlternativeTitle struct {
	CleanTitle      string     `json:"cleanTitle"`
	Id              int32      `json:"id,omitempty"`
	MovieMetadataId int32      `json:"movieMetadataId,omitempty"`
	SourceType      SourceType `json:"sourceType,omitempty"`
	Title           string     `json:"title"`
}

type MovieMetadata struct {
	AlternativeTitles  []AlternativeTitle `json:"alternativeTitles"`
	Certification      string             `json:"certification"`
	CleanOriginalTitle string             `json:"cleanOriginalTitle"`
	CleanTitle         string             `json:"cleanTitle"`
	CollectionTitle    string             `json:"collectionTitle"`
	CollectionTmdbId   int32              `json:"collectionTmdbId,omitempty"`
	DigitalRelease     time.Time          `json:"digitalRelease"`
	Genres             []string           `json:"genres"`
	Id                 int32              `json:"id,omitempty"`
	Images             []MediaCover       `json:"images"`
	ImdbId             string             `json:"imdbId"`
	InCinemas          time.Time          `json:"inCinemas"`
	IsRecentMovie      bool               `json:"isRecentMovie,omitempty"`
	LastInfoSync       time.Time          `json:"lastInfoSync"`
	OriginalLanguage   Language           `json:"originalLanguage,omitempty"`
	OriginalTitle      string             `json:"originalTitle"`
	Overview           string             `json:"overview"`
	PhysicalRelease    time.Time          `json:"physicalRelease"`
	Popularity         float32            `json:"popularity,omitempty"`
	Ratings            Ratings            `json:"ratings,omitempty"`
	Recommendations    []int32            `json:"recommendations"`
	Runtime            int32              `json:"runtime,omitempty"`
	SecondaryYear      int32              `json:"secondaryYear"`
	SortTitle          string             `json:"sortTitle"`
	Status             MovieStatusType    `json:"status,omitempty"`
	Studio             string             `json:"studio"`
	Title              string             `json:"title"`
	TmdbId             int32              `json:"tmdbId,omitempty"`
	Translations       []MovieTranslation `json:"translations"`
	Website            string             `json:"website"`
	Year               int32              `json:"year,omitempty"`
	YouTubeTrailerId   string             `json:"youTubeTrailerId"`
}

type MovieStatusType string

const (
	Announced MovieStatusType = "announced"
	Deleted   MovieStatusType = "deleted"
	InCinemas MovieStatusType = "inCinemas"
	Released  MovieStatusType = "released"
	Tba       MovieStatusType = "tba"
)

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

type MediaCover struct {
	CoverType MediaCoverTypes `json:"coverType,omitempty"`
	RemoteUrl string          `json:"remoteUrl"`
	Url       string          `json:"url"`
}

type MovieCollection struct {
	Added               time.Time       `json:"added,omitempty"`
	CleanTitle          string          `json:"cleanTitle"`
	Id                  int32           `json:"id,omitempty"`
	Images              []MediaCover    `json:"images"`
	LastInfoSync        time.Time       `json:"lastInfoSync"`
	MinimumAvailability MovieStatusType `json:"minimumAvailability,omitempty"`
	Monitored           bool            `json:"monitored,omitempty"`
	Movies              []MovieMetadata `json:"movies"`
	Overview            string          `json:"overview"`
	QualityProfileId    int32           `json:"qualityProfileId,omitempty"`
	RootFolderPath      string          `json:"rootFolderPath"`
	SearchOnAdd         bool            `json:"searchOnAdd,omitempty"`
	SortTitle           string          `json:"sortTitle"`
	Tags                []int32         `json:"tags"`
	Title               string          `json:"title"`
	TmdbId              int32           `json:"tmdbId,omitempty"`
}

type SourceType string

const (
	SourceTypeIndexer  SourceType = "indexer"
	SourceTypeMappings SourceType = "mappings"
	SourceTypeTmdb     SourceType = "tmdb"
	SourceTypeUser     SourceType = "user"
)

type MonitorTypes string

const (
	MonitorTypesMovieAndCollection MonitorTypes = "movieAndCollection"
	MonitorTypesMovieOnly          MonitorTypes = "movieOnly"
	MonitorTypesNone               MonitorTypes = "none"
)

type AlternativeTitleResource struct {
	CleanTitle      string     `json:"cleanTitle"`
	Id              int32      `json:"id,omitempty"`
	MovieMetadataId int32      `json:"movieMetadataId,omitempty"`
	SourceType      SourceType `json:"sourceType,omitempty"`
	Title           string     `json:"title"`
}

type AddMovieMethod string

const (
	AddMovieMethodCollection AddMovieMethod = "collection"
	AddMovieMethodList       AddMovieMethod = "list"
	AddMovieMethodManual     AddMovieMethod = "manual"
)

type AddMovieOptions struct {
	AddMethod                  AddMovieMethod `json:"addMethod,omitempty"`
	IgnoreEpisodesWithFiles    bool           `json:"ignoreEpisodesWithFiles,omitempty"`
	IgnoreEpisodesWithoutFiles bool           `json:"ignoreEpisodesWithoutFiles,omitempty"`
	Monitor                    MonitorTypes   `json:"monitor,omitempty"`
	SearchForMovie             bool           `json:"searchForMovie,omitempty"`
}

type PrivacyLevel string

const (
	PrivacyLevelApiKey   PrivacyLevel = "apiKey"
	PrivacyLevelNormal   PrivacyLevel = "normal"
	PrivacyLevelPassword PrivacyLevel = "password"
	PrivacyLevelUserName PrivacyLevel = "userName"
)

type SelectOption struct {
	DividerAfter bool   `json:"dividerAfter,omitempty"`
	Hint         string `json:"hint"`
	Name         string `json:"name"`
	Order        int32  `json:"order,omitempty"`
	Value        int32  `json:"value,omitempty"`
}

type Field struct {
	Advanced                    bool           `json:"advanced,omitempty"`
	HelpLink                    string         `json:"helpLink"`
	HelpText                    string         `json:"helpText"`
	HelpTextWarning             string         `json:"helpTextWarning"`
	Hidden                      string         `json:"hidden"`
	IsFloat                     bool           `json:"isFloat,omitempty"`
	Label                       string         `json:"label"`
	Name                        string         `json:"name"`
	Order                       int32          `json:"order,omitempty"`
	Placeholder                 string         `json:"placeholder"`
	Privacy                     PrivacyLevel   `json:"privacy,omitempty"`
	Section                     string         `json:"section"`
	SelectOptions               []SelectOption `json:"selectOptions"`
	SelectOptionsProviderAction string         `json:"selectOptionsProviderAction"`
	Type                        string         `json:"type"`
	Unit                        string         `json:"unit"`
	Value                       interface{}    `json:"value"`
}

type CustomFormatSpecificationSchema struct {
	Fields             []Field                           `json:"fields"`
	Id                 int32                             `json:"id,omitempty"`
	Implementation     string                            `json:"implementation"`
	ImplementationName string                            `json:"implementationName"`
	InfoLink           string                            `json:"infoLink"`
	Name               string                            `json:"name"`
	Negate             bool                              `json:"negate,omitempty"`
	Presets            []CustomFormatSpecificationSchema `json:"presets"`
	Required           bool                              `json:"required,omitempty"`
}

type CustomFormatResource struct {
	Id                              int32                             `json:"id,omitempty"`
	IncludeCustomFormatWhenRenaming bool                              `json:"includeCustomFormatWhenRenaming"`
	Name                            string                            `json:"name"`
	Specifications                  []CustomFormatSpecificationSchema `json:"specifications"`
}

type MediaInfoResource struct {
	AudioBitrate          int64   `json:"audioBitrate,omitempty"`
	AudioChannels         float64 `json:"audioChannels,omitempty"`
	AudioCodec            string  `json:"audioCodec"`
	AudioLanguages        string  `json:"audioLanguages"`
	AudioStreamCount      int32   `json:"audioStreamCount,omitempty"`
	Id                    int32   `json:"id,omitempty"`
	Resolution            string  `json:"resolution"`
	RunTime               string  `json:"runTime"`
	ScanType              string  `json:"scanType"`
	Subtitles             string  `json:"subtitles"`
	VideoBitDepth         int32   `json:"videoBitDepth,omitempty"`
	VideoBitrate          int64   `json:"videoBitrate,omitempty"`
	VideoCodec            string  `json:"videoCodec"`
	VideoDynamicRange     string  `json:"videoDynamicRange"`
	VideoDynamicRangeType string  `json:"videoDynamicRangeType"`
	VideoFps              float64 `json:"videoFps,omitempty"`
}

type Modifier string

const (
	ModifierBrdisk   Modifier = "brdisk"
	ModifierNone     Modifier = "none"
	ModifierRawhd    Modifier = "rawhd"
	ModifierRegional Modifier = "regional"
	ModifierRemux    Modifier = "remux"
	ModifierScreener Modifier = "screener"
)

type QualitySource string

const (
	Bluray    QualitySource = "bluray"
	Cam       QualitySource = "cam"
	Dvd       QualitySource = "dvd"
	Telecine  QualitySource = "telecine"
	Telesync  QualitySource = "telesync"
	Tv        QualitySource = "tv"
	Unknown   QualitySource = "unknown"
	Webdl     QualitySource = "webdl"
	Webrip    QualitySource = "webrip"
	Workprint QualitySource = "workprint"
)

type Quality struct {
	Id         int32         `json:"id,omitempty"`
	Modifier   Modifier      `json:"modifier,omitempty"`
	Name       string        `json:"name"`
	Resolution int32         `json:"resolution,omitempty"`
	Source     QualitySource `json:"source,omitempty"`
}

type Revision struct {
	IsRepack bool  `json:"isRepack,omitempty"`
	Real     int32 `json:"real,omitempty"`
	Version  int32 `json:"version,omitempty"`
}

type QualityModel struct {
	Quality  Quality  `json:"quality,omitempty"`
	Revision Revision `json:"revision,omitempty"`
}

type MovieFileResource struct {
	CustomFormatScore   int32                  `json:"customFormatScore,omitempty"`
	CustomFormats       []CustomFormatResource `json:"customFormats"`
	DateAdded           time.Time              `json:"dateAdded,omitempty"`
	Edition             string                 `json:"edition"`
	Id                  int32                  `json:"id,omitempty"`
	IndexerFlags        int32                  `json:"indexerFlags,omitempty"`
	Languages           []Language             `json:"languages"`
	MediaInfo           MediaInfoResource      `json:"mediaInfo,omitempty"`
	MovieId             int32                  `json:"movieId,omitempty"`
	OriginalFilePath    string                 `json:"originalFilePath"`
	Path                string                 `json:"path"`
	Quality             QualityModel           `json:"quality,omitempty"`
	QualityCutoffNotMet bool                   `json:"qualityCutoffNotMet,omitempty"`
	RelativePath        string                 `json:"relativePath"`
	ReleaseGroup        string                 `json:"releaseGroup"`
	SceneName           string                 `json:"sceneName"`
	Size                int64                  `json:"size,omitempty"`
}

type Movie struct {
	AddOptions            AddMovieOptions            `json:"addOptions,omitempty"`
	Added                 time.Time                  `json:"added,omitempty"`
	AlternateTitles       []AlternativeTitleResource `json:"alternateTitles"`
	Certification         string                     `json:"certification"`
	CleanTitle            string                     `json:"cleanTitle"`
	Collection            MovieCollection            `json:"collection,omitempty"`
	DigitalRelease        time.Time                  `json:"digitalRelease"`
	Folder                string                     `json:"folder"`
	FolderName            string                     `json:"folderName"`
	Genres                []string                   `json:"genres"`
	HasFile               bool                       `json:"hasFile,omitempty"`
	Id                    int32                      `json:"id,omitempty"`
	Images                []MediaCover               `json:"images"`
	ImdbId                string                     `json:"imdbId"`
	InCinemas             time.Time                  `json:"inCinemas"`
	IsAvailable           bool                       `json:"isAvailable,omitempty"`
	MinimumAvailability   MovieStatusType            `json:"minimumAvailability,omitempty"`
	Monitored             bool                       `json:"monitored,omitempty"`
	MovieFile             MovieFileResource          `json:"movieFile,omitempty"`
	OriginalLanguage      Language                   `json:"originalLanguage,omitempty"`
	OriginalTitle         string                     `json:"originalTitle"`
	Overview              string                     `json:"overview"`
	Path                  string                     `json:"path"`
	PhysicalRelease       time.Time                  `json:"physicalRelease"`
	PhysicalReleaseNote   string                     `json:"physicalReleaseNote"`
	Popularity            float32                    `json:"popularity,omitempty"`
	QualityProfileId      int32                      `json:"qualityProfileId,omitempty"`
	Ratings               Ratings                    `json:"ratings,omitempty"`
	RemotePoster          string                     `json:"remotePoster"`
	RootFolderPath        string                     `json:"rootFolderPath"`
	Runtime               int32                      `json:"runtime,omitempty"`
	SecondaryYear         int32                      `json:"secondaryYear"`
	SecondaryYearSourceId int32                      `json:"secondaryYearSourceId,omitempty"`
	SizeOnDisk            int64                      `json:"sizeOnDisk"`
	SortTitle             string                     `json:"sortTitle"`
	Status                MovieStatusType            `json:"status,omitempty"`
	Studio                string                     `json:"studio"`
	Tags                  []int32                    `json:"tags"`
	Title                 string                     `json:"title"`
	TitleSlug             string                     `json:"titleSlug"`
	TmdbId                int32                      `json:"tmdbId,omitempty"`
	Website               string                     `json:"website"`
	Year                  int32                      `json:"year,omitempty"`
	YouTubeTrailerId      string                     `json:"youTubeTrailerId"`
}

type Tag struct {
	Id    int32    `json:"id"`
	Label string `json:"label"`
}

type RadarrErrorMessage string

const ErrorMessageMovieAlreadyAdded RadarrErrorMessage = "This movie has already been added"

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

type Client struct {
	client *http.Client
	config *config.Config
}

func NewClient(cfg *config.Config) *Client {
	c := &http.Client{
		Transport: &config.ApiKeyTransport{ApiKey: cfg.Radarr.ApiKey},
	}

	return &Client{
		config: cfg,
		client: c,
	}
}

func (c *Client) AddMovie(title string, tmdbid int32, tags []int32) error {
	m := Movie{
		Title:               title,
		MinimumAvailability: MovieStatusType(c.config.Radarr.MinimumAvailability),
		QualityProfileId:    c.config.Radarr.QualityProfileID,
		Monitored:           c.config.Radarr.Monitored,
		TmdbId:              tmdbid,
		RootFolderPath:      c.config.Radarr.RootFolderPath,
		AddOptions: AddMovieOptions{
			IgnoreEpisodesWithFiles:    false,
			IgnoreEpisodesWithoutFiles: false,
			Monitor:                    MonitorTypes(c.config.Radarr.MonitorType),
			SearchForMovie:             c.config.Radarr.SearchForMovie,
			AddMethod:                  AddMovieMethodManual,
		},
		Tags: tags,
	}

	p, err := json.Marshal(m)
	if err != nil {
		return err
	}

	_, me, err := c.SendPostRequest(c.config.Radarr.Url.JoinPath("/api/v3/movie").String(), p)
	if err != nil {
		return err
	}

	if me != nil {
		if me.ErrorMessage == string(ErrorMessageMovieAlreadyAdded) {
			mm, err := c.GetMovie(m.TmdbId)
			if err != nil {
				return err
			}

			if mm[0].HaveTag(tags[0]) {
				return nil
			} else {
				err = c.PutSeries(m.Id, mm, tags)
				if err != nil {
					return err
				}

				return nil
			}
		}
	}

	return nil
}

func (c *Client) PutSeries(id int32, m []Movie, tags []int32) error {
	u := c.config.Radarr.Url.JoinPath(fmt.Sprintf("/api/v3/movie/%v", id))
	m[0].Tags = append(m[0].Tags, tags[0])
	body, err := json.MarshalIndent(m[0], "", "  ")
	if err != nil {
		return err
	}

	_, err = c.SendPutRequest(u.String(), body)
	if err != nil {
		return err
	}

	return nil
}

func (m *Movie) HaveTag(id int32) bool {
	for _, v := range m.Tags {
		if v == id {
			return true
		}
	}

	return false
}

func (c *Client) GetMovie(tmdbid int32) ([]Movie, error) {
	m := []Movie{}
	u := c.config.Radarr.Url.JoinPath("/api/v3/movie")
	params := u.Query()
	params.Add("tmdbId", fmt.Sprintf("%v", tmdbid))
	u.RawQuery = params.Encode()

	data, err := c.SendGetRequest(u.String())
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (c *Client) AddTag(label string) (int32, error) {
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

func (c *Client) TagExists(label string) (bool, int32, error) {
	var tags []Tag
	url := c.config.Radarr.Url.JoinPath("/api/v3/tag").String()
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
		me := []RadarrError{}
		err = json.Unmarshal(rb, &me)
		if err != nil {
			return nil, nil, fmt.Errorf("error: %v responsePost from Radarr:\n%v", err, string(rb))
		}

		return nil, &me[0], nil
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
		return nil, fmt.Errorf("responseGet from Radarr:\n%v", string(rb))
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
		return nil, fmt.Errorf("\nrequest:\n%v\nresponsePut from Radarr:\n%v", string(body), string(rb))
	}

	return rb, nil
}
