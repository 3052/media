package molotov

import "net/http"

func (r *refresh) view() (*http.Response, error) {
   req, err := http.NewRequest(
      "", "https://fapi.molotov.tv/v2/channels/531/programs/15082/view", nil,
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("x-molotov-agent", molotov_agent)
   req.URL.RawQuery = "access_token=" + r.AccessToken
   return http.DefaultClient.Do(req)
}
