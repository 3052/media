package roku

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "strings"
)

// Widevine fetches the license key
func (p *Playback) Widevine(payload []byte) ([]byte, error) {
   resp, err := http.Post(
      p.Drm.Widevine.LicenseServer, "application/x-protobuf",
      bytes.NewReader(payload),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (p *Playback) Dash() (*Dash, error) {
   resp, err := http.Get(p.Url)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Dash
   result.Body, err = io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   result.Url = resp.Request.URL
   return &result, nil
}

type Dash struct {
   Body []byte
   Url  *url.URL
}

const user_agent = "trc-googletv; production; 0"

// Connection represents the active session with the Roku API.
type Connection struct {
   AuthToken string
}

// LinkCode represents the activation code displayed to the user.
type LinkCode struct {
   Code string
}

// User represents the persistent saved account token.
type User struct {
   Token string
}

// String generates the user instructions
func (l *LinkCode) String() string {
   var data strings.Builder
   data.WriteString("1 Visit the URL\n")
   data.WriteString("  therokuchannel.com/link\n")
   data.WriteString("\n")
   data.WriteString("2 Enter the activation code\n")
   data.WriteString("  ")
   data.WriteString(l.Code)
   return data.String()
}

// Playback represents the media playback metadata.
type Playback struct {
   Drm struct {
      Widevine struct {
         LicenseServer string
      }
   }
   Url string // MPD
}

// NewConnection initializes a session. User can be nil
func NewConnection(current *User) (*Connection, error) {
   var req http.Request
   req.Header = http.Header{}
   req.Header.Set("user-agent", user_agent)
   if current != nil {
      req.Header.Set("x-roku-content-token", current.Token)
   }
   req.URL = &url.URL{
      Scheme: "https",
      Host:   "googletv.web.roku.com",
      Path:   "/api/v1/account/token",
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   result := &Connection{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

// GetUser exchanges the LinkCode for a permanent User token.
func (c *Connection) User(link *LinkCode) (*User, error) {
   var req http.Request
   req.Header = http.Header{}
   req.Header.Set("user-agent", user_agent)
   req.Header.Set("x-roku-content-token", c.AuthToken)
   req.URL = &url.URL{
      Scheme: "https",
      Host:   "googletv.web.roku.com",
      Path:   "/api/v1/account/activation/" + link.Code,
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   result := &User{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

// RequestLinkCode requests a new activation code from the server
func (c *Connection) LinkCode() (*LinkCode, error) {
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
   req.Header.Set("x-roku-content-token", c.AuthToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   result := &LinkCode{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

// Playback fetches the DASH manifest and DRM information.
func (c *Connection) Playback(rokuId string) (*Playback, error) {
   payload, err := json.Marshal(map[string]string{
      "mediaFormat": "DASH",
      "providerId":  "rokuavod",
      "rokuId":      rokuId,
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://googletv.web.roku.com/api/v3/playback",
      bytes.NewReader(payload),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("content-type", "application/json")
   req.Header.Set("user-agent", user_agent)
   req.Header.Set("x-roku-content-token", c.AuthToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   result := &Playback{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}
