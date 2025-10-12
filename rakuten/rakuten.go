package rakuten

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

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

// geo block
func (a *Address) streamInfo(
   content_id, audio_language, player string, video Quality,
) (*StreamInfo, error) {
   data, err := json.Marshal(map[string]string{
      "audio_language":              audio_language,
      "audio_quality":               "2.0",
      "device_serial":               "not implemented",
      "subtitle_language":           "MIS",
      "video_type":                  "stream",
      "device_identifier":           device_identifier,
      "content_id":                  content_id,
      "device_stream_video_quality": string(video),
      "player":                      device_identifier + player,
      "classification_id":           strconv.Itoa(a.classification_id()),
      "content_type":                a.ContentType,
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
   var value struct {
      Data struct {
         StreamInfos []StreamInfo `json:"stream_infos"`
      }
      Errors []struct {
         Message string
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   if len(value.Errors) >= 1 { // you can trigger this with wrong location
      return nil, errors.New(value.Errors[0].Message)
   }
   return &value.Data.StreamInfos[0], nil
}

type StreamInfo struct {
   // THIS URL GETS LOCKED TO DEVICE ON FIRST REQUEST
   LicenseUrl string `json:"license_url"`
   // MPD
   Url string
}

func (c *Content) String() string {
   var b strings.Builder
   b.WriteString("title = ")
   b.WriteString(c.Title)
   b.WriteString("\ncontent id = ")
   b.WriteString(c.Id)
   id := map[string]struct{}{}
   for _, stream := range c.ViewOptions.Private.Streams {
      for _, language := range stream.AudioLanguages {
         _, ok := id[language.Id]
         if !ok {
            b.WriteString("\naudio language = ")
            b.WriteString(language.Id)
            id[language.Id] = struct{}{}
         }
      }
   }
   return b.String()
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

// github.com/pandvan/rakuten-m3u-generator/blob/master/rakuten.py
func (a *Address) classification_id() int {
   switch a.MarketCode {
   case "cz":
      return 272
   case "fr":
      return 23
   case "pl":
      return 277
   case "se":
      return 282
   case "uk":
      return 18
   }
   return 0
}

func (a *Address) Set(data string) error {
   web, err := url.Parse(data)
   if err != nil {
      return err
   }
   a.ContentId = web.Query().Get("content_id")
   a.ContentType = web.Query().Get("content_type")
   a.MarketCode = strings.TrimPrefix(web.Path, "/")
   a.TvShowId = web.Query().Get("tv_show_id")
   return nil
}

func (a *Address) Movie() (*Content, error) {
   req, _ := http.NewRequest("", "https://gizmo.rakuten.tv", nil)
   req.URL.Path = "/v3/movies/" + a.ContentId
   req.URL.RawQuery = url.Values{
      "classification_id": {
         strconv.Itoa(a.classification_id()),
      },
      "device_identifier": {device_identifier},
      "market_code":       {a.MarketCode},
   }.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var value struct {
      Data Content
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   return &value.Data, nil
}

func (a *Address) Episodes(season_id string) ([]Content, error) {
   req, _ := http.NewRequest("", "https://gizmo.rakuten.tv", nil)
   req.URL.Path = "/v3/seasons/" + season_id
   req.URL.RawQuery = url.Values{
      "classification_id": {
         strconv.Itoa(a.classification_id()),
      },
      "device_identifier": {device_identifier},
      "market_code":       {a.MarketCode},
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

// rakuten.tv/se?content_type=movies&content_id=i-heart-huckabees
// rakuten.tv/uk?content_type=tv_shows&tv_show_id=clink
type Address struct {
   ContentId   string
   ContentType string
   MarketCode  string
   TvShowId    string
}

func (s *Season) String() string {
   var b strings.Builder
   b.WriteString("show title = ")
   b.WriteString(s.TvShowTitle)
   b.WriteString("\nseason id = ")
   b.WriteString(s.Id)
   return b.String()
}

type Season struct {
   TvShowTitle string `json:"tv_show_title"`
   Id          string
}

const device_identifier = "atvui40"

func (a *Address) Seasons() ([]Season, error) {
   req, _ := http.NewRequest("", "https://gizmo.rakuten.tv", nil)
   req.URL.Path = "/v3/tv_shows/" + a.TvShowId
   req.URL.RawQuery = url.Values{
      "classification_id": {
         strconv.Itoa(a.classification_id()),
      },
      "device_identifier": {device_identifier},
      "market_code":       {a.MarketCode},
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

const (
   Fhd Quality = "FHD"
   Hd  Quality = "HD"
)

func (a *Address) Wvm(
   content_id, audio_language string, video Quality,
) (*StreamInfo, error) {
   return a.streamInfo(content_id, audio_language, ":DASH-CENC:WVM", video)
}

func (a *Address) Pr(
   content_id, audio_language string, video Quality,
) (*StreamInfo, error) {
   return a.streamInfo(content_id, audio_language, ":DASH-CENC:PR", video)
}
