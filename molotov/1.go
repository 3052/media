package molotov

import (
   "encoding/json"
   "io"
   "net/http"
)

func (r *refresh) unmarshal(data Byte[refresh]) error {
   return json.Unmarshal(data, r)
}

type Byte[T any] []byte

type refresh struct {
   AccessToken string `json:"access_token"`
   RefreshToken string `json:"refresh_token"`
}

// authorization server issues a new refresh token, in which case the
// client MUST discard the old refresh token and replace it with the new
// refresh token
func (n login) refresh() (Byte[refresh], error) {
   req, _ := http.NewRequest("", "https://fapi.molotov.tv", nil)
   req.URL.Path = "/v3/auth/refresh/" + n.RefreshToken
   req.Header.Set("x-molotov-agent", molotov_agent)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}
