package binge

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
)

func (t *token) service() error {
   data, err := json.Marshal(map[string]string{"client_id": client_id})
   if err != nil {
      return err
   }
   req, err := http.NewRequest(
      "POST", "https://tokenservice.streamotion.com.au/oauth/token",
      bytes.NewReader(data),
   )
   if err != nil {
      return err
   }
   req.Header = http.Header{
      "content-type": {"application/json"},
      "authorization": {"Bearer " + t.AccessToken},
   }
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(t)
}

func new_token(username, password string) (Byte[token], error) {
   data, err := json.Marshal( map[string]string{
      "client_id": client_id,
      "grant_type": "http://auth0.com/oauth/grant-type/password-realm",
      "password": password,
      "realm": "prod-martian-database",
      "username": username,
      // need this otherwise service request fails
      "audience": "streamotion.com.au",
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

// android?
//const client_id = "QQdtPlVtx1h9BkO09BDM2OrFi5vTPCty"

// web
const client_id = "pM87TUXKQvSSu93ydRjDTqBgdYeCbdhZ"

// new refresh token is not returned, so we can keep old
func (t *token) refresh() error {
   data, err := json.Marshal(map[string]string{
      "client_id": client_id,
      "grant_type": "refresh_token",
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
   err = json.NewDecoder(resp.Body).Decode(t)
   if err != nil {
      return err
   }
   if t.ErrorDescription != "" {
      return errors.New(t.ErrorDescription)
   }
   return nil
}

type token struct {
   AccessToken string `json:"access_token"`
   ErrorDescription string `json:"error_description"`
   IdToken string `json:"id_token"`
   RefreshToken string `json:"refresh_token"`
}

type Byte[T any] []byte

func (t *token) unmarshal(data Byte[token]) error {
   return json.Unmarshal(data, t)
}
