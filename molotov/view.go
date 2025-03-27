package molotov

import (
   "encoding/json"
   "net/http"
)

type view struct {
   Program struct {
      Video struct {
         Id string
      }
   }
}

func (r *refresh) view() (*view, error) {
   req, err := http.NewRequest(
      "", "https://fapi.molotov.tv/v2/channels/531/programs/15082/view", nil,
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("x-molotov-agent", molotov_agent)
   req.URL.RawQuery = "access_token=" + r.AccessToken
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   view1 := &view{}
   err = json.NewDecoder(resp.Body).Decode(view1)
   if err != nil {
      return nil, err
   }
   return view1, nil
}
