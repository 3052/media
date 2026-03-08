package rakuten

import (
   "bytes"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

func makeStreamRequest(marketCode, contentType, contentId string, player PlayerType, quality VideoQuality, audioLanguage string) (*StreamData, error) {
   classId, ok := classificationMap[marketCode]
   if !ok {
      return nil, fmt.Errorf("unsupported market code: %s", marketCode)
   }
   data, err := json.Marshal(map[string]string{
      "audio_language":              audioLanguage,
      "audio_quality":               "2.0",
      "classification_id":           strconv.Itoa(classId),
      "content_id":                  contentId,
      "content_type":                contentType,
      "device_identifier":           DeviceId,
      "device_serial":               "not implemented",
      "device_stream_video_quality": string(quality),
      "player":                      string(player),
      "subtitle_language":           "MIS",
      "video_type":                  "stream",
   })
   if err != nil {
      return nil, err
   }
   resp, err := http.Post(
      "https://gizmo.rakuten.tv/v3/avod/streamings", "application/json",
      bytes.NewBuffer(data),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var wrapper struct {
      Data   StreamData `json:"data"`
      Errors []struct {
         Message string
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&wrapper)
   if err != nil {
      return nil, err
   }
   if len(wrapper.Errors) >= 1 {
      return nil, errors.New(wrapper.Errors[0].Message)
   }
   return &wrapper.Data, nil
}
func buildUrl(marketCode, endpoint, id string) (string, error) {
   classId, ok := classificationMap[marketCode]
   if !ok {
      return "", fmt.Errorf("unsupported market code %v", marketCode)
   }
   url_data := url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   join("/v3/", endpoint, "/", id),
      RawQuery: url.Values{
         "classification_id": {strconv.Itoa(classId)},
         "device_identifier": {DeviceId},
         "market_code":       {marketCode},
      }.Encode(),
   }
   return url_data.String(), nil
}

func (s StreamData) Dash() (*Dash, error) {
   resp, err := http.Get(s.StreamInfos[0].Url)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   body, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   return &Dash{Body: body, Url: resp.Request.URL}, nil
}

// EpisodeStream requests the stream for a specific TV Show Episode (POST).
func (m *Media) EpisodeStream(episodeId, audioLanguage string, player PlayerType, quality VideoQuality) (*StreamData, error) {
   if m.Type != TvShows {
      return nil, errors.New("cannot request an episode stream for non-tv-show content")
   }
   return makeStreamRequest(m.MarketCode, "episodes", episodeId, player, quality, audioLanguage)
}

// RequestMovie fetches movie details (GET).
func (m *Media) RequestMovie() (*VideoItem, error) {
   if m.Type != Movies {
      return nil, errors.New("cannot request movie details for a non-movie content type")
   }
   fullURL, err := buildUrl(m.MarketCode, "movies", m.Id)
   if err != nil {
      return nil, err
   }
   resp, err := http.Get(fullURL)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var wrapper struct {
      Data VideoItem `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
      return nil, err
   }
   return &wrapper.Data, nil
}

// RequestTvShow fetches TV show details like seasons (GET).
func (m *Media) RequestTvShow() (*TvShowData, error) {
   if m.Type != TvShows {
      return nil, errors.New("cannot request tv show details for a non-tv show content type")
   }
   fullURL, err := buildUrl(m.MarketCode, "tv_shows", m.Id)
   if err != nil {
      return nil, err
   }
   resp, err := http.Get(fullURL)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var wrapper struct {
      Data TvShowData `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
      return nil, err
   }
   return &wrapper.Data, nil
}

// RequestSeason fetches episodes for a specific season (GET).
// This method is only applicable to TV Shows.
func (m *Media) RequestSeason(seasonId string) (*SeasonData, error) {
   if m.Type != TvShows {
      return nil, errors.New("cannot request season for a non-tv show content type")
   }
   fullURL, err := buildUrl(m.MarketCode, "seasons", seasonId)
   if err != nil {
      return nil, err
   }
   resp, err := http.Get(fullURL)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var wrapper struct {
      Data SeasonData `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
      return nil, err
   }
   return &wrapper.Data, nil
}

func (s StreamData) Widevine(data []byte) ([]byte, error) {
   resp, err := http.Post(
      s.StreamInfos[0].LicenseUrl, "application/x-protobuf",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

// MovieStream requests the stream for this movie (POST).
func (m *Media) MovieStream(audioLanguage string, player PlayerType, quality VideoQuality) (*StreamData, error) {
   if m.Type != Movies {
      return nil, errors.New("cannot request a movie stream for non-movie content")
   }
   return makeStreamRequest(m.MarketCode, "movies", m.Id, player, quality, audioLanguage)
}

// join takes a variable number of strings and returns them combined into one string without separators.
func join(strs ...string) string {
   return strings.Join(strs, "")
}

func (t TvShowData) String() string {
   var data strings.Builder
   for i, season := range t.Seasons {
      if i >= 1 {
         data.WriteByte('\n')
      }
      data.WriteString("season id = ")
      data.WriteString(season.Id)
   }
   return data.String()
}

type TvShowData struct {
   Seasons []struct {
      Id string `json:"id"`
   } `json:"seasons"`
}

// It returns the ID, Title, and a unique list of available audio languages.
func (v *VideoItem) String() string {
   seen := make(map[string]bool)
   var data strings.Builder
   data.WriteString("title = ")
   data.WriteString(v.Title)
   data.WriteString("\nid = ")
   data.WriteString(v.Id)
   for _, stream := range v.ViewOptions.Private.Streams {
      for _, lang := range stream.AudioLanguages {
         if !seen[lang.Id] {
            seen[lang.Id] = true
            data.WriteString("\naudio language = ")
            data.WriteString(lang.Id)
         }
      }
   }
   return data.String()
}

type Dash struct {
   Body []byte
   Url  *url.URL
}

// extractMarketCode extracts the first segment of the path (e.g., "nl", "uk").
func extractMarketCode(path string) (string, error) {
   trimmed := strings.Trim(path, "/")
   // Check if we have anything left after trimming
   if trimmed == "" {
      return "", errors.New("could not determine market code from path")
   }
   segments := strings.Split(trimmed, "/")
   return segments[0], nil
}

// github.com/pandvan/rakuten-m3u-generator/blob/master/rakuten.py
var classificationMap = map[string]int{
   "cz": 272,
   "es": 5,
   "fr": 23,
   "ie": 41,
   "nl": 69,
   "pl": 277,
   "pt": 64,
   "se": 282,
   "uk": 18,
}

type VideoItem struct {
   Title       string `json:"title"`
   Id          string `json:"id"`
   ViewOptions struct {
      Private struct {
         Streams []struct {
            AudioLanguages []struct {
               Id string `json:"id"`
            } `json:"audio_languages"`
         } `json:"streams"`
      } `json:"private"`
   } `json:"view_options"`
}

type SeasonData struct {
   Episodes []VideoItem `json:"episodes"`
}

type StreamData struct {
   StreamInfos []struct {
      LicenseUrl string `json:"license_url"`
      Url        string `json:"url"`
   } `json:"stream_infos"`
}

const DeviceId = "atvui40"

type VideoQuality string

const (
   Fhd VideoQuality = "FHD"
   Hd  VideoQuality = "HD"
)

type PlayerType string

const (
   PlayReady PlayerType = DeviceId + ":DASH-CENC:PR"
   Widevine  PlayerType = DeviceId + ":DASH-CENC:WVM"
)

///

type Content int

const (
   Movies Content = iota
   TvShows
)

type Media struct {
   Id         string
   MarketCode string
   Type       Content
}

// Parse populates the Media struct from a raw URL.
func (m *Media) Parse(rawLink string) error {
   link, err := url.Parse(rawLink)
   if err != nil {
      return err
   }
   marketCode, err := extractMarketCode(link.Path)
   if err != nil {
      return err
   }
   m.MarketCode = marketCode

   // 1. Check Query Parameters
   query := link.Query()
   // 'contentType' here is the URL parameter value (e.g. "movies", "tv_shows")
   contentType := query.Get("content_type")
   if contentType == "movies" || contentType == "tv_shows" {
      var id string
      if contentType == "movies" {
         id = query.Get("content_id")
         if id == "" {
            return errors.New("url missing content_id param")
         }
         m.Type = Movies
      } else {
         id = query.Get("tv_show_id")
         if id == "" {
            return errors.New("url missing tv_show_id param")
         }
         m.Type = TvShows
      }
      m.Id = id
      return nil
   }

   // 2. Check Path Segments
   path := strings.Trim(link.Path, "/")
   segments := strings.Split(path, "/")
   for _, seg := range segments {
      if seg == "movies" || seg == "tv_shows" {
         id := segments[len(segments)-1]
         if id == seg {
            return fmt.Errorf("url does not contain a specific %s id", seg)
         }
         m.Id = id
         if seg == "movies" {
            m.Type = Movies
         } else {
            m.Type = TvShows
         }
         return nil
      }
   }
   return errors.New("not a movie or tv show url")
}
