package hulu

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "path"
   "strings"
)

type Playlist struct {
   Message   string
   StreamUrl string `json:"stream_url"` // MPD
   WvServer  string `json:"wv_server"`
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

func (a Authenticate) Playlist(deep *DeepLink) (Byte[Playlist], error) {
   data, err := json.Marshal(map[string]any{
      "content_eab_id":   deep.EabId,
      "deejay_device_id": 166,
      "unencrypted":      true,
      "version":          9999999,
      "playback": map[string]any{
         "version": 2, // needs to be exactly 2 for 1080p
         "manifest": map[string]string{
            "type": "DASH",
         },
         "drm": map[string]any{
            "selection_mode": "ALL",
            "values": []map[string]string{
               {
                  "security_level": "L3",
                  "type":           "WIDEVINE",
                  "version":        "MODULAR",
               },
            },
         },
         "segments": map[string]any{
            "selection_mode": "ALL",
            "values": []map[string]any{
               {
                  "type": "FMP4",
                  "encryption": map[string]string{
                     "mode": "CENC",
                     "type": "CENC",
                  },
               },
            },
         },
         "audio": map[string]any{
            "codecs": map[string]any{
               "selection_mode": "ALL",
               "values": []map[string]string{
                  {"type": "AAC"},
                  {"type": "EC3"},
               },
            },
         },
         "video": map[string]any{
            "codecs": map[string]any{
               "selection_mode": "ALL",
               "values": []map[string]any{
                  {
                     "height":  9999,
                     "level":   "9",
                     "profile": "HIGH",
                     "type":    "H264",
                     "width":   9999,
                  },
                  {
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
   req.Header = http.Header{
      "authorization": {"Bearer " + a.UserToken},
      "content-type":  {"application/json"},
   }
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

// hulu.com/watch/023c49bf-6a99-4c67-851c-4c9e7609cc1d
func (e *Entity) Set(data string) error {
   e[0] = path.Base(data)
   return nil
}

type Entity [1]string

func (e Entity) String() string {
   return e[0]
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
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      var data strings.Builder
      resp.Write(&data)
      return nil, errors.New(data.String())
   }
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

func (a Authenticate) DeepLink(id Entity) (*DeepLink, error) {
   req, _ := http.NewRequest("", "https://discover.hulu.com", nil)
   req.URL.Path = "/content/v5/deeplink/playback"
   req.URL.RawQuery = url.Values{
      "id":        {id[0]},
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

type Authenticate struct {
   DeviceToken string `json:"device_token"`
   UserToken   string `json:"user_token"`
}
