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

type Session struct {
   ResultData struct {
      Ctoken string // ONE TIME USE
   }
   ResultText string
}

func (d Device) Session(init1 *InitData, details1 *Details) (*Session, error) {
   data, err := json.Marshal(map[string]any{
      "contentID":  details1.Id,
      "drmMediaID": details1.VodItems[0].CasId,
      "streamType": "AST",
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://alkasvaspub.imagenio.telefonica.net",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.URL.Path = func() string {
      var b strings.Builder
      b.WriteString("/asvas/ccs/")
      b.WriteString(init1.AccountNumber)
      b.WriteByte('/')
      b.WriteString(device_type)
      b.WriteByte('/')
      b.WriteString(string(d))
      b.WriteString("/Session")
      return b.String()
   }()
   req.Header = http.Header{
      "content-type": {"application/json"},
      "x-hzid":       {init1.Token},
   }
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var value Session
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   if resp.StatusCode != http.StatusCreated {
      return nil, errors.New(value.ResultText)
   }
   return &value, nil
}

func (s Session) Widevine(data []byte) ([]byte, error) {
   req, err := http.NewRequest(
      "POST", "https://wv-ottlic-f3.imagenio.telefonica.net",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.URL.Path = "/TFAESP/wvls/contentlicenseservice/v1/licenses"
   req.Header.Set("nv-authorizations", s.ResultData.Ctoken)
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
func NewDetails(id int64) (Byte[Details], error) {
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
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   return io.ReadAll(resp.Body)
}

const device_type = "SMARTTV_OTT"

type Byte[T any] []byte

type Details struct {
   Id       int // contentID
   VodItems []struct {
      CasId    string // drmMediaID
      UrlVideo string
   }
}

func (d *Details) Unmarshal(data Byte[Details]) error {
   return json.Unmarshal(data, d)
}

// EVEN THE CONTENT IS GEO BLOCKED
func (d *Details) Mpd() (*http.Response, error) {
   req, err := http.NewRequest("", d.VodItems[0].UrlVideo, nil)
   if err != nil {
      return nil, err
   }
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   return resp, nil
}

type Device string

func (d *Device) Unmarshal(data Byte[Device]) error {
   return json.Unmarshal(data, d)
}

type InitData struct {
   AccountNumber string
   Token         string
}

type Oferta struct {
   AccountNumber string
}

func (o Oferta) InitData(device1 Device) (*InitData, error) {
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
   init1 := &InitData{}
   err = json.NewDecoder(resp.Body).Decode(init1)
   if err != nil {
      return nil, err
   }
   return init1, nil
}

func (t *Token) Unmarshal(data Byte[Token]) error {
   return json.Unmarshal(data, t)
}

// 10 days
type Token struct {
   AccessToken string `json:"access_token"`
   ExpiresIn   int64  `json:"expires_in"`
}

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

func (t *Token) Device(oferta1 *Oferta) (Byte[Device], error) {
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

func (t *Token) Oferta() (*Oferta, error) {
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
      Ofertas []Oferta
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   return &value.Ofertas[0], nil
}
