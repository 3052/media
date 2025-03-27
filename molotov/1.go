package molotov

import "net/http"

// new refresh token is returned
func (n login) refresh() (*http.Response, error) {
   req, _ := http.NewRequest("", "https://fapi.molotov.tv", nil)
   req.URL.Path = "/v3/auth/refresh/" + n.RefreshToken
   req.Header.Set("x-molotov-agent", molotov_agent)
   return http.DefaultClient.Do(req)
}
