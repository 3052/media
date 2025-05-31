package rakuten

import (
   "encoding/json"
   "errors"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

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

const device_identifier = "atvui40"

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

type episode struct {
   Id    string
   Title string
}

func (e *episode) String() string {
   var b strings.Builder
   b.WriteString("title = ")
   b.WriteString(e.Title)
   b.WriteString("\nid = ")
   b.WriteString(e.Id)
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

type quality string

const (
   fhd quality = "FHD"
   hd  quality = "HD"
)

type streamings struct {
   AudioLanguage            string  `json:"audio_language"`
   AudioQuality             string  `json:"audio_quality"`
   ClassificationId         int     `json:"classification_id"`
   ContentId                string  `json:"content_id"`
   ContentType              string  `json:"content_type"`
   DeviceIdentifier         string  `json:"device_identifier"`
   DeviceSerial             string  `json:"device_serial"`
   DeviceStreamVideoQuality quality `json:"device_stream_video_quality"`
   Player                   string  `json:"player"`
   SubtitleLanguage         string  `json:"subtitle_language"`
   VideoType                string  `json:"video_type"`
}

func (e *episode) streamings(video quality) streamings {
   var s streamings
   s.ContentType = "episodes"
   s.DeviceStreamVideoQuality = video
   s.ContentId = e.Id
   s.AudioQuality = "2.0"
   s.DeviceIdentifier = device_identifier
   s.DeviceSerial = "not implemented"
   s.Player = device_identifier + ":DASH-CENC:WVM"
   s.SubtitleLanguage = "MIS"
   s.VideoType = "stream"
   
   s.AudioLanguage = ""
   s.ClassificationId = 0
   
   return s
}
