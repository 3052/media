package movistar

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

func (t *Token) Unmarshal(data Byte[Token]) error {
   return json.Unmarshal(data, t)
}

// 10 days
type Token struct {
   AccessToken string `json:"access_token"`
   ExpiresIn   int64  `json:"expires_in"`
}

type Byte[T any] []byte

const device_type = "SMARTTV_OTT"

type Details struct {
   Id       int // contentID
   VodItems []struct {
      CasId    string // drmMediaID
      UrlVideo string // MPD mullvad
   }
}

func (d *Details) New(id int64) error {
   req, _ := http.NewRequest("", "https://ottcache.dof6.com", nil)
   req.URL.Path = func() string {
      b := []byte("/movistarplus/amazon.tv/contents/")
      b = strconv.AppendInt(b, id, 10)
      b = append(b, "/details"...)
      return string(b)
   }()
   req.URL.RawQuery = "mdrm=true"
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(d)
}

// mullvad pass
func NewToken(username, password string) (Byte[Token], error) {
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

///

type init_data struct {
   AccountNumber string
   Token         string
}

// mullvad pass
func (o oferta) init_data(device1 device) (*init_data, error) {
   data, err := json.Marshal(map[string]string{
      "accountNumber": o.AccountNumber,
      "deviceType":    device_type, // NEEDED FOR /Session
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://clientservices.dof6.com?qspVersion=ssp",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("content-type", "application/json")
   req.URL.Path = func() string {
      var b strings.Builder
      b.WriteString("/movistarplus/amazon.tv/sdp/mediaPlayers/")
      b.WriteString(string(device1))
      b.WriteString("/initData")
      return b.String()
   }()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   init1 := &init_data{}
   err = json.NewDecoder(resp.Body).Decode(init1)
   if err != nil {
      return nil, err
   }
   return init1, nil
}

func (d *device) unmarshal(data Byte[device]) error {
   return json.Unmarshal(data, d)
}

type device string

// mullvad pass
func (t *Token) device(oferta1 *oferta) (Byte[device], error) {
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
func (t *Token) oferta() (*oferta, error) {
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
