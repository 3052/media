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

// this is old device that returns 4K MPD:
// https://vodmanifest.hulustream.com
// newer devices return 2K MPD:
// https://dynamic-manifest.hulustream.com
const (
   //deejay_device_id = 166
   //version          = 9999999
   
   deejay_device_id = 204
   version          = 4
)

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

func (a Authenticate) Playlist(deep *DeepLink) (Byte[Playlist], error) {
   data, err := json.Marshal(map[string]any{
      "content_eab_id":   deep.EabId,
      "deejay_device_id": deejay_device_id,
      "unencrypted": true,
      "version":     version,
      "playback": map[string]any{
         "audio": map[string]any{
            "codecs": map[string]any{
               "selection_mode": "ALL",
               "values": []any{
                  map[string]string{"type": "AAC"},
                  map[string]string{"type": "EC3"},
               },
            },
         },
         "drm": map[string]any{
            // 1080p (FHD) L3
            // 1440p (QHD) L1
            // 2160p (UHD) L1
            "multi_key": true,
            "selection_mode": "ALL",
            "values": []any{
               map[string]string{
                  "security_level": "L3",
                  "type":           "WIDEVINE",
                  "version":        "MODULAR",
               },
               map[string]string{
                  "security_level": "SL2000",
                  "type":           "PLAYREADY",
                  "version":        "V2",
               },
            },
         },
         "version": 2, // needs to be exactly 2 for 1080p
         "manifest": map[string]string{
            "type": "DASH",
         },
         "segments": map[string]any{
            "selection_mode": "ALL",
            "values": []any{
               map[string]any{
                  "type": "FMP4",
                  "encryption": map[string]string{
                     "mode": "CENC",
                     "type": "CENC",
                  },
               },
            },
         },
         "video": map[string]any{
            "codecs": map[string]any{
               "selection_mode": "ALL",
               "values": []any{
                  map[string]any{
                     "height":  9999,
                     "level":   "9",
                     "profile": "HIGH",
                     "type":    "H264",
                     "width":   9999,
                  },
                  map[string]any{
                     "height":  9999,
                     "level":   "9",
                     "profile": "MAIN_10",
                     "tier":    "MAIN",
                     "type":    "H265",
                     "width":   9999,
                  },
               },
            },
         },
      },
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://play.hulu.com/v6/playlist", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+a.UserToken)
   req.Header.Set("content-type", "application/json")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (p *Playlist) License(data []byte) ([]byte, error) {
   resp, err := http.Post(
      p.WvServer, "application/x-protobuf", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

type Authenticate struct {
   DeviceToken string `json:"device_token"`
   UserToken   string `json:"user_token"`
}

type Byte[T any] []byte

type DeepLink struct {
   EabId   string `json:"eab_id"`
   Message string
}

func NewAuthenticate(email, password string) (Byte[Authenticate], error) {
   resp, err := http.PostForm(
      "https://auth.hulu.com/v2/livingroom/password/authenticate", url.Values{
         "friendly_name": {"!"},
         "password":      {password},
         "serial_number": {"!"},
         "user_email":    {email},
      },
   )
   if err != nil {
      return nil, err
   }
   if resp.StatusCode != http.StatusOK {
      var data strings.Builder
      resp.Write(&data)
      return nil, errors.New(data.String())
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (a *Authenticate) Unmarshal(data Byte[Authenticate]) error {
   var value struct {
      Data Authenticate
   }
   err := json.Unmarshal(data, &value)
   if err != nil {
      return err
   }
   *a = value.Data
   return nil
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

func (p *Playlist) Unmarshal(data Byte[Playlist]) error {
   err := json.Unmarshal(data, p)
   if err != nil {
      return err
   }
   if p.Message != "" {
      return errors.New(p.Message)
   }
   return nil
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
