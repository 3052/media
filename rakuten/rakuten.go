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

type SeasonData struct {
   Episodes []VideoItem `json:"episodes"`
}

type StreamData struct {
   StreamInfos []struct {
      LicenseUrl string `json:"license_url"`
      Url        string `json:"url"`
   } `json:"stream_infos"`
}

// It returns the ID, Title, and a unique list of available audio languages
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

type Dash struct {
   Body []byte
   Url  *url.URL
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

// join takes a variable number of strings and returns them combined into one
// string without separators
func join(strs ...string) string {
   return strings.Join(strs, "")
}

type Content struct {
   Id         string
   MarketCode string
   Type       string
}

func (c *Content) Parse(urlData string) error {
   urlParse, err := url.Parse(urlData)
   if err != nil {
      return err
   }
   // Trim prefix once and extract the market code
   path := strings.TrimPrefix(urlParse.Path, "/")
   c.MarketCode, _, _ = strings.Cut(path, "/")
   // 1. Check Query Parameters
   query := urlParse.Query()
   contentType := query.Get("content_type")
   switch contentType {
   case "movies":
      c.Id = query.Get("content_id")
      c.Type = contentType
      return nil
   case "tv_shows":
      c.Id = query.Get("tv_show_id")
      c.Type = contentType
      return nil
   }
   // 2. Check Path Segments
   segments := strings.Split(path, "/")
   for _, segment := range segments {
      switch segment {
      case "movies", "tv_shows":
         c.Id = segments[len(segments)-1]
         c.Type = segment
         return nil
      }
   }
   return errors.New("not a movie or tv show url")
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

///

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
   var result struct {
      Data   StreamData `json:"data"`
      Errors []struct {
         Message string
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Errors) >= 1 {
      return nil, errors.New(result.Errors[0].Message)
   }
   return &result.Data, nil
}

// EpisodeStream requests the stream for a specific TV Show Episode (POST).
func (c *Content) EpisodeStream(episodeId, audioLanguage string, player PlayerType, quality VideoQuality) (*StreamData, error) {
   if c.Type != "tv_shows" {
      return nil, errors.New("cannot request an episode stream for non-tv-show content")
   }
   return makeStreamRequest(c.MarketCode, "episodes", episodeId, player, quality, audioLanguage)
}

// RequestMovie fetches movie details (GET).
func (c *Content) RequestMovie() (*VideoItem, error) {
   if c.Type != "movies" {
      return nil, errors.New("cannot request movie details for a non-movie content type")
   }
   fullUrl, err := buildUrl(c.MarketCode, "movies", c.Id)
   if err != nil {
      return nil, err
   }
   resp, err := http.Get(fullUrl)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var result struct {
      Data VideoItem `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }
   return &result.Data, nil
}

// RequestSeason fetches episodes for a specific season (GET).
// This method is only applicable to TV Shows.
func (c *Content) RequestSeason(seasonId string) (*SeasonData, error) {
   if c.Type != "tv_shows" {
      return nil, errors.New("cannot request season for a non-tv show content type")
   }
   fullUrl, err := buildUrl(c.MarketCode, "seasons", seasonId)
   if err != nil {
      return nil, err
   }
   resp, err := http.Get(fullUrl)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var result struct {
      Data SeasonData `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }
   return &result.Data, nil
}

// RequestTvShow fetches TV show details like seasons (GET).
func (c *Content) RequestTvShow() (*TvShowData, error) {
   if c.Type != "tv_shows" {
      return nil, errors.New("cannot request tv show details for a non-tv show content type")
   }
   fullUrl, err := buildUrl(c.MarketCode, "tv_shows", c.Id)
   if err != nil {
      return nil, err
   }
   resp, err := http.Get(fullUrl)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var result struct {
      Data TvShowData `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }
   return &result.Data, nil
}

// MovieStream requests the stream for this movie (POST).
func (c *Content) MovieStream(audioLanguage string, player PlayerType, quality VideoQuality) (*StreamData, error) {
   if c.Type != "movies" {
      return nil, errors.New("cannot request a movie stream for non-movie content")
   }
   return makeStreamRequest(c.MarketCode, "movies", c.Id, player, quality, audioLanguage)
}
