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

type StreamInfo struct {
   LicenseUrl string `json:"license_url"`
   Url        string // MPD
}

func (s *Streamings) Hd() {
   s.DeviceStreamVideoQuality = "HD"
}

func (s *Streamings) Fhd() {
   s.DeviceStreamVideoQuality = "FHD"
}

func (c *Content) Streamings() Streamings {
   return Streamings{ContentId: c.Id, ContentType: c.Type}
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

type Season struct {
   Episodes []Content
}

type Byte[T any] []byte

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

// hard geo block
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
   resp, err := http.DefaultClient.Do(req)
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
   // you can trigger this with wrong location
   if len(value.Errors) >= 1 {
      return nil, errors.New(value.Errors[0].Message)
   }
   return &value.Data.StreamInfos[0], nil
}

// github.com/pandvan/rakuten-m3u-generator/blob/master/rakuten.py
func (a *Address) ClassificationId() (int, bool) {
   switch a.MarketCode {
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

func (s Season) Content(web *Address) (*Content, bool) {
   for _, episode := range s.Episodes {
      if episode.Id == web.ContentId {
         return &episode, true
      }
   }
   return nil, false
}

func (a *Address) Season(classification_id int) (Byte[Season], error) {
   req, _ := http.NewRequest("", "https://gizmo.rakuten.tv", nil)
   req.URL.Path = "/v3/seasons/" + a.SeasonId
   req.URL.RawQuery = url.Values{
      "classification_id": {strconv.Itoa(classification_id)},
      "device_identifier": {"atvui40"},
      "market_code":       {a.MarketCode},
   }.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (a *Address) Movie(classification_id int) (Byte[Content], error) {
   req, _ := http.NewRequest("", "https://gizmo.rakuten.tv", nil)
   req.URL.Path = "/v3/movies/" + a.ContentId
   req.URL.RawQuery = url.Values{
      "classification_id": {strconv.Itoa(classification_id)},
      "device_identifier": {"atvui40"},
      "market_code":       {a.MarketCode},
   }.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

type Address struct {
   ContentId  string
   MarketCode string
   SeasonId   string
}

func (a *Address) Set(data string) error {
   data = strings.TrimPrefix(data, "https://")
   data = strings.TrimPrefix(data, "www.")
   data = strings.TrimPrefix(data, "rakuten.tv")
   data = strings.TrimPrefix(data, "/")
   a.MarketCode, data, _ = strings.Cut(data, "/")
   var found bool
   data, a.ContentId, found = strings.Cut(data, "movies/")
   if !found {
      data = strings.TrimPrefix(data, "player/episodes/stream/")
      a.SeasonId, a.ContentId, _ = strings.Cut(data, "/")
   }
   return nil
}
