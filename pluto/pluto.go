package pluto

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "strconv"
)

var (
   app_name         = "androidtv"
   drm_capabilities = "widevine:L1"
)

func Widevine(data []byte) ([]byte, error) {
   resp, err := http.Post(
      "https://service-concierge.clusters.pluto.tv/v1/wv/alt",
      "application/x-protobuf", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (s *Series) Mpd() (*url.URL, []byte, error) {
   req, err := http.NewRequest("", s.Servers.StitcherDash, nil)
   if err != nil {
      return nil, nil, err
   }
   req.URL.Path = "/v2" + s.Vod[0].Stitched.Paths[0].Path
   req.URL.RawQuery = "jwt=" + s.SessionToken
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, nil, err
   }
   defer resp.Body.Close()
   data, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, nil, err
   }
   return resp.Request.URL, data, nil
}

func (s *Series) Fetch(id string) error {
   req, _ := http.NewRequest("", "https://boot.pluto.tv/v4/start", nil)
   req.URL.RawQuery = url.Values{
      "appName":           {app_name},
      "appVersion":        {"9"},
      "clientID":          {"9"},
      "clientModelNumber": {"9"},
      "deviceMake":        {"9"},
      "deviceModel":       {"9"},
      "deviceVersion":     {"9"},
      "drmCapabilities":   {drm_capabilities},
      "seriesIDs":         {id},
   }.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   err = json.NewDecoder(resp.Body).Decode(s)
   if err != nil {
      return err
   }
   if s.Vod[0].Id != id {
      return errors.New("id mismatch")
   }
   return nil
}

type Series struct {
   Servers struct {
      StitcherDash string
   }
   SessionToken string
   Vod          []Vod
}

type Vod struct {
   Stitched struct {
      Paths []struct {
         Path string
      }
   }
   Id      string
   Seasons []struct {
      Number   int64
      Episodes []struct {
         Number int64
         Name   string
         Id     string `json:"_id"`
      }
   }
}

func (v *Vod) String() string {
   var (
      data  []byte
      lines bool
   )
   for _, season := range v.Seasons {
      for _, episode := range season.Episodes {
         if lines {
            data = append(data, "\n\n"...)
         } else {
            lines = true
         }
         data = append(data, "season = "...)
         data = strconv.AppendInt(data, season.Number, 10)
         data = append(data, "\nepisode = "...)
         data = strconv.AppendInt(data, episode.Number, 10)
         data = append(data, "\nname = "...)
         data = append(data, episode.Name...)
         data = append(data, "\nid = "...)
         data = append(data, episode.Id...)
      }
   }
   return string(data)
}
