package movistar

import (
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "strings"
)

func (d *device) unmarshal(data Byte[device]) error {
   return json.Unmarshal(data, d)
}

type device string

// mullvad pass
func (t *token) device(oferta1 *oferta) (Byte[device], error) {
   req, err := http.NewRequest(
      "POST", "https://auth.dof6.com?qspVersion=ssp", nil,
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+t.AccessToken)
   req.URL.Path = func() string {
      var b strings.Builder
      b.WriteString("/movistarplus/amazon.tv/accounts/")
      b.WriteString(oferta1.AccountNumber)
      b.WriteString("/devices/")
      return b.String()
   }()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusCreated {
      return nil, errors.New(resp.Status)
   }
   return io.ReadAll(resp.Body)
}
type oferta struct {
   AccountNumber string
}

// mullvad pass
func (t *token) oferta() (*oferta, error) {
   req, _ := http.NewRequest("", "https://auth.dof6.com", nil)
   req.URL.Path = "/movistarplus/api/devices/amazon.tv/users/authenticate"
   req.Header.Set("authorization", "Bearer "+t.AccessToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var value struct {
      Ofertas []oferta
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   return &value.Ofertas[0], nil
}

// mullvad pass
func new_token(username, password string) (Byte[token], error) {
   resp, err := http.PostForm(
      "https://auth.dof6.com/auth/oauth2/token?deviceClass=amazon.tv",
      url.Values{
         "grant_type": {"password"},
         "password":   {password},
         "username":   {username},
      },
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   return io.ReadAll(resp.Body)
}

// 10 days
type token struct {
   AccessToken string `json:"access_token"`
   ExpiresIn   int64  `json:"expires_in"`
}

type Byte[T any] []byte

func (t *token) unmarshal(data Byte[token]) error {
   return json.Unmarshal(data, t)
}
