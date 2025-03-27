package molotov

import (
   "bytes"
   "encoding/json"
   "io"
   "net/http"
)

const molotov_agent = `{ "app_build": 4, "app_id": "browser_app" }`

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

// authorization server issues a new refresh token, in which case the
// client MUST discard the old refresh token and replace it with the new
// refresh token
func (r *refresh) refresh() (Byte[refresh], error) {
   req, _ := http.NewRequest("", "https://fapi.molotov.tv", nil)
   req.URL.Path = "/v3/auth/refresh/" + r.RefreshToken
   req.Header.Set("x-molotov-agent", molotov_agent)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

type Byte[T any] []byte

type refresh struct {
   AccessToken string `json:"access_token"`
   RefreshToken string `json:"refresh_token"`
}

type login struct {
   Auth refresh
}

func (r *refresh) unmarshal(data Byte[refresh]) error {
   return json.Unmarshal(data, r)
}
