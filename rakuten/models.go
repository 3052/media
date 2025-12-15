package rakuten

import (
   "fmt"
   "net/url"
   "strings"
)

// ParseURL attempts to parse a URL and populate the Movie struct.
// Usage:
//
//   var m Movie
//   err := m.ParseURL("https://rakuten.tv/...")
func (m *Movie) ParseURL(rawURL string) error {
   u, err := url.Parse(rawURL)
   if err != nil {
      return fmt.Errorf("invalid url: %w", err)
   }

   marketCode, err := extractMarketCode(u.Path)
   if err != nil {
      return err
   }

   // 1. Check Query Parameters
   // Pattern: ?content_type=movies&content_id=...
   q := u.Query()
   if q.Get("content_type") == "movies" {
      id := q.Get("content_id")
      if id == "" {
         return fmt.Errorf("url missing content_id param")
      }

      m.MarketCode = marketCode
      m.Id = id
      return nil
   }

   // 2. Check Path Segments
   // Pattern: /nl/movies/id OR /nl/player/movies/stream/id
   path := strings.Trim(u.Path, "/")
   segments := strings.Split(path, "/")

   for _, seg := range segments {
      if seg == "movies" {
         // Assuming the ID is the last segment
         id := segments[len(segments)-1]
         if id == "movies" {
            return fmt.Errorf("url does not contain a specific movie id")
         }

         m.MarketCode = marketCode
         m.Id = id
         return nil
      }
   }

   return fmt.Errorf("not a movie url")
}

// ParseURL attempts to parse a URL and populate the TvShow struct.
// Usage:
//
//   var t TvShow
//   err := t.ParseURL("https://rakuten.tv/...")
func (t *TvShow) ParseURL(rawURL string) error {
   u, err := url.Parse(rawURL)
   if err != nil {
      return fmt.Errorf("invalid url: %w", err)
   }

   marketCode, err := extractMarketCode(u.Path)
   if err != nil {
      return err
   }

   // 1. Check Query Parameters
   // Pattern: ?content_type=tv_shows&tv_show_id=...
   q := u.Query()
   if q.Get("content_type") == "tv_shows" {
      id := q.Get("tv_show_id")
      if id == "" {
         return fmt.Errorf("url missing tv_show_id param")
      }

      t.MarketCode = marketCode
      t.Id = id
      return nil
   }

   // 2. Check Path Segments
   // Pattern: /nl/tv_shows/id
   path := strings.Trim(u.Path, "/")
   segments := strings.Split(path, "/")

   for _, seg := range segments {
      if seg == "tv_shows" {
         // Assuming the ID is the last segment
         id := segments[len(segments)-1]
         if id == "tv_shows" {
            return fmt.Errorf("url does not contain a specific tv show id")
         }

         t.MarketCode = marketCode
         t.Id = id
         return nil
      }
   }

   return fmt.Errorf("not a tv show url")
}

// extractMarketCode extracts the first segment of the path (e.g., "nl", "uk").
func extractMarketCode(path string) (string, error) {
   trimmed := strings.Trim(path, "/")

   // Check if we have anything left after trimming
   if trimmed == "" {
      return "", fmt.Errorf("could not determine market code from path")
   }

   segments := strings.Split(trimmed, "/")
   return segments[0], nil
}
// String implements the fmt.Stringer interface.
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

// classificationMap maps market codes to their internal classification IDs.
var classificationMap = map[string]int{
   "cz": 272,
   "dk": 283,
   "es": 5,
   "fr": 23,
   "nl": 69,
   "pl": 277,
   "pt": 64,
   "se": 282,
   "uk": 18,
}

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
