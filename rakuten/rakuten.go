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

// github.com/pandvan/rakuten-m3u-generator/blob/master/rakuten.py
var classificationMap = map[string]int{
   "cz": 272,
   "dk": 283,
   "es": 5,
   "fr": 23,
   "ie": 41,
   "nl": 69,
   "pl": 277,
   "pt": 64,
   "se": 282,
   "uk": 18,
}

// join takes a variable number of strings and returns them combined into one string without separators.
func join(strs ...string) string {
   return strings.Join(strs, "")
}

type AudioLanguage struct {
   Id string `json:"id"`
}

type SeasonData struct {
   Episodes []VideoItem `json:"episodes"`
}

type Stream struct {
   AudioLanguages []AudioLanguage `json:"audio_languages"`
}

type StreamData struct {
   StreamInfos []struct {
      LicenseUrl string `json:"license_url"`
      Url        string `json:"url"`
   } `json:"stream_infos"`
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

type VideoItem struct {
   Title       string      `json:"title"`
   Id          string      `json:"id"`
   ViewOptions ViewOptions `json:"view_options"`
}

type ViewOptions struct {
   Private struct {
      Streams []Stream `json:"streams"`
   } `json:"private"`
}

type Dash struct {
   Body []byte
   Url  *url.URL
}

// RequestMovie fetches movie details (GET).
func (m *Media) RequestMovie() (*VideoItem, error) {
   if m.Type != Movies {
      return nil, errors.New("cannot request movie details for a non-movie content type")
   }
   fullURL, err := buildURL(m.MarketCode, "movies", m.Id)
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
   fullURL, err := buildURL(m.MarketCode, "tv_shows", m.Id)
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
   fullURL, err := buildURL(m.MarketCode, "seasons", seasonId)
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

func (s StreamData) Dash() (*Dash, error) {
   resp, err := http.Get(s.StreamInfos[0].Url)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Dash
   result.Body, err = io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   result.Url = resp.Request.URL
   return &result, nil
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

func buildURL(marketCode, endpoint, id string) (string, error) {
   classID, ok := classificationMap[marketCode]
   if !ok {
      return "", fmt.Errorf("unsupported market code %v", marketCode)
   }
   params := url.Values{}
   params.Add("device_identifier", DeviceID)
   params.Add("market_code", marketCode)
   params.Add("classification_id", strconv.Itoa(classID))
   return join(
      "https://gizmo.rakuten.tv/v3/",
      endpoint,
      "/",
      id,
      "?",
      params.Encode(),
   ), nil
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

// DeviceID is the default identifier used for requests.
const DeviceID = "atvui40"
