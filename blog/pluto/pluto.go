package pluto

import (
   "bytes"
   "encoding/json"
   "io"
   "net/http"
   "net/url"
)

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

type Series struct {
   Servers struct {
      StitcherDash string
   }
   SessionToken string
   Vod []struct {
      Stitched struct {
         Paths []struct {
            Path string
         }
      }
   }
}

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

func (s *Series) Fetch(id string) error {
   req, _ := http.NewRequest("", "https://boot.pluto.tv/v4/start", nil)
   req.URL.RawQuery = url.Values{
      "appName": {"androidtv"},
      "appVersion": {"9"},
      "clientID": {"9"},
      "clientModelNumber": {"9"},
      "deviceMake": {"9"},
      "deviceModel": {"9"},
      "deviceVersion": {"9"},
      "drmCapabilities": {"widevine:L1"},
      "seriesIDs": {id},
   }.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(s)
}
