package rakuten

import (
   "bytes"
   "encoding/json"
   "errors"
   "fmt"
   "net/http"
   "net/url"
   "strings"
)

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

// VideoItem represents a movie or episode.
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

// DeviceID is the default identifier used for requests.
const DeviceID = "atvui40"

type VideoQuality string

const (
   Fhd VideoQuality = "FHD"
   Hd  VideoQuality = "HD"
)

type PlayerType string

const (
   PlayReady PlayerType = DeviceID + ":DASH-CENC:PR"
   Widevine  PlayerType = DeviceID + ":DASH-CENC:WVM"
)

type StreamRequestPayload struct {
   AudioLanguage            string       `json:"audio_language"`
   AudioQuality             string       `json:"audio_quality"`
   ClassificationId         int          `json:"classification_id"`
   ContentId                string       `json:"content_id"`
   ContentType              string       `json:"content_type"`
   DeviceIdentifier         string       `json:"device_identifier"`
   DeviceSerial             string       `json:"device_serial"`
   DeviceStreamVideoQuality VideoQuality `json:"device_stream_video_quality"`
   Player                   PlayerType   `json:"player"`
   SubtitleLanguage         string       `json:"subtitle_language"`
   VideoType                string       `json:"video_type"`
}

// ContentTypeCategory defines the category of the media (Movie vs TV Show).
// Renamed to avoid collision with 'contentType' variables.
type ContentTypeCategory int

const (
   Movies ContentTypeCategory = iota
   TvShows
)

// Media represents a piece of content, which can be a Movie or a TV Show.
type Media struct {
   Id         string // Matches "content_id" or "tv_show_id" in URLs
   MarketCode string
   Type       ContentTypeCategory
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

// MovieStream requests the stream for this movie (POST).
func (m *Media) MovieStream(audioLanguage string, player PlayerType, quality VideoQuality) (*StreamData, error) {
   if m.Type != Movies {
      return nil, errors.New("cannot request a movie stream for non-movie content")
   }
   // For movies, the API expects content_type="movies"
   return makeStreamRequest(m.MarketCode, "movies", m.Id, player, quality, audioLanguage)
}

// EpisodeStream requests the stream for a specific TV Show Episode (POST).
func (m *Media) EpisodeStream(episodeId, audioLanguage string, player PlayerType, quality VideoQuality) (*StreamData, error) {
   if m.Type != TvShows {
      return nil, errors.New("cannot request an episode stream for non-tv-show content")
   }
   // For TV episodes, the API expects content_type="episodes" (different from URL "tv_shows")
   return makeStreamRequest(m.MarketCode, "episodes", episodeId, player, quality, audioLanguage)
}

// makeStreamRequest handles the common logic for the POST stream request.
// The 'contentType' argument here refers to the API payload value.
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
      ContentType:              contentType, // Field matches struct, value matches arg
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
