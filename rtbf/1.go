package rtbf

import (
   "encoding/json"
   "errors"
   "net/http"
   "net/url"
)

type Login struct {
   ErrorMessage string
   SessionInfo  struct {
      CookieValue string
   }
}

func (l *Login) Fetch(id, password string) error {
   resp, err := http.PostForm(
      "https://login.auvio.rtbf.be/accounts.login", url.Values{
         "APIKey":   {api_key},
         "loginID":  {id},
         "password": {password},
      },
   )
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   err = json.NewDecoder(resp.Body).Decode(l)
   if err != nil {
      return err
   }
   if l.ErrorMessage != "" {
      return errors.New(l.ErrorMessage)
   }
   return nil
}
