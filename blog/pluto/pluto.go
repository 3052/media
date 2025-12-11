package main

import (
   "encoding/json"
   "io"
   "log"
   "net/http"
   "net/url"
   "strings"
)

func (s *series) alfa() (*url.URL, []byte, error) {
   req, err := http.NewRequest("", s.Servers.StitcherDash, nil)
   if err != nil {
      return nil, nil, err
   }
   req.URL.Path = "/v2" + s.Vod[0].Stitched.Paths[0].Path
   req.URL.RawQuery = "jwt=" + s.SessionToken
   req.Header.Set(
      "user-agent",
      //2025/12/10 22:47:25  duration="PT30.066S" true
      "Mozilla/5",
   )
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

type series struct {
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
      Id string
      Seasons []struct {
         Number   int64
         Episodes []struct {
            Number int64
            Name   string
            Id     string `json:"_id"`
         }
      }
   }
}

func main() {
   var req http.Request
   req.Header = http.Header{}
   req.Header.Add("Accept", "*/*")
   req.Header.Add("Accept-Language", "en-US,en;q=0.5")
   req.Header.Add("Cache-Control", "no-cache")
   req.Header.Add("Connection", "keep-alive")
   req.Header.Add("Content-Length", "0")
   req.Header.Add("Host", "boot.pluto.tv")
   req.Header.Add("Origin", "https://pluto.tv")
   req.Header.Add("Pragma", "no-cache")
   req.Header.Add("Priority", "u=4")
   req.Header.Add("Referer", "https://pluto.tv/")
   req.Header.Add("Sec-Fetch-Dest", "empty")
   req.Header.Add("Sec-Fetch-Mode", "cors")
   req.Header.Add("Sec-Fetch-Site", "same-site")
   req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:128.0) Gecko/20100101 Firefox/128.0")
   req.ProtoMajor = 1
   req.ProtoMinor = 1
   req.URL = &url.URL{}
   req.URL.Host = "boot.pluto.tv"
   req.URL.Path = "/v4/start"
   value := url.Values{}
   value["appLaunchCount"] = []string{"0"}
   value["appName"] = []string{"web"}
   value["appVersion"] = []string{"9.18.0-32296d47c9882754e360f1b28a33027c54cbad16"}
   value["blockingMode"] = []string{""}
   value["clientID"] = []string{"e0292ffd-7e8b-4607-ab89-fcd441a74b40"}
   value["clientModelNumber"] = []string{"1.0.0"}
   value["clientTime"] = []string{"2025-12-10T00:06:54.759Z"}
   value["deviceMake"] = []string{"firefox"}
   value["deviceModel"] = []string{"web"}
   value["deviceType"] = []string{"web"}
   value["deviceVersion"] = []string{"128.0.0"}
   value["drmCapabilities"] = []string{"widevine:L3"}
   value["lastAppLaunchDate"] = []string{"2025-12-10T00:06:54.758Z"}
   value["notificationVersion"] = []string{"1"}
   value["seriesIDs"] = []string{"6495eff09263a40013cf63a5"}
   value["serverSideAds"] = []string{"false"}
   req.URL.RawQuery = value.Encode()
   req.URL.Scheme = "https"
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      panic(err)
   }
   defer resp.Body.Close()
   var result series
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      panic(err)
   }
   const duration = ` duration="PT30.066S"`
   _, data, err := result.alfa()
   if err != nil {
      panic(err)
   }
   log.Println(duration, strings.Contains(string(data), duration))
}
