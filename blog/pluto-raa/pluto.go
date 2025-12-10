package main

import (
   "encoding/json"
   "fmt"
   "io"
   "log"
   "net/http"
   "net/url"
   "strings"
)

func main() {
   req, _ := http.NewRequest("", "https://boot.pluto.tv/v4/start", nil)
   req.URL.RawQuery = url.Values{
      "appName": {"androidtv"},
      "appVersion": {"5.53.0-leanback"},
      "clientID": {"720234b6-ce56-462a-892a-cf0d80c51469_2a547545129d6564"},
      "clientModelNumber": {"sdk_google_atv_x86"},
      "deviceMake": {"unknown"},
      "deviceModel": {"sdk_google_atv_x86"},
      "deviceVersion": {"9_28"},
      "drmCapabilities": {"widevine:L1"},
      "seriesIDs": {"6495eff09263a40013cf63a5"},
   }.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      panic(err)
   }
   defer resp.Body.Close()
   var result start
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      panic(err)
   }
   data, err := get(result.String())
   if err != nil {
      panic(err)
   }
   if strings.Contains(string(data), ` height="1080"`) {
      log.Print("pass")
   } else {
      log.Print("fail")
   }
}

func get(address string) ([]byte, error) {
   resp, err := http.Get(address)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (s *start) String() string {
   return fmt.Sprint(
      s.Servers.StitcherDash,
      "/v2", s.Vod[0].Stitched.Paths[0].Path,
      "?jwt=", s.SessionToken,
   )
}

type start struct {
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
