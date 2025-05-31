package rakuten

import (
   "encoding/json"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

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
