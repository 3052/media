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

func (m *Media) streamInfo(
   content_id, audio_language, player string, video Quality,
) (*StreamInfo, error) {
   classificationID, err := m.classification_id()
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
      "content_type":                m.ContentType,
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

func (m *Media) Wvm(
   content_id, audio_language string, video Quality,
) (*StreamInfo, error) {
   return m.streamInfo(content_id, audio_language, ":DASH-CENC:WVM", video)
}

func (m *Media) Pr(
   content_id, audio_language string, video Quality,
) (*StreamInfo, error) {
   return m.streamInfo(content_id, audio_language, ":DASH-CENC:PR", video)
}
func (m *Media) Parse(source string) error {
   parsed, err := url.Parse(source)
   if err != nil {
      return err
   }
   path := strings.Split(strings.Trim(parsed.Path, "/"), "/")
   query := parsed.Query()
   if len(path) > 0 {
      m.MarketCode = path[0]
   }
   if content_type := query.Get("content_type"); content_type != "" {
      m.ContentType = content_type
      m.ContentId = query.Get("content_id")
      m.TvShowId = query.Get("tv_show_id")
   } else {
      if len(path) > 1 {
         m.ContentType = path[1]
      }
      if len(path) > 2 {
         if m.ContentType == "tv_shows" {
            m.TvShowId = path[2]
         } else {
            m.ContentId = path[2]
         }
      }
   }
   return nil
}

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
func (m *Media) classification_id() (int, error) {
   switch m.MarketCode {
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
type Media struct {
   ContentId   string
   ContentType string
   MarketCode  string
   TvShowId    string
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
