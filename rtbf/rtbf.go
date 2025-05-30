package rtbf

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "strings"
)

func (e *Entitlement) Widevine(data []byte) ([]byte, error) {
   req, err := http.NewRequest(
      "POST", "https://rbm-rtbf.live.ott.irdeto.com", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.URL.Path = "/licenseServer/widevine/v1/rbm-rtbf/license"
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
   return io.ReadAll(resp.Body)
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

type Content struct {
   AssetId string
   Media   *struct {
      AssetId string
   }
}

func (c *Content) get_asset_id() string {
   if c.AssetId != "" {
      return c.AssetId
   }
   return c.Media.AssetId
}

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

func (n *Login) Jwt() (*Jwt, error) {
   resp, err := http.PostForm(
      "https://login.auvio.rtbf.be/accounts.getJWT", url.Values{
         "APIKey":      {api_key},
         "login_token": {n.SessionInfo.CookieValue},
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

func (n *Login) Unmarshal(data Byte[Login]) error {
   err := json.Unmarshal(data, n)
   if err != nil {
      return err
   }
   if n.ErrorMessage != "" {
      return errors.New(n.ErrorMessage)
   }
   return nil
}

func (e *Entitlement) Dash() (*Format, bool) {
   for _, format1 := range e.Formats {
      if format1.Format == "DASH" {
         return &format1, true
      }
   }
   return nil, false
}

type Entitlement struct {
   AssetId   string
   PlayToken string
   Formats   []Format
}

type Byte[T any] []byte

func (e *Entitlement) Unmarshal(data Byte[Entitlement]) error {
   return json.Unmarshal(data, e)
}

func (g *GigyaLogin) Entitlement(c *Content) (Byte[Entitlement], error) {
   req, _ := http.NewRequest("", "https://exposure.api.redbee.live", nil)
   req.URL.Path = func() string {
      var data strings.Builder
      data.WriteString("/v2/customer/RTBF/businessunit/Auvio/entitlement/")
      data.WriteString(c.get_asset_id())
      data.WriteString("/play")
      return data.String()
   }()
   req.Header.Set("x-forwarded-for", "91.90.123.17")
   req.Header.Set("authorization", "Bearer " + g.SessionToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

type Address [1]string

func (a *Address) New(data string) {
   data = strings.TrimPrefix(data, "https://")
   a[0] = strings.TrimPrefix(data, "auvio.rtbf.be")
}

func (a Address) Content() (*Content, error) {
   resp, err := http.Get(
      "https://bff-service.rtbf.be/auvio/v1.23/pages" + a[0],
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var value struct {
      Data struct {
         Content Content
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   return &value.Data.Content, nil
}
