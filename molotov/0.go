package molotov

import (
   "bytes"
   "encoding/json"
   "net/http"
)

const molotov_agent = `{ "app_build": 4, "app_id": "browser_app" }`

type login struct {
   RefreshToken string `json:"refresh_token"`
}

func (n *login) New(email, password string) error {
   data, err := json.Marshal(map[string]string{
      "grant_type": "password",
      "email": email,
      "password": password,
   })
   if err != nil {
      return err
   }
   req, err := http.NewRequest(
      "POST", "https://fapi.molotov.tv/v3.1/auth/login",
      bytes.NewReader(data),
   )
   if err != nil {
      return err
   }
   req.Header.Set("x-molotov-agent", molotov_agent)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(n)
}
