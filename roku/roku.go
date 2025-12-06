package roku

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "strings"
)

func (a *Activation) String() string {
   var data strings.Builder
   data.WriteString("1 Visit the URL\n")
   data.WriteString("  therokuchannel.com/link\n")
   data.WriteString("\n")
   data.WriteString("2 Enter the activation code\n")
   data.WriteString("  ")
   data.WriteString(a.Code)
   return data.String()
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

const user_agent = "trc-googletv; production; 0"

type AccountToken struct {
   AuthToken string
}

// code can be nil
func (c *Code) AccountToken() (*AccountToken, error) {
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
   token := &AccountToken{}
   err = json.NewDecoder(resp.Body).Decode(token)
   if err != nil {
      return nil, err
   }
   return token, nil
}

type Activation struct {
   Code string
}

func (a *AccountToken) Activation() (*Activation, error) {
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
   req.Header.Set("content-type", "application/json")
   req.Header.Set("user-agent", user_agent)
   req.Header.Set("x-roku-content-token", a.AuthToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   value := &Activation{}
   err = json.NewDecoder(resp.Body).Decode(value)
   if err != nil {
      return nil, err
   }
   return value, nil
}

type Code struct {
   Token string
}

func (a *AccountToken) Code(act *Activation) (*Code, error) {
   req, _ := http.NewRequest("", "https://googletv.web.roku.com", nil)
   req.URL.Path = "/api/v1/account/activation/" + act.Code
   req.Header.Set("user-agent", user_agent)
   req.Header.Set("x-roku-content-token", a.AuthToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   value := &Code{}
   err = json.NewDecoder(resp.Body).Decode(value)
   if err != nil {
      return nil, err
   }
   return value, nil
}

func (a *AccountToken) Playback(rokuId string) (*Playback, error) {
   data, err := json.Marshal(map[string]string{
      "mediaFormat": "DASH",
      "providerId":  "rokuavod",
      "rokuId":      rokuId,
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
   play := &Playback{}
   err = json.NewDecoder(resp.Body).Decode(play)
   if err != nil {
      return nil, err
   }
   return play, nil
}
