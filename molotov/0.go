package molotov

import (
   "bytes"
   "encoding/json"
   "io"
   "net/http"
)

func (n *login) unmarshal(data Byte[login]) error {
   return json.Unmarshal(data, n)
}

type login struct {
   RefreshToken string `json:"refresh_token"`
}

const molotov_agent = `{ "app_build": 4, "app_id": "browser_app" }`

type Byte[T any] []byte

func new_login(email, password string) (Byte[login], error) {
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
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}
