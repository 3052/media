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

type Cache struct {
   Connection *Connection
   LinkCode   *LinkCode
   Mpd        *url.URL
   MpdBody    []byte
   User       *User
}

func (p *Playback) Mpd(storage *Cache) error {
   resp, err := http.Get(p.Url)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   storage.MpdBody, err = io.ReadAll(resp.Body)
   if err != nil {
      return err
   }
   storage.Mpd = resp.Request.URL
   return nil
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

// GetUser exchanges the LinkCode for a permanent User token.
func (c *Connection) GetUser(link *LinkCode) (*User, error) {
   req, _ := http.NewRequest("", "https://googletv.web.roku.com", nil)
   req.URL.Path = "/api/v1/account/activation/" + link.Code
   req.Header.Set("user-agent", userAgent)
   req.Header.Set("x-roku-content-token", c.AuthToken)

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   currentUser := &User{}
   err = json.NewDecoder(resp.Body).Decode(currentUser)
   if err != nil {
      return nil, err
   }
   return currentUser, nil
}

const userAgent = "trc-googletv; production; 0"

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
   req.Header.Set("user-agent", userAgent)
   req.Header.Set("x-roku-content-token", c.AuthToken)

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }

   playResult := &Playback{}
   err = json.NewDecoder(resp.Body).Decode(playResult)
   if err != nil {
      return nil, err
   }
   return playResult, nil
}

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

// String generates the user instructions
func (l *LinkCode) String() string {
   var output strings.Builder
   output.WriteString("1 Visit the URL\n")
   output.WriteString("  therokuchannel.com/link\n")
   output.WriteString("\n")
   output.WriteString("2 Enter the activation code\n")
   output.WriteString("  ")
   output.WriteString(l.Code)
   return output.String()
}

// RequestLinkCode requests a new activation code from the server
func (c *Connection) RequestLinkCode() (*LinkCode, error) {
   payload, err := json.Marshal(map[string]string{"platform": "googletv"})
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://googletv.web.roku.com/api/v1/account/activation",
      bytes.NewReader(payload),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("content-type", "application/json")
   req.Header.Set("user-agent", userAgent)
   req.Header.Set("x-roku-content-token", c.AuthToken)

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   link := &LinkCode{}
   err = json.NewDecoder(resp.Body).Decode(link)
   if err != nil {
      return nil, err
   }
   return link, nil
}

// NewConnection initializes a session. User can be nil
func (u *User) NewConnection() (*Connection, error) {
   req, _ := http.NewRequest("", "https://googletv.web.roku.com", nil)
   req.URL.Path = "/api/v1/account/token"
   req.Header.Set("user-agent", userAgent)
   if u != nil {
      req.Header.Set("x-roku-content-token", u.Token)
   }

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   connector := &Connection{}
   err = json.NewDecoder(resp.Body).Decode(connector)
   if err != nil {
      return nil, err
   }
   return connector, nil
}
