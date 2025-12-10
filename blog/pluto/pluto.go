package pluto

import (
   "encoding/json"
   "net/http"
   "net/url"
   "strings"
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

func (s *Series) String() string {
   var data strings.Builder
   data.WriteString(s.Servers.StitcherDash)
   data.WriteString("/v2")
   data.WriteString(s.Vod[0].Stitched.Paths[0].Path)
   data.WriteString("?jwt=")
   data.WriteString(s.SessionToken)
   return data.String()
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
