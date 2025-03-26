package max

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "strings"
)

func (n *Login) Playback(edit_id string) (Byte[Playback], error) {
   value := map[string]any{
      "editId": edit_id,
      "consumptionType":      "streaming",
      "appBundle":            "",         // required
      "applicationSessionId": "",         // required
      "firstPlay":            false,      // required
      "gdpr":                 false,      // required
      "playbackSessionId":    "",         // required
      "userPreferences":      struct{}{}, // required
      "capabilities": map[string]any{
         "manifests": map[string]any{
            "formats": map[string]any{
               "dash": struct{}{}, // required
            }, // required
         }, // required
      }, // required
      "deviceInfo": map[string]any{
         "player": map[string]any{
            "mediaEngine": map[string]string{
               "name":    "", // required
               "version": "", // required
            }, // required
            "playerView": map[string]int{
               "height": 0, // required
               "width":  0, // required
            }, // required
            "sdk": map[string]string{
               "name":    "", // required
               "version": "", // required
            }, // required
         }, // required
      }, // required
   }
   data, err := json.MarshalIndent(value, "", " ")
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest("POST", prd_api, bytes.NewReader(data))
   if err != nil {
      return nil, err
   }
   req.URL.Path = func() string {
      var b bytes.Buffer
      b.WriteString("/playback-orchestrator/any/playback-orchestrator/v1")
      b.WriteString("/playbackInfo")
      return b.String()
   }()
   // .Set to match .Get
   req.Header.Set("content-type", "application/json")
   req.Header.Set("authorization", "Bearer "+n.Data.Attributes.Token)
   req.Header.Set("proxy", "true")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

const (
   device_info  = "!/!(!/!;!/!;!/!)"
   disco_client = "!:!:beam:!"
   prd_api      = "https://default.prd.api.discomax.com"
)

type Playback struct {
   Drm struct {
      Schemes struct {
         Widevine struct {
            LicenseUrl string
         }
      }
   }
   Errors []struct {
      Detail string
   }
   Fallback struct {
      Manifest struct {
         Url Url // MPD
      }
   }
}

func (p *Playback) Unmarshal(data Byte[Playback]) error {
   err := json.Unmarshal(data, p)
   if err != nil {
      return err
   }
   if len(p.Errors) >= 1 {
      return errors.New(p.Errors[0].Detail)
   }
   return nil
}

func (i *Initiate) String() string {
   var b strings.Builder
   b.WriteString("target URL = ")
   b.WriteString(i.TargetUrl)
   b.WriteString("\nlinking code = ")
   b.WriteString(i.LinkingCode)
   return b.String()
}

type Initiate struct {
   LinkingCode string
   TargetUrl   string
}

func (p *Playback) Widevine(data []byte) ([]byte, error) {
   resp, err := http.Post(
      p.Drm.Schemes.Widevine.LicenseUrl, "application/x-protobuf",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   return io.ReadAll(resp.Body)
}

func (s *St) New() error {
   req, _ := http.NewRequest("", prd_api+"/token?realm=bolt", nil)
   req.Header = http.Header{
      "x-device-info":  {device_info},
      "x-disco-client": {disco_client},
   }
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   for _, cookie := range resp.Cookies() {
      if cookie.Name == "st" {
         (*s)[0] = cookie
         return nil
      }
   }
   return http.ErrNoCookie
}

type St [1]*http.Cookie

func (s *St) Set(data string) error {
   var err error
   (*s)[0], err = http.ParseSetCookie(data)
   if err != nil {
      return err
   }
   return nil
}

func (s St) String() string {
   return s[0].String()
}

type Login struct {
   Data struct {
      Attributes struct {
         Token string
      }
   }
}

type Url [1]string

func (u *Url) UnmarshalText(data []byte) error {
   (*u)[0] = strings.Replace(string(data), "_fallback", "", 1)
   return nil
}

// you must
// /authentication/linkDevice/initiate
// first or this will always fail
func (s St) Login() (Byte[Login], error) {
   req, _ := http.NewRequest("POST", prd_api, nil)
   req.URL.Path = "/authentication/linkDevice/login"
   req.AddCookie(s[0])
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (n *Login) Unmarshal(data Byte[Login]) error {
   return json.Unmarshal(data, n)
}

type Byte[T any] []byte

func (s St) Initiate() (*Initiate, error) {
   req, _ := http.NewRequest("POST", prd_api, nil)
   req.URL.Path = "/authentication/linkDevice/initiate"
   req.AddCookie(s[0])
   req.Header.Set("x-device-info", device_info)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var value struct {
      Data struct {
         Attributes Initiate
      }
      Errors []struct {
         Detail string
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   if len(value.Errors) >= 1 {
      return nil, errors.New(value.Errors[0].Detail)
   }
   return &value.Data.Attributes, nil
}
