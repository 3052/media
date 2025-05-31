package rakuten

import (
   "bytes"
   "encoding/json"
   "errors"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

func (t *tv_show) info(
   content_id, audio_language string, video quality,
) (*stream_info, error) {
   data, err := json.Marshal(map[string]string{
      "audio_language":              audio_language,
      "audio_quality":               "2.0",
      "classification_id":           strconv.Itoa(t.classification_id()),
      "content_id":                  content_id,
      "content_type":                "episodes",
      "device_identifier":           device_identifier,
      "device_serial":               "not implemented",
      "device_stream_video_quality": string(video),
      "player":                      device_identifier + ":DASH-CENC:WVM",
      "subtitle_language":           "MIS",
      "video_type":                  "stream",
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://gizmo.rakuten.tv/v3/avod/streamings",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("content-type", "application/json")
   req.Header.Set("proxy", "true")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var value struct {
      Data struct {
         StreamInfos []stream_info `json:"stream_infos"`
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
func (e *episode) String() string {
   var b strings.Builder
   b.WriteString("title = ")
   b.WriteString(e.Title)
   for _, stream := range e.ViewOptions.Private.Streams {
      for _, language := range stream.AudioLanguages {
         b.WriteString("\nlanguage = ")
         b.WriteString(language.Id)
      }
   }
   b.WriteString("\nid = ")
   b.WriteString(e.Id)
   return b.String()
}

type episode struct {
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

const device_identifier = "atvui40"

type quality string

const (
   fhd quality = "FHD"
   hd  quality = "HD"
)

func (t *tv_show) seasons() ([]season, error) {
   req, _ := http.NewRequest("", "https://gizmo.rakuten.tv", nil)
   req.URL.Path = "/v3/tv_shows/" + t.tv_show_id
   req.URL.RawQuery = url.Values{
      "classification_id": {
         strconv.Itoa(t.classification_id()),
      },
      "device_identifier": {device_identifier},
      "market_code":       {t.market_code},
   }.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var value struct {
      Data struct {
         Seasons []season
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

type season struct {
   TvShowTitle string `json:"tv_show_title"`
   Id          string
}

func (s *season) String() string {
   var b strings.Builder
   b.WriteString("show title = ")
   b.WriteString(s.TvShowTitle)
   b.WriteString("\nid = ")
   b.WriteString(s.Id)
   return b.String()
}

func (t *tv_show) episodes(season_id string) ([]episode, error) {
   req, _ := http.NewRequest("", "https://gizmo.rakuten.tv", nil)
   req.URL.Path = "/v3/seasons/" + season_id
   req.URL.RawQuery = url.Values{
      "classification_id": {
         strconv.Itoa(t.classification_id()),
      },
      "device_identifier": {device_identifier},
      "market_code":       {t.market_code},
   }.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var value struct {
      Data struct {
         Episodes []episode
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   return value.Data.Episodes, nil
}

type tv_show struct {
   market_code string
   tv_show_id  string
}

func (t *tv_show) Set(data string) error {
   tv_url, err := url.Parse(data)
   if err != nil {
      return err
   }
   t.market_code = strings.TrimPrefix(tv_url.Path, "/")
   t.tv_show_id = tv_url.Query().Get("tv_show_id")
   return nil
}

func (t *tv_show) classification_id() int {
   switch t.market_code {
   case "at":
      return 300
   case "ch":
      return 319
   case "cz":
      return 272
   case "de":
      return 307
   case "fr":
      return 23
   case "ie":
      return 41
   case "nl":
      return 69
   case "pl":
      return 277
   case "se":
      return 282
   case "uk":
      return 18
   }
   return 0
}

type stream_info struct {
   LicenseUrl string `json:"license_url"`
   Url        string // MPD
}
