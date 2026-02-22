package peacock

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "strings"
)

// userToken is good for one day
type Token struct {
   Description string
   UserToken   string
}

func (t *Token) Playout(variantId string) (*Playout, error) {
   body, err := json.Marshal(map[string]any{
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
      // "contentId": "GMO_00000000261361_02_HDSDR",
      "providerVariantId": variantId,
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
   req.Header.Set("x-skyott-usertoken", t.UserToken)
   req.Header.Set(
      "x-sky-signature",
      generate_sky_ott(req.Method, req.URL.Path, req.Header, body),
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

var Territory = "US"

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

type Playout struct {
   Asset struct {
      Endpoints []AssetEndpoint
   }
   Description string
   Protection  struct {
      LicenceAcquisitionUrl string
   }
}

// 1080p L3
func (p *Playout) Widevine(body []byte) ([]byte, error) {
   req, err := http.NewRequest(
      "POST", p.Protection.LicenceAcquisitionUrl, bytes.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set(
      "x-sky-signature",
      generate_sky_ott(req.Method, req.URL.Path, req.Header, body),
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

func (p *Playout) Fastly() (*AssetEndpoint, error) {
   for _, endpoint := range p.Asset.Endpoints {
      if endpoint.Cdn == "FASTLY" {
         return &endpoint, nil
      }
   }
   return nil, errors.New("FASTLY endpoint not found")
}

func (t *Token) Fetch(idSession *http.Cookie) error {
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
   req.Header.Set(
      "x-sky-signature", generate_sky_ott(req.Method, req.URL.Path, nil, body),
   )
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   err = json.NewDecoder(resp.Body).Decode(t)
   if err != nil {
      return err
   }
   if t.Description != "" {
      return errors.New(t.Description)
   }
   return nil
}

type Dash struct {
   Body []byte
   Url  *url.URL
}

type AssetEndpoint struct {
   Cdn string
   Url string
}
