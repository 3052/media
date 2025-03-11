package binge

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
)

// new refresh token is not returned, so we can keep old
func (t *token) refresh() error {
   data, err := json.Marshal(map[string]string{
      "grant_type": "refresh_token",
      "scope": "openid offline_access drm:low email",
      "client_id": "QQdtPlVtx1h9BkO09BDM2OrFi5vTPCty",
      "refresh_token": t.RefreshToken,
   })
   if err != nil {
      return err
   }
   resp, err := http.Post(
      "https://auth.streamotion.com.au/oauth/token", "application/json",
      bytes.NewReader(data),
   )
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(t)
}

type token struct {
   AccessToken string `json:"access_token"`
   IdToken string `json:"id_token"`
   RefreshToken string `json:"refresh_token"`
}

type Byte[T any] []byte

func new_token(username, password string) (Byte[token], error) {
   data, err := json.Marshal( map[string]string{
      "client_id": "QQdtPlVtx1h9BkO09BDM2OrFi5vTPCty",
      "grant_type": "http://auth0.com/oauth/grant-type/password-realm",
      "password": password,
      "realm": "prod-martian-database",
      "username": username,
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

func (t *token) unmarshal(data Byte[token]) error {
   return json.Unmarshal(data, t)
}
