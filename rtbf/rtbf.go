package rtbf

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "log"
   "net/http"
   "net/url"
   "strings"
)

func GetPath(rawUrl string) (string, error) {
   u, err := url.Parse(rawUrl)
   if err != nil {
      return "", err
   }
   if u.Scheme == "" {
      return "", errors.New("invalid URL: scheme is missing")
   }
   return u.Path, nil
}

func FetchAssetId(path string) (string, error) {
   resp, err := http.Get(
      "https://bff-service.rtbf.be/auvio/v1.23/pages" + path,
   )
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return "", errors.New(resp.Status)
   }
   var page struct {
      Data struct {
         Content struct {
            AssetId string
            Media   *struct {
               AssetId string
            }
         }
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&page)
   if err != nil {
      return "", err
   }
   content := page.Data.Content
   if content.AssetId != "" {
      return content.AssetId, nil
   }
   if content.Media != nil {
      return content.Media.AssetId, nil
   }
   return "", errors.New("assetId not found")
}

var Transport = http.Transport{
   Proxy: func(req *http.Request) (*url.URL, error) {
      log.Println(req.Method, req.URL)
      return http.ProxyFromEnvironment(req)
   },
}

func (e *Entitlement) Widevine(data []byte) ([]byte, error) {
   req, err := http.NewRequest(
      "POST", "https://exposure.api.redbee.live", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.URL.Path = "/v2/license/customer/RTBF/businessunit/Auvio/widevine"
   req.URL.RawQuery = url.Values{
      "contentId":  {e.AssetId},
      "ls_session": {e.PlayToken},
   }.Encode()
   req.Header.Set("content-type", "application/x-protobuf")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      var value struct {
         Message string
      }
      err = json.NewDecoder(resp.Body).Decode(&value)
      if err != nil {
         return nil, err
      }
      return nil, errors.New(value.Message)
   }
   return io.ReadAll(resp.Body)
}

type Entitlement struct {
   AssetId   string
   Formats   []Format
   Message   string
   PlayToken string
}

func (e *Entitlement) Unmarshal(data Byte[Entitlement]) error {
   err := json.Unmarshal(data, e)
   if err != nil {
      return err
   }
   if e.Message != "" {
      return errors.New(e.Message)
   }
   return nil
}

type Format struct {
   Format       string
   MediaLocator string // MPD
}

type GigyaLogin struct {
   SessionToken string
}

type Jwt struct {
   ErrorMessage string
   IdToken      string `json:"id_token"`
}

type Login struct {
   ErrorMessage string
   SessionInfo  struct {
      CookieValue string
   }
}

// hard coded in JavaScript
const api_key = "4_Ml_fJ47GnBAW6FrPzMxh0w"

func (j *Jwt) Login() (*GigyaLogin, error) {
   data, err := json.Marshal(map[string]any{
      "device": map[string]string{
         "deviceId": "",
         "type":     "WEB",
      },
      "jwt": j.IdToken,
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://exposure.api.redbee.live", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.URL.Path = "/v2/customer/RTBF/businessunit/Auvio/auth/gigyaLogin"
   req.Header.Set("content-type", "application/json")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   gigya := &GigyaLogin{}
   err = json.NewDecoder(resp.Body).Decode(gigya)
   if err != nil {
      return nil, err
   }
   return gigya, nil
}

func (l *Login) Jwt() (*Jwt, error) {
   resp, err := http.PostForm(
      "https://login.auvio.rtbf.be/accounts.getJWT", url.Values{
         "APIKey":      {api_key},
         "login_token": {l.SessionInfo.CookieValue},
      },
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var token Jwt
   err = json.NewDecoder(resp.Body).Decode(&token)
   if err != nil {
      return nil, err
   }
   if token.ErrorMessage != "" {
      return nil, errors.New(token.ErrorMessage)
   }
   return &token, nil
}

func NewLogin(id, password string) (Byte[Login], error) {
   resp, err := http.PostForm(
      "https://login.auvio.rtbf.be/accounts.login", url.Values{
         "APIKey":   {api_key},
         "loginID":  {id},
         "password": {password},
      },
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (l *Login) Unmarshal(data Byte[Login]) error {
   err := json.Unmarshal(data, l)
   if err != nil {
      return err
   }
   if l.ErrorMessage != "" {
      return errors.New(l.ErrorMessage)
   }
   return nil
}

func (e *Entitlement) Dash() (*Format, bool) {
   for _, format_var := range e.Formats {
      if format_var.Format == "DASH" {
         return &format_var, true
      }
   }
   return nil, false
}

type Byte[T any] []byte

func (g *GigyaLogin) Entitlement(assetId string) (Byte[Entitlement], error) {
   req, _ := http.NewRequest("", "https://exposure.api.redbee.live", nil)
   req.URL.Path = func() string {
      var data strings.Builder
      data.WriteString("/v2/customer/RTBF/businessunit/Auvio/entitlement/")
      data.WriteString(assetId)
      data.WriteString("/play")
      return data.String()
   }()
   req.Header.Set("x-forwarded-for", "91.90.123.17")
   req.Header.Set("authorization", "Bearer "+g.SessionToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}
