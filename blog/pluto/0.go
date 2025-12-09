package main

import (
   "net/http"
   "net/url"
   "fmt"
   "encoding/json"
)

func main() {
   var req http.Request
   req.Header = http.Header{}
   req.URL = &url.URL{}
   req.URL.Host = "boot.pluto.tv"
   req.URL.Path = "/v4/start"
   value := url.Values{}
   value["appName"] = []string{"web"}
   value["appVersion"] = []string{"9.18.0-32296d47c9882754e360f1b28a33027c54cbad16"}
   value["clientID"] = []string{"eb722591-733e-459e-8d58-97c7bbbdf18c"}
   value["clientModelNumber"] = []string{"1.0.0"}
   value["drmCapabilities"] = []string{"widevine:L3"}
   value["seriesIDs"] = []string{"6495eff09263a40013cf63a5"}
   req.URL.RawQuery = value.Encode()
   req.URL.Scheme = "https"
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      panic(err)
   }
   defer resp.Body.Close()
   var result struct {
      Vod []struct {
         Stitched struct {
            Paths []struct {
               Path string
            }
         }
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      panic(err)
   }
   fmt.Printf("%+v\n", result)
}

/*

https://cfd-v4-service-stitcher-dash-use1-1.prd.pluto.tv/v2/stitch/dash/episode/6495eff09263a40013cf63a5/main.mpd

*/

