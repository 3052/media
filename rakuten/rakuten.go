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
