package molotov

import (
   "bytes"
   "encoding/json"
   "net/http"
)

const molotov_agent = `{ "app_build": 4, "app_id": "browser_app" }`

func zero(email, password string) (*http.Response, error) {
   data, err := json.Marshal(map[string]string{
      "grant_type": "password",
      "email": email,
      "password": password,
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://fapi.molotov.tv/v3.1/auth/login",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("x-molotov-agent", molotov_agent)
   return http.DefaultClient.Do(req)
}
