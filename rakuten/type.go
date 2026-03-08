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

// Media represents a piece of content, which can be a Movie or a TV Show.
type Media struct {
   Id         string // Matches "content_id" or "tv_show_id" in URLs
   MarketCode string
   Type       ContentType
}

func ParseMedia(rawLink string) (*Media, error) {
   link, err := url.Parse(rawLink)
   if err != nil {
      return nil, err
   }
   // Assuming extractMarketCode is defined elsewhere in your package
   marketCode, err := extractMarketCode(link.Path)
   if err != nil {
      return nil, err
   }
   // Initialize the struct here
   m := Media{MarketCode: marketCode}
   // 1. Check Query Parameters
   query := link.Query()
   contentType := query.Get("content_type")
   if contentType == "movies" || contentType == "tv_shows" {
      var id string
      if contentType == "movies" {
         id = query.Get("content_id")
         if id == "" {
            return nil, errors.New("url missing content_id param")
         }
      } else {
         id = query.Get("tv_show_id")
         if id == "" {
            return nil, errors.New("url missing tv_show_id param")
         }
      }
      m.Id = id
      m.Type = ContentType(contentType)
      return &m, nil
   }
   // 2. Check Path Segments
   path := strings.Trim(link.Path, "/")
   segments := strings.Split(path, "/")
   for _, seg := range segments {
      if seg == "movies" || seg == "tv_shows" {
         id := segments[len(segments)-1]
         if id == seg {
            return nil, fmt.Errorf("url does not contain a specific %s id", seg)
         }
         m.Id = id
         m.Type = ContentType(seg)
         return &m, nil
      }
   }
   return nil, errors.New("not a movie or tv show url")
}

type ContentType string

const (
   Movies  ContentType = "movies"
   TvShows ContentType = "tv_shows"
)

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

// MovieStream requests the stream for this movie (POST).
// The movie's own ID is used as the contentId.
func (m *Media) MovieStream(audioLanguage string, player PlayerType, quality VideoQuality) (*StreamData, error) {
   if m.Type != Movies {
      return nil, errors.New("cannot request a movie stream for non-movie content")
   }
   return makeStreamRequest(m.MarketCode, "movies", m.Id, player, quality, audioLanguage)
}

// EpisodeStream requests the stream for a specific TV Show Episode (POST).
func (m *Media) EpisodeStream(episodeId, audioLanguage string, player PlayerType, quality VideoQuality) (*StreamData, error) {
   if m.Type != TvShows {
      return nil, errors.New("cannot request an episode stream for non-tv-show content")
   }
   // For TV content, the standard Rakuten content_type for streaming an episode is "episodes"
   return makeStreamRequest(m.MarketCode, "episodes", episodeId, player, quality, audioLanguage)
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
