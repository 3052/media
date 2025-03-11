package binge

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
)

func (a *auth) token() (*token_service, error) {
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
   req.Header = http.Header{
      "content-type": {"application/json"},
      "authorization": {"Bearer " + a.AccessToken},
   }
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   token := &token_service{}
   err = json.NewDecoder(resp.Body).Decode(token)
   if err != nil {
      return nil, err
   }
   return token, nil
}

// android?
//const client_id = "QQdtPlVtx1h9BkO09BDM2OrFi5vTPCty"

// web
const client_id = "pM87TUXKQvSSu93ydRjDTqBgdYeCbdhZ"

// new refresh token is not returned, so we can keep old
func (a *auth) refresh() error {
   data, err := json.Marshal(map[string]string{
      "client_id": client_id,
      "grant_type": "refresh_token",
      "refresh_token": a.RefreshToken,
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
   err = json.NewDecoder(resp.Body).Decode(a)
   if err != nil {
      return err
   }
   if a.ErrorDescription != "" {
      return errors.New(a.ErrorDescription)
   }
   return nil
}

type Byte[T any] []byte

func (a *auth) unmarshal(data Byte[auth]) error {
   return json.Unmarshal(data, a)
}

func new_auth(username, password string) (Byte[auth], error) {
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

type auth struct {
   AccessToken string `json:"access_token"`
   ErrorDescription string `json:"error_description"`
   IdToken string `json:"id_token"`
   RefreshToken string `json:"refresh_token"`
}

type token_service struct {
   AccessToken string `json:"access_token"`
}
