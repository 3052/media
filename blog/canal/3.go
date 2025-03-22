package canal

import (
   "bytes"
   "encoding/json"
   "net/http"
)

type play struct {
   Url string
}

// us fail
// czech republic mullvad fail
// czech republic nord fail
// czech republic smart proxy pass
func (s session) play() (*play, error) {
   data, err := json.Marshal(map[string]any{
      "player": map[string]any{
         "capabilities": map[string]any{
            "drmSystems": []string{"Widevine"},
            "mediaTypes": []string{"DASH"},
         },
      },
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://tvapi-hlm2.solocoo.tv", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.URL.Path = "/v1/assets/1EBvrU5Q2IFTIWSC2_4cAlD98U0OR0ejZm_dgGJi/play"
   req.Header = http.Header{
      "authorization": {"Bearer " + s.Token},
      "content-type": {"application/json"},
   }
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   play1 := &play{}
   err = json.NewDecoder(resp.Body).Decode(play1)
   if err != nil {
      return nil, err
   }
   return play1, nil
}
