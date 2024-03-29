package roku

import (
   "bytes"
   "encoding/json"
   "errors"
   "net/http"
)

func (c CrossSite) csrf() (*http.Cookie, bool) {
   for _, cookie := range c.cookies {
      if cookie.Name == "_csrf" {
         return cookie, true
      }
   }
   return nil, false
}

func (Playback) RequestHeader() (http.Header, error) {
   return http.Header{}, nil
}

func (p Playback) RequestUrl() (string, bool) {
   return p.DRM.Widevine.LicenseServer, true
}

func (c CrossSite) Playback(id string) (*Playback, error) {
   csrf, ok := c.csrf()
   if !ok {
      return nil, http.ErrNoCookie
   }
   body, err := func() ([]byte, error) {
      m := map[string]string{
         "mediaFormat": "mpeg-dash",
         "providerId": "rokuavod",
         "rokuId": id,
      }
      return json.Marshal(m)
   }()
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://therokuchannel.roku.com/api/v3/playback",
      bytes.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   // we could use Request.AddCookie, but we would need to call it after this,
   // otherwise it would be clobbered
   req.Header = http.Header{
      "CSRF-Token": {c.token},
      "Content-Type": {"application/json"},
      "Cookie": {csrf.Raw},
   }
   res, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return nil, errors.New(res.Status)
   }
   play := new(Playback)
   if err := json.NewDecoder(res.Body).Decode(play); err != nil {
      return nil, err
   }
   return play, nil
}

type Playback struct {
   DRM struct {
      Widevine struct {
         LicenseServer string
      }
   }
}

func (Playback) RequestBody(b []byte) ([]byte, error) {
   return b, nil
}

func (Playback) ResponseBody(b []byte) ([]byte, error) {
   return b, nil
}

type CrossSite struct {
   cookies []*http.Cookie
   token string
}

type MediaVideo struct {
   DrmAuthentication *struct{}
   URL string
   VideoType string
}
