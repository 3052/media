package rakuten

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "log"
   "net/http"
   "net/url"
   "path"
   "strconv"
   "strings"
)

var Transport = http.Transport{
   Protocols: &http.Protocols{}, // github.com/golang/go/issues/25793
   Proxy: func(req *http.Request) (*url.URL, error) {
      switch path.Ext(req.URL.Path) {
      case ".isma", ".ismv":
      default:
         log.Println(req.Method, req.URL)
      }
      return http.ProxyFromEnvironment(req)
   },
}

// github.com/pandvan/rakuten-m3u-generator/blob/master/rakuten.py
func (a *Address) classification_id() (int, error) {
   switch a.MarketCode {
   case "cz":
      return 272, nil
   case "dk":
      return 283, nil
   case "fr":
      return 23, nil
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

type StreamInfo struct {
   // THIS URL GETS LOCKED TO DEVICE ON FIRST REQUEST
   LicenseUrl string `json:"license_url"`
   // MPD
   Url string
}

type Quality string

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

const device_identifier = "atvui40"

type Season struct {
   TvShowTitle string `json:"tv_show_title"`
   Id          string
}

// https://rakuten.tv/fr/movies/michael-clayton
// https://rakuten.tv/fr/tv_shows/une-femme-d-honneur
// https://rakuten.tv/fr?content_type=movies&content_id=michael-clayton
// https://rakuten.tv/fr?content_type=tv_shows&tv_show_id=une-femme-d-honneur&content_id=une-femme-d-honneur-1
type Address struct {
   ContentId   string
   ContentType string
   MarketCode  string
   TvShowId    string
}

func (a *Address) Seasons() ([]Season, error) {
   classificationID, err := a.classification_id()
   if err != nil {
      return nil, err
   }
   req, _ := http.NewRequest("", "https://gizmo.rakuten.tv", nil)
   req.URL.Path = "/v3/tv_shows/" + a.TvShowId
   req.URL.RawQuery = url.Values{
      "classification_id": {
         strconv.Itoa(classificationID),
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

func (a *Address) Movie() (*Content, error) {
   classificationID, err := a.classification_id()
   if err != nil {
      return nil, err
   }
   req, _ := http.NewRequest("", "https://gizmo.rakuten.tv", nil)
   req.URL.Path = "/v3/movies/" + a.ContentId
   req.URL.RawQuery = url.Values{
      "classification_id": {
         strconv.Itoa(classificationID),
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

func (a *Address) streamInfo(
   content_id, audio_language, player string, video Quality,
) (*StreamInfo, error) {
   classificationID, err := a.classification_id()
   if err != nil {
      return nil, err
   }
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
      "classification_id":           strconv.Itoa(classificationID),
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

func (a *Address) Episodes(season_id string) ([]Content, error) {
   classificationID, err := a.classification_id()
   if err != nil {
      return nil, err
   }
   req, _ := http.NewRequest("", "https://gizmo.rakuten.tv", nil)
   req.URL.Path = "/v3/seasons/" + season_id
   req.URL.RawQuery = url.Values{
      "classification_id": {
         strconv.Itoa(classificationID),
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

func (s *Season) String() string {
   var b strings.Builder
   b.WriteString("show title = ")
   b.WriteString(s.TvShowTitle)
   b.WriteString("\nseason id = ")
   b.WriteString(s.Id)
   return b.String()
}

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

func (a *Address) Parse(data string) error {
   parsed, err := url.Parse(data)
   if err != nil {
      return err
   }
   path := strings.Split(strings.Trim(parsed.Path, "/"), "/")
   query := parsed.Query()
   if len(path) > 0 {
      a.MarketCode = path[0]
   }
   if content_type := query.Get("content_type"); content_type != "" {
      a.ContentType = content_type
      a.ContentId = query.Get("content_id")
      a.TvShowId = query.Get("tv_show_id")
   } else {
      if len(path) > 1 {
         a.ContentType = path[1]
      }
      if len(path) > 2 {
         if a.ContentType == "tv_shows" {
            a.TvShowId = path[2]
         } else {
            a.ContentId = path[2]
         }
      }
   }
   return nil
}
