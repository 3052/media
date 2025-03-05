package kanopy

import (
   "bytes"
   "encoding/json"
   "io"
   "net/http"
   "strconv"
)

func (n *Login) Widevine(m *Manifest, data []byte) ([]byte, error) {
   req, err := http.NewRequest(
      "POST", "https://www.kanopy.com", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.URL.Path = "/kapi/licenses/widevine/" + m.DrmLicenseId
   req.Header = http.Header{
      "user-agent":    {user_agent},
      "x-version":     {x_version},
      "authorization": {"Bearer " + n.Jwt},
   }
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (m *Manifest) Mpd() (*http.Response, error) {
   req, err := http.NewRequest("", m.Url, nil)
   if err != nil {
      return nil, err
   }
   req.Header.Set("user-agent", "Mozilla")
   return http.DefaultClient.Do(req)
}

const (
   user_agent = "!"
   x_version  = "!/!/!/!"
)

type Membership struct {
   DomainId int
}

func (n *Login) Membership() (*Membership, error) {
   req, _ := http.NewRequest("", "https://www.kanopy.com", nil)
   req.URL.Path = "/kapi/memberships"
   req.URL.RawQuery = "userId=" + strconv.Itoa(n.UserId)
   req.Header = http.Header{
      "authorization": {"Bearer " + n.Jwt},
      "user-agent":    {user_agent},
      "x-version":     {x_version},
   }
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var value struct {
      List []Membership
   }
   err = json.NewDecoder(resp.Body).Decode(&value)
   if err != nil {
      return nil, err
   }
   return &value.List[0], nil
}

func NewLogin(email, password string) (Byte[Login], error) {
   data, err := json.Marshal(map[string]any{
      "credentialType": "email",
      "emailUser": map[string]string{
         "email":    email,
         "password": password,
      },
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://www.kanopy.com/kapi/login", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header = http.Header{
      "content-type": {"application/json"},
      "user-agent":   {user_agent},
   }
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

type Byte[T any] []byte

func (n *Login) Unmarshal(data Byte[Login]) error {
   return json.Unmarshal(data, n)
}

type Manifest struct {
   DrmLicenseId string
   ManifestType string
   Url          string
}

// good for 10 years
type Login struct {
   Jwt    string
   UserId int
}

func (p *Plays) Dash() (*Manifest, bool) {
   for _, value := range p.Manifests {
      if value.ManifestType == "dash" {
         return &value, true
      }
   }
   return nil, false
}

type Plays struct {
   ErrorMsgLong string `json:"error_msg_long"`
   Manifests []Manifest
}

func (n *Login) Plays(member *Membership, video_id int) (Byte[Plays], error) {
   data, err := json.Marshal(map[string]int{
      "domainId": member.DomainId,
      "userId":   n.UserId,
      "videoId":  video_id,
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
   req.Header = http.Header{
      "authorization": {"Bearer " + n.Jwt},
      "content-type":  {"application/json"},
      "user-agent":    {user_agent},
      "x-version":     {x_version},
   }
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (p *Plays) Unmarshal(data Byte[Plays]) error {
   return json.Unmarshal(data, p)
}
