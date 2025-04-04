package roku

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "strings"
)

func (a *AccountToken) Code(act *Activation) (Byte[Code], error) {
   req, _ := http.NewRequest("", "https://googletv.web.roku.com", nil)
   req.URL.Path = "/api/v1/account/activation/" + act.Code
   // .Get .Set
   req.Header.Set("user-agent", user_agent)
   req.Header.Set("x-roku-content-token", a.AuthToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (a *AccountToken) Activation() (Byte[Activation], error) {
   data, err := json.Marshal(map[string]string{"platform": "googletv"})
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://googletv.web.roku.com/api/v1/account/activation",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   // .Get .Set
   req.Header.Set("content-type", "application/json")
   req.Header.Set("user-agent", user_agent)
   req.Header.Set("x-roku-content-token", a.AuthToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (a *AccountToken) Playback(roku_id string) (Byte[Playback], error) {
   data, err := json.Marshal(map[string]string{
      "mediaFormat": "DASH",
      "providerId":  "rokuavod",
      "rokuId":      roku_id,
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://googletv.web.roku.com/api/v3/playback",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   // .Set to match .Get
   req.Header.Set("content-type", "application/json")
   req.Header.Set("user-agent", user_agent)
   req.Header.Set("x-roku-content-token", a.AuthToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   return io.ReadAll(resp.Body)
}

type Byte[T any] []byte

// code can be nil
func (c *Code) AccountToken() (Byte[AccountToken], error) {
   req, _ := http.NewRequest("", "https://googletv.web.roku.com", nil)
   req.URL.Path = "/api/v1/account/token"
   req.Header.Set("user-agent", user_agent)
   if c != nil {
      req.Header.Set("x-roku-content-token", c.Token)
   }
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (a *AccountToken) Unmarshal(data Byte[AccountToken]) error {
   return json.Unmarshal(data, a)
}

const user_agent = "trc-googletv; production; 0"

type AccountToken struct {
   AuthToken string
}

func (a *Activation) String() string {
   var b strings.Builder
   b.WriteString("1 Visit the URL\n")
   b.WriteString("  therokuchannel.com/link\n")
   b.WriteString("\n")
   b.WriteString("2 Enter the activation code\n")
   b.WriteString("  ")
   b.WriteString(a.Code)
   return b.String()
}

func (a *Activation) Unmarshal(data Byte[Activation]) error {
   return json.Unmarshal(data, a)
}

type Activation struct {
   Code string
}

type Code struct {
   Token string
}

func (c *Code) Unmarshal(data Byte[Code]) error {
   return json.Unmarshal(data, c)
}

type Playback struct {
   Drm struct {
      Widevine struct {
         LicenseServer string
      }
   }
   Url string // MPD
}

func (p *Playback) Widevine(data []byte) ([]byte, error) {
   resp, err := http.Post(
      p.Drm.Widevine.LicenseServer, "application/x-protobuf",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (p *Playback) Unmarshal(data Byte[Playback]) error {
   return json.Unmarshal(data, p)
}
