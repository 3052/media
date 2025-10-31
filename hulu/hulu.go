package hulu

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "strings"
)

func (p *Playlist) PlayReady(data []byte) ([]byte, error) {
   resp, err := http.Post(
      p.DashPrServer, "", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   data, err = io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   if resp.StatusCode != http.StatusOK {
      var value struct {
         Message string
      }
      err = json.Unmarshal(data, &value)
      if err != nil {
         return nil, err
      }
      return nil, errors.New(value.Message)
   }
   return data, nil
}

func (p *Playlist) Widevine(data []byte) ([]byte, error) {
   resp, err := http.Post(
      p.WvServer, "application/x-protobuf", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

// hulu.com/movie/05e76ad8-c3dd-4c3e-bab9-df3cf71c6871
// hulu.com/movie/alien-romulus-05e76ad8-c3dd-4c3e-bab9-df3cf71c6871
func Id(raw_url string) (string, error) {
   last_slash := strings.LastIndex(raw_url, "/")
   if last_slash == -1 {
      return "", errors.New("no slash found in URL")
   }
   last_part := raw_url[last_slash+1:]
   len_last := len(last_part)
   const len_uuid = 36
   if len_last > len_uuid {
      if last_part[len_last-len_uuid-1] == '-' {
         return last_part[len_last-len_uuid:], nil
      }
   }
   return last_part, nil
}

type Playlist struct {
   DashPrServer string `json:"dash_pr_server"`
   WvServer     string `json:"wv_server"`
   Message      string
   StreamUrl    string `json:"stream_url"` // MPD
}

type Authenticate struct {
   DeviceToken string `json:"device_token"`
   UserToken   string `json:"user_token"`
}

type DeepLink struct {
   EabId   string `json:"eab_id"`
   Message string
}

// returns user_token only
func (a *Authenticate) Refresh() error {
   resp, err := http.PostForm(
      "https://auth.hulu.com/v1/device/device_token/authenticate", url.Values{
         "action":       {"token_refresh"},
         "device_token": {a.DeviceToken},
      },
   )
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(a)
}

func (a Authenticate) DeepLink(id string) (*DeepLink, error) {
   req, _ := http.NewRequest("", "https://discover.hulu.com", nil)
   req.URL.Path = "/content/v5/deeplink/playback"
   req.URL.RawQuery = url.Values{
      "id":        {id},
      "namespace": {"entity"},
   }.Encode()
   req.Header.Set("authorization", "Bearer "+a.UserToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var deep DeepLink
   err = json.NewDecoder(resp.Body).Decode(&deep)
   if err != nil {
      return nil, err
   }
   if deep.Message != "" {
      return nil, errors.New(deep.Message)
   }
   return &deep, nil
}

var deejay = []struct {
   resolution  string
   device_id   int
   key_version int
}{
   {
      resolution:  "2160p",
      device_id:   210,
      key_version: 1,
   },
   {
      resolution:  "2160p",
      device_id:   208,
      key_version: 1,
   },
   {
      resolution:  "2160p",
      device_id:   204,
      key_version: 4,
   },
   {
      resolution:  "2160p",
      device_id:   188,
      key_version: 17,
   },
   {
      resolution:  "720p",
      device_id:   214,
      key_version: 1,
   },
   {
      resolution:  "720p",
      device_id:   191,
      key_version: 1,
   },
   {
      resolution:  "720p",
      device_id:   190,
      key_version: 1,
   },
   {
      resolution:  "720p",
      device_id:   142,
      key_version: 1,
   },
   {
      resolution:  "720p",
      device_id:   109,
      key_version: 1,
   },
}
