package kanopy

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "path"
   "strconv"
   "strings"
)

type Video struct {
   Alias   string
   VideoId int
}

// https://kanopy.com/video/6440418
// https://kanopy.com/video/genius-party
// https://kanopy.com/en/video/genius-party
// https://kanopy.com/en/product/genius-party
func (v *Video) Parse(inputUrl string) error {
   parsedUrl, err := url.Parse(inputUrl)
   if err != nil {
      return err
   }
   if !strings.Contains(parsedUrl.Host, "kanopy.com") {
      return errors.New("invalid domain")
   }
   // Get the directory of the path (removes the final identifier).
   // e.g., "/en/product/genius-party" -> "/en/product"
   dir := path.Dir(parsedUrl.Path)
   // Update: Check if the directory ends with "/video" OR "/product".
   // This supports:
   // - /video/{id}
   // - /en/video/{id}
   // - /en/product/{id}
   if !strings.HasSuffix(dir, "/video") && !strings.HasSuffix(dir, "/product") {
      return errors.New("invalid path structure")
   }
   identifier := path.Base(parsedUrl.Path)
   numericId, err := strconv.Atoi(identifier)
   if err != nil {
      v.Alias = identifier
   } else {
      v.VideoId = numericId
   }
   return nil
}

const x_version = "!/!/!/!"

func (l *Login) Video(alias string) (*Video, error) {
   var req http.Request
   req.Header = http.Header{}
   req.URL = &url.URL{
      Scheme: "https",
      Host:   "www.kanopy.com",
      Path:   "/kapi/videos/alias/" + alias,
   }
   req.Header.Set("x-version", x_version)
   req.Header.Set("authorization", "Bearer "+l.Jwt)
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Video Video
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result.Video, nil
}

// good for 10 years
type Login struct {
   Jwt    string
   UserId int
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
   req.Header.Set("authorization", "Bearer "+l.Jwt)
   req.Header.Set("content-type", "application/json")
   req.Header.Set("user-agent", user_agent)
   req.Header.Set("x-version", x_version)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Plays
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.ErrorMsgLong != "" {
      return nil, errors.New(result.ErrorMsgLong)
   }
   return &result, nil
}

const user_agent = "!"

type Membership struct {
   DomainId int
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

func (l *Login) Widevine(manifest *PlayManifest, data []byte) ([]byte, error) {
   req, err := http.NewRequest(
      "POST", "https://www.kanopy.com", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.URL.Path = "/kapi/licenses/widevine/" + manifest.DrmLicenseId
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

func (p *PlayManifest) Dash() (*Dash, error) {
   req, err := http.NewRequest("", p.Url, nil)
   if err != nil {
      return nil, err
   }
   req.Header.Set("user-agent", "Mozilla")
   resp, err := http.DefaultClient.Do(req)
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

type Dash struct {
   Body []byte
   Url  *url.URL
}

func (l *Login) Membership() (*Membership, error) {
   var req http.Request
   req.Header = http.Header{}
   req.Header.Set("authorization", "Bearer "+l.Jwt)
   req.Header.Set("user-agent", user_agent)
   req.Header.Set("x-version", x_version)
   req.URL = &url.URL{
      Scheme:   "https",
      Host:     "www.kanopy.com",
      Path:     "/kapi/memberships",
      RawQuery: "userId=" + strconv.Itoa(l.UserId),
   }
   resp, err := http.DefaultClient.Do(&req)
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

type Plays struct {
   Captions []struct {
      Files []struct {
         Url string
      }
   }
   ErrorMsgLong string `json:"error_msg_long"`
   Manifests    []PlayManifest
}

type PlayManifest struct {
   DrmLicenseId string
   ManifestType string
   Url          string
}

func (p *Plays) Dash() (*PlayManifest, error) {
   for _, manifest := range p.Manifests {
      if manifest.ManifestType == "dash" {
         return &manifest, nil
      }
   }
   return nil, errors.New("dash manifest not found")
}
