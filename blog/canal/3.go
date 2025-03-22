package canal

import (
   "bytes"
   "encoding/json"
   "net/http"
)

func (t token) Play() (*http.Response, error) {
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
      "authorization": {"Bearer " + t.Token},
      "content-type": {"application/json"},
   }
   return http.DefaultClient.Do(req)
}
