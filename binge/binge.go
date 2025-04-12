package binge

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "strconv"
)

func (a *Auth) Unmarshal(data Byte[Auth]) error {
   err := json.Unmarshal(data, a)
   if err != nil {
      return err
   }
   if a.ErrorDescription != "" {
      return errors.New(a.ErrorDescription)
   }
   return nil
}

func NewAuth(username, password string) (Byte[Auth], error) {
   data, err := json.Marshal(map[string]string{
      "client_id":  client_id,
      "grant_type": "http://auth0.com/oauth/grant-type/password-realm",
      "password":   password,
      "realm":      "prod-martian-database",
      "username":   username,
      // need this otherwise service request fails
      "audience": "streamotion.com.au",
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://auth.streamotion.com.au/oauth/token",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("content-type", "application/json")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (p *Play) Unmarshal(data Byte[Play]) error {
   return json.Unmarshal(data, p)
}

type Play struct {
   Streams []Stream
}

type Stream struct {
   LicenseAcquisitionUrl struct {
      ComWidevineAlpha string `json:"com.widevine.alpha"`
   }
   Manifest        string // MPD
   Provider        string
   StreamingFormat string
}

func (t TokenService) Play(asset_id int) (Byte[Play], error) {
   data, err := json.Marshal(map[string]any{
      "application": map[string]string{
         "name": "binge",
      },
      "assetId": strconv.Itoa(asset_id),
      "device": map[string]string{
         "id": "!",
      },
      "player": map[string]string{
         "name": "VideoFS",
      },
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://play.binge.com.au/api/v3/play", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer " + t.AccessToken)
   req.Header.Set("content-type", "application/json")
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

func (t TokenService) Widevine(stream1 *Stream, data []byte) ([]byte, error) {
   req, err := http.NewRequest(
      "POST", stream1.LicenseAcquisitionUrl.ComWidevineAlpha,
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+t.AccessToken)
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

func (t *TokenService) Unmarshal(data Byte[TokenService]) error {
   return json.Unmarshal(data, t)
}

type TokenService struct {
   AccessToken string `json:"access_token"`
}
// SEGMENTS ARE GEO BLOCK WITH ALL PROVIDER
func (p Play) Dash() (*Stream, bool) {
   for _, stream1 := range p.Streams {
      if stream1.StreamingFormat == "dash" {
         return &stream1, true
      }
   }
   return nil, false
}

// web
const client_id = "pM87TUXKQvSSu93ydRjDTqBgdYeCbdhZ"

func (a *Auth) Token() (Byte[TokenService], error) {
   data, err := json.Marshal(map[string]string{"client_id": client_id})
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://tokenservice.streamotion.com.au/oauth/token",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer " + a.AccessToken)
   req.Header.Set("content-type", "application/json")
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

func (a *Auth) Refresh() (Byte[Auth], error) {
   data, err := json.Marshal(map[string]string{
      "client_id":     client_id,
      "grant_type":    "refresh_token",
      "refresh_token": a.RefreshToken,
   })
   if err != nil {
      return nil, err
   }
   resp, err := http.Post(
      "https://auth.streamotion.com.au/oauth/token", "application/json",
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

type Auth struct {
   AccessToken      string `json:"access_token"`
   ErrorDescription string `json:"error_description"`
   IdToken          string `json:"id_token"`
   RefreshToken     string `json:"refresh_token"`
}

type Byte[T any] []byte
