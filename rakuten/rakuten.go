package rakuten

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "log"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

type Content struct {
   ViewOptions struct {
      Private struct {
         Streams []struct {
            AudioLanguages []struct {
               Id string
            } `json:"audio_languages"`
         }
      }
   } `json:"view_options"`
   Id   string
   Type string
}

func (s *Season) Unmarshal(data Byte[Season]) error {
   var value struct {
      Data Season
   }
   err := json.Unmarshal(data, &value)
   if err != nil {
      return err
   }
   *s = value.Data
   return nil
}

type Season struct {
   Episodes []Content
}

type Byte[T any] []byte

type Path struct {
   SeasonId   string
   MarketCode string
   ContentId  string
}

func (p *Path) New(data string) {
   data = strings.TrimPrefix(data, "https://")
   data = strings.TrimPrefix(data, "www.")
   data = strings.TrimPrefix(data, "rakuten.tv")
   data = strings.TrimPrefix(data, "/")
   p.MarketCode, data, _ = strings.Cut(data, "/")
   var found bool
   data, p.ContentId, found = strings.Cut(data, "movies/")
   if !found {
      data = strings.TrimPrefix(data, "player/episodes/stream/")
      p.SeasonId, p.ContentId, _ = strings.Cut(data, "/")
   }
}

func (p *Path) Season(classification_id int) (Byte[Season], error) {
   req, _ := http.NewRequest("", "https://gizmo.rakuten.tv", nil)
   req.URL.Path = "/v3/seasons/" + p.SeasonId
   req.URL.RawQuery = url.Values{
      "device_identifier": {"atvui40"},
      "classification_id": {strconv.Itoa(classification_id)},
      "market_code":       {p.MarketCode},
   }.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

// github.com/pandvan/rakuten-m3u-generator/blob/master/rakuten.py
func (p *Path) ClassificationId() (int, bool) {
   switch p.MarketCode {
   case "at":
      return 300, true
   case "ch":
      return 319, true
   case "cz":
      return 272, true
   case "de":
      return 307, true
   case "fr":
      return 23, true
   case "ie":
      return 41, true
   case "nl":
      return 69, true
   case "pl":
      return 277, true
   case "se":
      return 282, true
   case "uk":
      return 18, true
   }
   return 0, false
}

type Streamings struct {
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

func (c *Content) Unmarshal(data Byte[Content]) error {
   var value struct {
      Data   Content
      Errors []struct {
         Message string
      }
   }
   err := json.Unmarshal(data, &value)
   if err != nil {
      return err
   }
   *c = value.Data
   return nil
}

func (c *Content) String() string {
   var (
      audio = map[string]struct{}{}
      b     []byte
   )
   for _, stream := range c.ViewOptions.Private.Streams {
      for _, language := range stream.AudioLanguages {
         _, ok := audio[language.Id]
         if !ok {
            if b != nil {
               b = append(b, '\n')
            }
            b = append(b, "audio language = "...)
            b = append(b, language.Id...)
            audio[language.Id] = struct{}{}
         }
      }
   }
   b = append(b, "\nid = "...)
   b = append(b, c.Id...)
   b = append(b, "\ntype = "...)
   b = append(b, c.Type...)
   return string(b)
}

type StreamInfo struct {
   LicenseUrl string `json:"license_url"`
   Url        string // MPD
}

func (c *Content) Streamings() Streamings {
   return Streamings{ContentId: c.Id, ContentType: c.Type}
}

func (s *Streamings) Hd() {
   s.DeviceStreamVideoQuality = "HD"
}

func (s *Streamings) Fhd() {
   s.DeviceStreamVideoQuality = "FHD"
}

func (s Season) Content(path1 *Path) (*Content, bool) {
   for _, episode := range s.Episodes {
      if episode.Id == path1.ContentId {
         return &episode, true
      }
   }
   return nil, false
}

func (s *Streamings) Info(
   audio_language string, classification_id int,
) (*StreamInfo, error) {
   s.AudioLanguage = audio_language
   s.AudioQuality = "2.0"
   s.ClassificationId = classification_id
   s.DeviceIdentifier = "atvui40"
   s.DeviceSerial = "not implemented"
   s.Player = "atvui40:DASH-CENC:WVM"
   s.SubtitleLanguage = "MIS"
   s.VideoType = "stream"
   data, err := json.Marshal(s)
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
         Id          string
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
   // you can trigger this with wrong location
   if len(value.Errors) >= 1 {
      return nil, errors.New(value.Errors[0].Message)
   }
   log.Println("id", value.Data.Id)
   return &value.Data.StreamInfos[0], nil
}

///

func (p *Path) Movie(classification_id int) (Byte[Content], error) {
   req, _ := http.NewRequest("", "https://gizmo.rakuten.tv", nil)
   req.URL.Path = "/v3/movies/" + p.ContentId
   req.URL.RawQuery = url.Values{
      "classification_id": {strconv.Itoa(classification_id)},
      "device_identifier": {"atvui40"},
      "market_code":       {p.MarketCode},
   }.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (s *StreamInfo) License(data []byte) ([]byte, error) {
   resp, err := http.Post(
      s.LicenseUrl, "application/x-protobuf", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}
