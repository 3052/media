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

func (l *Login) PlayResponse(member *Membership, videoId int) (*PlayResponse, error) {
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
   req.Header.Set("authorization", "Bearer "+l.Jwt)
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
   play := &PlayResponse{}
   err = json.NewDecoder(resp.Body).Decode(play)
   if err != nil {
      return nil, err
   }
   return play, nil
}

func (s *StreamInfo) Mpd() (*url.URL, []byte, error) {
   req, err := http.NewRequest("", s.Url, nil)
   if err != nil {
      return nil, nil, err
   }
   req.Header.Set("user-agent", "Mozilla")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, nil, err
   }
   defer resp.Body.Close()
   data, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, nil, err
   }
   return resp.Request.URL, data, nil
}

const (
   user_agent = "!"
   x_version  = "!/!/!/!"
)

type StreamInfo struct {
   DrmLicenseId string
   ManifestType string
   Url          string
}

type Membership struct {
   DomainId int
}

func (p *PlayResponse) Dash() (*StreamInfo, bool) {
   for _, info := range p.Manifests {
      if info.ManifestType == "dash" {
         return &info, true
      }
   }
   return nil, false
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

type PlayResponse struct {
   ErrorMsgLong string `json:"error_msg_long"`
   Manifests    []StreamInfo
}

func (l *Login) Membership() (*Membership, error) {
   req, _ := http.NewRequest("", "https://www.kanopy.com", nil)
   req.URL.Path = "/kapi/memberships"
   req.URL.RawQuery = "userId=" + strconv.Itoa(l.UserId)
   req.Header.Set("authorization", "Bearer "+l.Jwt)
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
   var result struct {
      List []Membership
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result.List[0], nil
}

func (l *Login) Widevine(info *StreamInfo, data []byte) ([]byte, error) {
   req, err := http.NewRequest(
      "POST", "https://www.kanopy.com", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.URL.Path = "/kapi/licenses/widevine/" + info.DrmLicenseId
   req.Header.Set("user-agent", user_agent)
   req.Header.Set("x-version", x_version)
   req.Header.Set("authorization", "Bearer "+l.Jwt)
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
