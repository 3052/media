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

type Mpd struct {
   Body []byte
   Url  *url.URL
}

func (p *Playlist) Mpd() (*Mpd, error) {
   resp, err := http.Get(p.StreamUrl)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   data, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   return &Mpd{data, resp.Request.URL}, nil
}

func (s *Session) Fetch(email, password string) error {
   resp, err := http.PostForm(
      "https://auth.hulu.com/v2/livingroom/password/authenticate", url.Values{
         "friendly_name": {"!"},
         "password":      {password},
         "serial_number": {"!"},
         "user_email":    {email},
      },
   )
   if err != nil {
      return err
   }
   if resp.StatusCode != http.StatusOK {
      var data strings.Builder
      err = resp.Write(&data)
      if err != nil {
         return err
      }
      return errors.New(data.String())
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(s)
}

// returns user_token only
func (s *Session) TokenRefresh() error {
   resp, err := http.PostForm(
      "https://auth.hulu.com/v1/device/device_token/authenticate", url.Values{
         "action":       {"token_refresh"},
         "device_token": {s.Data.DeviceToken},
      },
   )
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(s)
}

type Playlist struct {
   DashPrServer string `json:"dash_pr_server"`
   WvServer     string `json:"wv_server"`
   Message      string
   StreamUrl    string `json:"stream_url"` // MPD
}

type Session struct {
   Data struct {
      DeviceToken string `json:"device_token"`
      UserToken   string `json:"user_token"`
   }
}

func (s *Session) DeepLink(id string) (*DeepLink, error) {
   req, _ := http.NewRequest("", "https://discover.hulu.com", nil)
   req.URL.Path = "/content/v5/deeplink/playback"
   req.URL.RawQuery = url.Values{
      "id":        {id},
      "namespace": {"entity"},
   }.Encode()
   req.Header.Set("authorization", "Bearer "+s.Data.UserToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result DeepLink
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.Message != "" {
      return nil, errors.New(result.Message)
   }
   return &result, nil
}

// 1080p (FHD) L3, SL2000
// 1440p (QHD) L1, SL3000
// 2160p (UHD) L1, SL3000
func (s *Session) Playlist(deep *DeepLink) (*Playlist, error) {
   data, err := json.Marshal(map[string]any{
      "deejay_device_id": deejay[0].device_id,
      "version":          deejay[0].key_version,
      "content_eab_id":   deep.EabId,
      "unencrypted":      true,
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
            "multi_key":      true, // NEED THIS FOR 4K UHD
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
   req.Header.Set("authorization", "Bearer "+s.Data.UserToken)
   req.Header.Set("content-type", "application/json")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Playlist
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.Message != "" {
      return nil, errors.New(result.Message)
   }
   return &result, nil
}

type DeepLink struct {
   EabId   string `json:"eab_id"`
   Message string
}

// hulu.com/movie/05e76ad8-c3dd-4c3e-bab9-df3cf71c6871
// hulu.com/movie/alien-romulus-05e76ad8-c3dd-4c3e-bab9-df3cf71c6871
func Id(link string) (string, error) {
   slash := strings.LastIndex(link, "/")
   if slash == -1 {
      return "", errors.New("no slash found in URL")
   }
   part := link[slash+1:]
   len_part := len(part)
   const len_uuid = 36
   if len_part > len_uuid {
      if part[len_part-len_uuid-1] == '-' {
         return part[len_part-len_uuid:], nil
      }
   }
   return part, nil
}

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
      var result struct {
         Message string
      }
      err = json.Unmarshal(data, &result)
      if err != nil {
         return nil, err
      }
      return nil, errors.New(result.Message)
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
   data, err = io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   if resp.StatusCode != http.StatusOK {
      var result struct {
         Message string
      }
      err = json.Unmarshal(data, &result)
      if err != nil {
         return nil, err
      }
      return nil, errors.New(result.Message)
   }
   return data, nil
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
