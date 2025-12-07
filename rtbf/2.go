package rtbf

import (
   "encoding/json"
   "errors"
   "net/http"
   "net/url"
)

type Jwt struct {
   ErrorMessage string
   IdToken      string `json:"id_token"`
}

func (l *Login) Jwt() (*Jwt, error) {
   resp, err := http.PostForm(
      "https://login.auvio.rtbf.be/accounts.getJWT", url.Values{
         "APIKey":      {api_key},
         "login_token": {l.SessionInfo.CookieValue},
      },
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var token Jwt
   err = json.NewDecoder(resp.Body).Decode(&token)
   if err != nil {
      return nil, err
   }
   if token.ErrorMessage != "" {
      return nil, errors.New(token.ErrorMessage)
   }
   return &token, nil
}
