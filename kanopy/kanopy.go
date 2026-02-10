package kanopy

import (
   "bytes"
   "encoding/json"
   "io"
   "net/http"
   "strconv"
)

func (n *Login) Plays(member *Membership, video_id int) (*Plays, error) {
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
   value := &Plays{}
   err = json.NewDecoder(resp.Body).Decode(value)
   if err != nil {
      return nil, err
   }
   return value, nil
}

func (p *Plays) Dash() (*Manifest, bool) {
   for _, value := range p.Manifests {
      if value.ManifestType == "dash" {
         return &value, true
      }
   }
   return nil, false
}

const (
   user_agent = "!"
   x_version  = "!/!/!/!"
)

func (n *Login) Unmarshal(data []byte) error {
   return json.Unmarshal(data, n)
}

type Membership struct {
   DomainId int
}

type Manifest struct {
   DrmLicenseId string
   ManifestType string
   Url          string
}

type Plays struct {
   ErrorMsgLong string `json:"error_msg_long"`
   Manifests []Manifest
}

// good for 10 years
type Login struct {
   Jwt    string
   UserId int
}

type Client struct {
   Manifest *Manifest
   Login    *Login
}

func (m *Manifest) Mpd() (*http.Response, error) {
   req, err := http.NewRequest("", m.Url, nil)
   if err != nil {
      return nil, err
   }
   req.Header.Set("user-agent", "Mozilla")
   return http.DefaultClient.Do(req)
}

func (c *Client) License(data []byte) ([]byte, error) {
   req, err := http.NewRequest(
      "POST", "https://www.kanopy.com", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.URL.Path = "/kapi/licenses/widevine/" + c.Manifest.DrmLicenseId
   req.Header = http.Header{
      "authorization": {"Bearer " + c.Login.Jwt},
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
func (Login) Marshal(email, password string) ([]byte, error) {
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
