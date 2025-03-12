package binge

import (
   "bytes"
   "encoding/json"
   "net/http"
)

// Akamai needs residential proxy (or Nord) and CloudFront/Fastly work with
// just Mullvad
func (p play) dash() (*stream, bool) {
   for _, stream1 := range p.Streams {
      if stream1.StreamingFormat == "dash" {
         return &stream1, true
      }
   }
   return nil, false
}

type stream struct {
   LicenseAcquisitionUrl struct {
      ComWidevineAlpha string `json:"com.widevine.alpha"`
   }
   Manifest string
   Provider string
   StreamingFormat string
}

type play struct {
   Streams []stream
}

func (t token_service) play() (*play, error) {
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

