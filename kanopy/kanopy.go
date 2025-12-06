package kanopy

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "strconv"
)

type Cache struct {
   Mpd      *url.URL
   MpdBody  []byte
   Login *Login
   Manifest *Manifest
}

func (m *Manifest) Mpd(storage *Cache) error {
   req, err := http.NewRequest("", m.Url, nil)
   if err != nil {
      return err
   }
   req.Header.Set("user-agent", "Mozilla")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   storage.MpdBody, err = io.ReadAll(resp.Body)
   if err != nil {
      return err
   }
   storage.Mpd = resp.Request.URL
   return nil
}

// good for 10 years
type Login struct {
   Jwt    string
   UserId int
}

func (l *Login) Fetch(email, password string) error {
   data, err := json.Marshal(map[string]any{
      "credentialType": "email",
      "emailUser": map[string]string{
         "email":    email,
         "password": password,
      },
   })
   if err != nil {
      return err
   }
   req, err := http.NewRequest(
      "POST", "https://www.kanopy.com/kapi/login", bytes.NewReader(data),
   )
   if err != nil {
      return err
   }
   req.Header.Set("content-type", "application/json")
   req.Header.Set("user-agent", user_agent)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return errors.New(resp.Status)
   }
   return json.NewDecoder(resp.Body).Decode(l)
}

type Plays struct {
   ErrorMsgLong string `json:"error_msg_long"`
   Manifests []Manifest
}

func (l *Login) Membership() (*Membership, error) {
   req, _ := http.NewRequest("", "https://www.kanopy.com", nil)
   req.URL.Path = "/kapi/memberships"
   req.URL.RawQuery = "userId=" + strconv.Itoa(l.UserId)
   req.Header.Set("authorization", "Bearer " + l.Jwt)
   req.Header.Set("user-agent", user_agent)
   req.Header.Set("x-version", x_version)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var value struct {
      List []Membership
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   return &value.List[0], nil
}

func (l *Login) Plays(member *Membership, videoId int) (*Plays, error) {
   data, err := json.Marshal(map[string]int{
      "domainId": member.DomainId,
      "userId":   l.UserId,
      "videoId":  videoId,
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://www.kanopy.com/kapi/plays", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer " + l.Jwt)
   req.Header.Set("content-type", "application/json")
   req.Header.Set("user-agent", user_agent)
   req.Header.Set("x-version", x_version)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   plays_var := &Plays{}
   err = json.NewDecoder(resp.Body).Decode(plays_var)
   if err != nil {
      return nil, err
   }
   return plays_var, nil
}

func (l *Login) Widevine(manifestVar *Manifest, data []byte) ([]byte, error) {
   req, err := http.NewRequest(
      "POST", "https://www.kanopy.com", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.URL.Path = "/kapi/licenses/widevine/" + manifestVar.DrmLicenseId
   req.Header.Set("user-agent", user_agent)
   req.Header.Set("x-version", x_version)
   req.Header.Set("authorization", "Bearer " + l.Jwt)
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

const (
   user_agent = "!"
   x_version  = "!/!/!/!"
)

type Manifest struct {
   DrmLicenseId string
   ManifestType string
   Url          string
}

type Membership struct {
   DomainId int
}

func (p *Plays) Dash() (*Manifest, bool) {
   for _, value := range p.Manifests {
      if value.ManifestType == "dash" {
         return &value, true
      }
   }
   return nil, false
}
