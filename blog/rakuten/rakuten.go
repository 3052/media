package rakuten

import (
   "encoding/json"
   "errors"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

type address struct {
   market_code string
   tv_show_id  string
}

func (a *address) seasons() ([]season, error) {
   req, _ := http.NewRequest("", "https://gizmo.rakuten.tv", nil)
   req.URL.Path = "/v3/tv_shows/" + a.tv_show_id
   req.URL.RawQuery = url.Values{
      "classification_id": {
         strconv.Itoa(a.classification_id()),
      },
      "device_identifier": {"web"},
      "market_code":       {a.market_code},
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

func (s *season) String() string {
   var b strings.Builder
   b.WriteString("show title = ")
   b.WriteString(s.TvShowTitle)
   b.WriteString("\nid = ")
   b.WriteString(s.Id)
   return b.String()
}

type season struct {
   TvShowTitle string `json:"tv_show_title"`
   Id          string
}

func (a *address) Set(data string) error {
   url2, err := url.Parse(data)
   if err != nil {
      return err
   }
   a.market_code = strings.TrimPrefix(url2.Path, "/")
   a.tv_show_id = url2.Query().Get("tv_show_id")
   return nil
}

func (a *address) classification_id() int {
   switch a.market_code {
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

func (s *streamings) Fhd() {
   s.DeviceStreamVideoQuality = "FHD"
}

func (e *episode) streamings() streamings {
   return streamings{ContentId: e.Id, ContentType: "episodes"}
}

type streamings struct {
   AudioLanguage            string `json:"audio_language"`
   AudioQuality             string `json:"audio_quality"`
   ClassificationId         int    `json:"classification_id"`
   ContentId                string `json:"content_id"`
   ContentType              string `json:"content_type"`
   DeviceIdentifier         string `json:"device_identifier"`
   DeviceSerial             string `json:"device_serial"`
   DeviceStreamVideoQuality string `json:"device_stream_video_quality"`
   Player                   string `json:"player"`
   SubtitleLanguage         string `json:"subtitle_language"`
   VideoType                string `json:"video_type"`
}

func (e *episode) String() string {
   var b strings.Builder
   b.WriteString("title = ")
   b.WriteString(e.Title)
   b.WriteString("\nid = ")
   b.WriteString(e.Id)
   return b.String()
}

type episode struct {
   Id    string
   Title string
}

func (a *address) episodes(season_id string) ([]episode, error) {
   req, _ := http.NewRequest("", "https://gizmo.rakuten.tv", nil)
   req.URL.Path = "/v3/seasons/" + season_id
   req.URL.RawQuery = url.Values{
      "classification_id": {
         strconv.Itoa(a.classification_id()),
      },
      "device_identifier": {"web"},
      "market_code":       {a.market_code},
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
