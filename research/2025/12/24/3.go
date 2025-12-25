package disney

import (
   "bytes"
   "errors"
   "io"
   "net/http"
   "net/url"
)

func (r refresh_token) obtain_license(data []byte) ([]byte, error) {
   var req http.Request
   req.Header = http.Header{}
   req.Method = "POST"
   req.URL = &url.URL{}
   req.URL.Host = "disney.playback.edge.bamgrid.com"
   req.URL.Path = "/widevine/v1/obtain-license"
   req.URL.Scheme = "https"
   req.Body = io.NopCloser(bytes.NewReader(data))
   req.Header.Set("Authorization", r.Extensions.Sdk.Token.AccessToken)
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   return io.ReadAll(resp.Body)
}
