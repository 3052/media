package movistar

import (
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
)

// XFF fail
// mullvad pass
// nord pass
func new_token(username, password string) (Byte[token], error) {
   resp, err := http.PostForm(
      "https://auth.dof6.com/auth/oauth2/token?deviceClass=amazon.tv",
      url.Values{
         "grant_type": {"password"},
         "password":   {password},
         "username":   {username},
      },
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

// 10 days
type token struct {
   AccessToken string `json:"access_token"`
   ExpiresIn   int64  `json:"expires_in"`
}

type Byte[T any] []byte

func (t *token) unmarshal(data Byte[token]) error {
   return json.Unmarshal(data, t)
}
