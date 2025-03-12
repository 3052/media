package binge

import (
   "bytes"
   "io"
   "net/http"
)

func (t token_service) widevine(stream1 *stream, data []byte) ([]byte, error) {
   req, err := http.NewRequest(
      "POST", stream1.LicenseAcquisitionUrl.ComWidevineAlpha,
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer " + t.AccessToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}
