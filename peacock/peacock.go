package peacock

import (
   "bytes"
   "crypto/hmac"
   "crypto/md5"
   "crypto/sha1"
   "encoding/base64"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "slices"
   "strings"
   "time"
)

func FetchIdSession(user, password string) (*http.Cookie, error) {
   data := url.Values{
      "userIdentifier": {user},
      "password":       {password},
   }.Encode()
   req, err := http.NewRequest(
      "POST", "https://rango.id.peacocktv.com/signin/service/international",
      strings.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("content-type", "application/x-www-form-urlencoded")
   req.Header.Set("x-skyott-proposition", "NBCUOTT")
   req.Header.Set("x-skyott-provider", "NBCU")
   req.Header.Set("x-skyott-territory", Territory)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusCreated {
      var data strings.Builder
      err = resp.Write(&data)
      if err != nil {
         return nil, err
      }
      return nil, errors.New(data.String())
   }
   for _, cookie := range resp.Cookies() {
      if cookie.Name == "idsession" {
         return cookie, nil
      }
   }
   return nil, http.ErrNoCookie
}

var Territory = "US"

func sign(method, path string, head http.Header, body []byte) string {
   timestamp := time.Now().Unix()
   text_headers := func() string {
      var s []string
      for k := range head {
         k = strings.ToLower(k)
         if strings.HasPrefix(k, "x-skyott-") {
            s = append(s, k+": "+head.Get(k)+"\n")
         }
      }
      slices.Sort(s)
      return strings.Join(s, "")
   }()
   headers_md5 := md5.Sum([]byte(text_headers))
   payload_md5 := md5.Sum(body)
   signature := func() string {
      h := hmac.New(sha1.New, []byte(sky_key))
      fmt.Fprintln(h, method)
      fmt.Fprintln(h, path)
      fmt.Fprintln(h)
      fmt.Fprintln(h, sky_client)
      fmt.Fprintln(h, sky_version)
      fmt.Fprintf(h, "%x\n", headers_md5)
      fmt.Fprintln(h, timestamp)
      fmt.Fprintf(h, "%x\n", payload_md5)
      hashed := h.Sum(nil)
      return base64.StdEncoding.EncodeToString(hashed[:])
   }
   sky_ott := func() string {
      data := []byte("SkyOTT")
      // must be quoted
      data = fmt.Appendf(data, " client=%q", sky_client)
      // must be quoted
      data = fmt.Appendf(data, ",signature=%q", signature())
      // must be quoted
      data = fmt.Appendf(data, `,timestamp="%v"`, timestamp)
      // must be quoted
      data = fmt.Appendf(data, ",version=%q", sky_version)
      return string(data)
   }
   return sky_ott()
}

func (a *AssetEndpoint) Dash() (*Dash, error) {
   resp, err := http.Get(a.Url)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Dash
   result.Body, err = io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   result.Url = resp.Request.URL
   return &result, nil
}

type Dash struct {
   Body []byte
   Url  *url.URL
}

type AssetEndpoint struct {
   Cdn string
   Url string
}

type Playout struct {
   Asset struct {
      Endpoints []AssetEndpoint
   }
   Description string
   Protection  struct {
      LicenceAcquisitionUrl string
   }
}
func (p *Playout) Widevine(body []byte) ([]byte, error) {
   req, err := http.NewRequest(
      "POST", p.Protection.LicenceAcquisitionUrl, bytes.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set(
      "x-sky-signature", sign(req.Method, req.URL.Path, req.Header, body),
   )
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   return io.ReadAll(resp.Body)
}

func (p *Playout) Fastly() (string, bool) {
   for _, endpoint := range p.Asset.Endpoints {
      if endpoint.Cdn == "FASTLY" {
         return endpoint.Url, true
      }
   }
   return "", false
}

func (a *AuthToken) Playout(contentId string) (*Playout, error) {
   body, err := json.Marshal(map[string]any{
      "contentId": contentId,
      "device": map[string]any{
         "capabilities": []any{
            map[string]string{
               "acodec":     "AAC",
               "container":  "ISOBMFF",
               "protection": "WIDEVINE",
               "transport":  "DASH",
               "vcodec":     "H264",
            },
         },
         "maxVideoFormat": "HD",
      },
      "personaParentalControlRating": 9,
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://ovp.peacocktv.com/video/playouts/vod",
      bytes.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   // `application/json` fails
   req.Header.Set("content-type", "application/vnd.playvod.v1+json")
   req.Header.Set("x-skyott-usertoken", a.UserToken)
   req.Header.Set(
      "x-sky-signature", sign(req.Method, req.URL.Path, req.Header, body),
   )
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Playout
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.Description != "" {
      return nil, errors.New(result.Description)
   }
   return &result, nil
}

const (
   sky_client  = "NBCU-ANDROID-v3"
   sky_key     = "JuLQgyFz9n89D9pxcN6ZWZXKWfgj2PNBUb32zybj"
   sky_version = "1.0"
)

func (a *AuthToken) Fetch(idSession *http.Cookie) error {
   body, err := json.Marshal(map[string]any{
      "auth": map[string]string{
         "authScheme":        "MESSO",
         "proposition":       "NBCUOTT",
         "provider":          "NBCU",
         "providerTerritory": Territory,
      },
      "device": map[string]string{
         // if empty /drm/widevine/acquirelicense will fail with
         // {
         //    "errorCode": "OVP_00306",
         //    "description": "Security failure"
         // }
         "drmDeviceId": "UNKNOWN",
         // if incorrect /video/playouts/vod will fail with
         // {
         //    "errorCode": "OVP_00311",
         //    "description": "Unknown deviceId"
         // }
         // changing this too often will result in a four hour block
         // {
         //    "errorCode": "OVP_00014",
         //    "description": "Maximum number of streaming devices exceeded"
         // }
         "id":       "PC",
         "platform": "ANDROIDTV",
         "type":     "TV",
      },
   })
   if err != nil {
      return err
   }
   req, err := http.NewRequest(
      "POST", "https://ovp.peacocktv.com/auth/tokens", bytes.NewReader(body),
   )
   if err != nil {
      return err
   }
   req.AddCookie(idSession)
   req.Header.Set("content-type", "application/vnd.tokens.v1+json")
   req.Header.Set("x-sky-signature", sign(req.Method, req.URL.Path, nil, body))
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   err = json.NewDecoder(resp.Body).Decode(a)
   if err != nil {
      return err
   }
   if a.Description != "" {
      return errors.New(a.Description)
   }
   return nil
}

// userToken is good for one day
type AuthToken struct {
   Description string
   UserToken   string
}
