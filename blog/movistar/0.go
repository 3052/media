package movistar

import (
   "encoding/json"
   "net/http"
   "net/url"
   "time"
)

// XFF fail
// mullvad pass
// nord pass
func (t *token) New(username, password string) error {
   resp, err := http.PostForm(
      "https://auth.dof6.com/auth/oauth2/token?deviceClass=amazon.tv",
      url.Values{
         "grant_type": {"password"},
         "password": {password},
         "username": {username},
      },
   )
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(t)
}

// 10 days
type token struct {
   AccessToken string `json:"access_token"`
   ExpiresIn int64 `json:"expires_in"`
}

func (t *token) duration() time.Duration {
   return time.Duration(t.ExpiresIn) * time.Second
}
