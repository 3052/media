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

func (s *StreamData) Dash() (*Dash, error) {
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

type Dash struct {
   Body []byte
   Url  *url.URL
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

func (m *Movie) ParseURL(rawLink string) error {
   link, err := url.Parse(rawLink)
   if err != nil {
      return err
   }
   marketCode, err := extractMarketCode(link.Path)
   if err != nil {
      return err
   }
   // 1. Check Query Parameters
   // Pattern: ?content_type=movies&content_id=...
   query := link.Query()
   if query.Get("content_type") == "movies" {
      id := query.Get("content_id")
      if id == "" {
         return errors.New("url missing content_id param")
      }

      m.MarketCode = marketCode
      m.Id = id
      return nil
   }

   // 2. Check Path Segments
   // Pattern: /nl/movies/id OR /nl/player/movies/stream/id
   path := strings.Trim(link.Path, "/")
   segments := strings.Split(path, "/")

   for _, seg := range segments {
      if seg == "movies" {
         // Assuming the ID is the last segment
         id := segments[len(segments)-1]
         if id == "movies" {
            return errors.New("url does not contain a specific movie id")
         }

         m.MarketCode = marketCode
         m.Id = id
         return nil
      }
   }
   return errors.New("not a movie url")
}

// ParseURL attempts to parse a URL and populate the TvShow struct.
// Usage:
//
//   var t TvShow
//   err := t.ParseURL("https://rakuten.tv/...")
func (t *TvShow) ParseURL(rawLink string) error {
   link, err := url.Parse(rawLink)
   if err != nil {
      return err
   }
   marketCode, err := extractMarketCode(link.Path)
   if err != nil {
      return err
   }
   // 1. Check Query Parameters
   // Pattern: ?content_type=tv_shows&tv_show_id=...
   query := link.Query()
   if query.Get("content_type") == "tv_shows" {
      id := query.Get("tv_show_id")
      if id == "" {
         return errors.New("url missing tv_show_id param")
      }
      t.MarketCode = marketCode
      t.Id = id
      return nil
   }
   // 2. Check Path Segments
   // Pattern: /nl/tv_shows/id
   path := strings.Trim(link.Path, "/")
   segments := strings.Split(path, "/")
   for _, seg := range segments {
      if seg == "tv_shows" {
         // Assuming the ID is the last segment
         id := segments[len(segments)-1]
         if id == "tv_shows" {
            return errors.New("url does not contain a specific tv show id")
         }
         t.MarketCode = marketCode
         t.Id = id
         return nil
      }
   }
   return errors.New("not a tv show url")
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

type AudioLanguage struct {
   Id string `json:"id"`
}

type Stream struct {
   AudioLanguages []AudioLanguage `json:"audio_languages"`
}

func (t TvShowData) String() string {
   var data strings.Builder
   for i, season := range t.Seasons {
      if i >= 1 {
         data.WriteByte('\n')
      }
      data.WriteString("id = ")
      data.WriteString(season.Id)
   }
   return data.String()
}

type TvShowData struct {
   Seasons []struct {
      Id string `json:"id"`
   } `json:"seasons"`
}

type StreamData struct {
   StreamInfos []struct {
      LicenseUrl string `json:"license_url"`
      Url        string `json:"url"`
   } `json:"stream_infos"`
}

// DeviceID is the default identifier used for requests.
const DeviceID = "atvui40"

// VideoQuality defines the allowed video qualities for streaming.
type VideoQuality string

var Quality = struct {
   FHD VideoQuality
   HD  VideoQuality
}{
   FHD: "FHD",
   HD:  "HD",
}

// PlayerType defines the allowed player types/DRM schemes.
type PlayerType string

var Player = struct {
   PlayReady PlayerType
   Widevine  PlayerType
}{
   PlayReady: DeviceID + ":DASH-CENC:PR",
   Widevine:  DeviceID + ":DASH-CENC:WVM",
}

// --- Shared Structs for Nested JSON ---

// --- Response Data Structs ---

type SeasonData struct {
   Episodes []VideoItem `json:"episodes"`
}

type StreamRequestPayload struct {
   AudioQuality             string       `json:"audio_quality"`
   DeviceIdentifier         string       `json:"device_identifier"`
   DeviceSerial             string       `json:"device_serial"`
   SubtitleLanguage         string       `json:"subtitle_language"`
   VideoType                string       `json:"video_type"`
   Player                   PlayerType   `json:"player"`
   ClassificationId         int          `json:"classification_id"`
   ContentType              string       `json:"content_type"`
   DeviceStreamVideoQuality VideoQuality `json:"device_stream_video_quality"`
   AudioLanguage            string       `json:"audio_language"`
   ContentId                string       `json:"content_id"`
}

// join takes a variable number of strings and returns them combined into one string without separators.
func join(strs ...string) string {
   return strings.Join(strs, "")
}

// makeStreamRequest handles the common logic for the POST stream request and
// parsing
func makeStreamRequest(marketCode, contentType, contentId string, player PlayerType, quality VideoQuality, audioLanguage string) (*StreamData, error) {
   classID, ok := classificationMap[marketCode]
   if !ok {
      return nil, fmt.Errorf("unsupported market code: %s", marketCode)
   }

   payload := StreamRequestPayload{
      AudioQuality:             "2.0",
      DeviceIdentifier:         DeviceID,
      DeviceSerial:             "not implemented",
      SubtitleLanguage:         "MIS",
      VideoType:                "stream",
      Player:                   player,
      ClassificationId:         classID,
      ContentType:              contentType,
      DeviceStreamVideoQuality: quality,
      AudioLanguage:            audioLanguage,
      ContentId:                contentId,
   }

   jsonData, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }

   apiURL := "https://gizmo.rakuten.tv/v3/avod/streamings"
   req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
   if err != nil {
      return nil, err
   }
   req.Header.Set("Content-Type", "application/json")
   resp, err := http.DefaultClient.Do(req)
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

// --- Movie Type and Methods ---

type Movie struct {
   Id         string // Matches "content_id" in URLs
   MarketCode string
}

// Request fetches movie details (GET).
func (m *Movie) Request() (*VideoItem, error) {
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

// RequestStream requests the stream for this movie (POST).
// Arguments: audioLanguage, player, quality.
func (m *Movie) RequestStream(audioLanguage string, player PlayerType, quality VideoQuality) (*StreamData, error) {
   return makeStreamRequest(m.MarketCode, "movies", m.Id, player, quality, audioLanguage)
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

// --- TvShow Type and Methods ---

type TvShow struct {
   Id         string // Matches "tv_show_id" in URLs
   MarketCode string
}

// Request fetches TV show details like seasons (GET).
func (t *TvShow) Request() (*TvShowData, error) {
   fullURL, err := buildURL(t.MarketCode, "tv_shows", t.Id)
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
func (t *TvShow) RequestSeason(seasonId string) (*SeasonData, error) {
   fullURL, err := buildURL(t.MarketCode, "seasons", seasonId)
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

// RequestStream requests the stream for a specific Episode (POST).
// Arguments: episodeId, audioLanguage, player, quality.
func (t *TvShow) RequestStream(episodeId string, audioLanguage string, player PlayerType, quality VideoQuality) (*StreamData, error) {
   // For TV content, the standard Rakuten content_type for streaming an episode is "episodes"
   return makeStreamRequest(t.MarketCode, "episodes", episodeId, player, quality, audioLanguage)
}
