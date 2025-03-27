package molotov

import (
   "net/http"
   "net/url"
)

func (r *refresh) assets(view1 *view) (*http.Response, error) {
   req, _ := http.NewRequest("", "https://fapi.molotov.tv/v2/me/assets", nil)
   req.URL.RawQuery = url.Values{
      "access_token": {r.AccessToken},
      "id": {view1.Program.Video.Id},
      "type": {"vod"},
   }.Encode()
   req.Header = http.Header{
      "x-forwarded-for": {"138.199.15.158"},
      "x-molotov-agent": {molotov_agent},
   }
   return http.DefaultClient.Do(req)
}
