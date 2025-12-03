package rakuten

import (
   "bytes"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "log"
   "net/http"
   "net/url"
   "path"
   "strconv"
   "strings"
)

func (c *Content) String() string {
   var data strings.Builder
   data.WriteString("title = ")
   data.WriteString(c.Title)
   data.WriteString("\ncontent id = ")
   data.WriteString(c.Id)
   id := map[string]struct{}{}
   for _, stream := range c.ViewOptions.Private.Streams {
      for _, language := range stream.AudioLanguages {
         _, ok := id[language.Id]
         if !ok {
            data.WriteString("\naudio language = ")
            data.WriteString(language.Id)
            id[language.Id] = struct{}{}
         }
      }
   }
   return data.String()
}

type Content struct {
   Id          string
   Title       string
   ViewOptions struct {
      Private struct {
         Streams []struct {
            AudioLanguages []struct {
               Id string
            } `json:"audio_languages"`
         }
      }
   } `json:"view_options"`
}

func (m *Media) Seasons() ([]Season, error) {
   classificationID, err := m.classification_id()
   if err != nil {
      return nil, err
   }
   req, _ := http.NewRequest("", "https://gizmo.rakuten.tv", nil)
   req.URL.Path = "/v3/tv_shows/" + m.TvShowId
   req.URL.RawQuery = url.Values{
      "classification_id": {
         strconv.Itoa(classificationID),
      },
      "device_identifier": {device_identifier},
      "market_code":       {m.MarketCode},
   }.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var value struct {
      Data struct {
         Seasons []Season
      }
      Errors []struct {
         Code string
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   if len(value.Errors) >= 1 {
      return nil, errors.New(value.Errors[0].Code)
   }
   return value.Data.Seasons, nil
}

type Quality string

func (m *Media) Movie() (*Content, error) {
   classificationID, err := m.classification_id()
   if err != nil {
      return nil, err
   }
   req, _ := http.NewRequest("", "https://gizmo.rakuten.tv", nil)
   req.URL.Path = "/v3/movies/" + m.ContentId
   req.URL.RawQuery = url.Values{
      "classification_id": {
         strconv.Itoa(classificationID),
      },
      "device_identifier": {device_identifier},
      "market_code":       {m.MarketCode},
   }.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var value struct {
      Data   Content
      Errors []struct {
         Message string
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   if len(value.Errors) >= 1 {
      return nil, errors.New(value.Errors[0].Message)
   }
   return &value.Data, nil
}

func (m *Media) Episodes(seasonId string) ([]Content, error) {
   classificationID, err := m.classification_id()
   if err != nil {
      return nil, err
   }
   req, _ := http.NewRequest("", "https://gizmo.rakuten.tv", nil)
   req.URL.Path = "/v3/seasons/" + seasonId
   req.URL.RawQuery = url.Values{
      "classification_id": {
         strconv.Itoa(classificationID),
      },
      "device_identifier": {device_identifier},
      "market_code":       {m.MarketCode},
   }.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var value struct {
      Data struct {
         Episodes []Content
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   return value.Data.Episodes, nil
}

const device_identifier = "atvui40"

const (
   Fhd Quality = "FHD"
   Hd  Quality = "HD"
)

func (m *Media) Pr(
   content_id, audio_language string, video Quality,
) (*StreamInfo, error) {
   return m.streamInfo(content_id, audio_language, ":DASH-CENC:PR", video)
}

func (m *Media) Wvm(
   content_id, audio_language string, video Quality,
) (*StreamInfo, error) {
   return m.streamInfo(content_id, audio_language, ":DASH-CENC:WVM", video)
}

func (m *Media) streamInfo(
   content_id, audio_language, player string, video Quality,
) (*StreamInfo, error) {
   classificationID, err := m.classification_id()
   if err != nil {
      return nil, err
   }
   data, err := json.Marshal(map[string]string{
      "audio_quality":               "2.0",
      "device_serial":               "not implemented",
      "subtitle_language":           "MIS",
      "video_type":                  "stream",
      "content_type":                m.ContentType,
      "device_identifier":           device_identifier,
      "audio_language":              audio_language,
      "content_id":                  content_id,
      "device_stream_video_quality": string(video),
      "player":                      device_identifier + player,
      "classification_id":           strconv.Itoa(classificationID),
   })
   if err != nil {
      return nil, err
   }
   resp, err := http.Post(
      "https://gizmo.rakuten.tv/v3/avod/streamings",
      "application/json", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   data, err = io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   if len(data) == 0 {
      return nil, errors.New(resp.Status)
   }
   var value struct {
      Data struct {
         StreamInfos []StreamInfo `json:"stream_infos"`
      }
      Errors []struct {
         Message string
      }
   }
   err = json.Unmarshal(data, &value)
   if err != nil {
      return nil, err
   }
   if len(value.Errors) >= 1 {
      return nil, errors.New(value.Errors[0].Message)
   }
   return &value.Data.StreamInfos[0], nil
}

type Media struct {
   ContentId   string
   ContentType string
   MarketCode  string
   TvShowId    string
}

func (m *Media) Parse(rawUrl string) error {
   parsed, err := url.Parse(rawUrl)
   if err != nil {
      return fmt.Errorf("failed to parse URL: %w", err)
   }
   if parsed.Scheme == "" {
      return errors.New("invalid URL: scheme is missing")
   }
   path := strings.Trim(parsed.Path, "/")
   pathParts := strings.Split(path, "/")
   // The first part of the path is the MarketCode.
   if len(pathParts) > 0 && len(pathParts[0]) == 2 {
      m.MarketCode = pathParts[0]
   } else {
      return fmt.Errorf("could not determine market code from URL path")
   }
   // First, try to parse content info from the path.
   if len(pathParts) > 1 {
      // Handle direct content links like /movies/... or /tv_shows/...
      if pathParts[1] == "movies" && len(pathParts) == 3 {
         m.ContentType = "movies"
         m.ContentId = pathParts[2]
         return nil
      }
      if pathParts[1] == "tv_shows" && len(pathParts) == 3 {
         m.ContentType = "tv_shows"
         m.TvShowId = pathParts[2]
         return nil
      }
      if pathParts[1] == "player" {
         if len(pathParts) == 5 {
            if pathParts[2] == "movies" && pathParts[3] == "stream" {
               m.ContentType = "movies"
               m.ContentId = pathParts[4]
               return nil
            }
         }
      }
   }
   // If not in the path, fall back to checking query parameters.
   query := parsed.Query()
   contentType := query.Get("content_type")
   if contentType != "" {
      m.ContentType = contentType
      if contentType == "movies" {
         if contentID := query.Get("content_id"); contentID != "" {
            m.ContentId = contentID
            return nil
         }
      } else if contentType == "tv_shows" {
         if tvShowID := query.Get("tv_show_id"); tvShowID != "" {
            m.TvShowId = tvShowID
            return nil
         }
      }
   }
   return fmt.Errorf("could not parse content type and ID from URL")
}

// github.com/pandvan/rakuten-m3u-generator/blob/master/rakuten.py
func (m *Media) classification_id() (int, error) {
   switch m.MarketCode {
   case "cz":
      return 272, nil
   case "dk":
      return 283, nil
   case "es":
      return 5, nil
   case "fr":
      return 23, nil
   case "nl":
      return 69, nil
   case "pl":
      return 277, nil
   case "pt":
      return 64, nil
   case "se":
      return 282, nil
   case "uk":
      return 18, nil
   }
   return 0, errors.New("unknown market code")
}

type Season struct {
   TvShowTitle string `json:"tv_show_title"`
   Id          string
}

type StreamInfo struct {
   // THIS URL GETS LOCKED TO DEVICE ON FIRST REQUEST
   LicenseUrl string `json:"license_url"`
   // MPD
   Url string
}

func (s *StreamInfo) Widevine(data []byte) ([]byte, error) {
   resp, err := http.Post(
      s.LicenseUrl, "application/x-protobuf", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (s *Season) String() string {
   var data strings.Builder
   data.WriteString("show title = ")
   data.WriteString(s.TvShowTitle)
   data.WriteString("\nseason id = ")
   data.WriteString(s.Id)
   return data.String()
}
