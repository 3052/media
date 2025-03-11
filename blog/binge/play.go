package binge

import (
   "bytes"
   "encoding/json"
   "net/http"
)

func (t token_service) play() (*http.Response, error) {
   data, err := json.Marshal(map[string]any{
      "assetId": "7738",
      "application": map[string]string{
         "name": "binge",
      },
      "device": map[string]string{
         "id": "50e785be-4c7f-4781-87e4-a3b4c75a3634",
      },
      "player": map[string]string{
         "name": "VideoFS",
      },
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://play.binge.com.au/api/v3/play", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header = http.Header{
      "content-type": {"application/json"},
      "authorization": {"Bearer " + t.AccessToken},
   }
   return http.DefaultClient.Do(req)
}
