package criterion

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
)

func (f *File) Widevine(data []byte) ([]byte, error) {
   req, err := http.NewRequest(
      "POST", "https://drm.vhx.com/v2/widevine", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.URL.RawQuery = "token=" + f.DrmAuthorizationToken
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (f *Files) Unmarshal(data Byte[Files]) error {
   return json.Unmarshal(data, f)
}

type Byte[T any] []byte

func (t Token) Files(video1 *Video) (Byte[Files], error) {
   req, err := http.NewRequest("", video1.Links.Files.Href, nil)
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+t.AccessToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

type File struct {
   DrmAuthorizationToken string `json:"drm_authorization_token"`
   Links                 struct {
      Source struct {
         Href string // MPD
      }
   } `json:"_links"`
   Method string
}

func (f Files) Dash() (*File, bool) {
   for _, file1 := range f {
      if file1.Method == "dash" {
         return &file1, true
      }
   }
   return nil, false
}

const client_id = "9a87f110f79cd25250f6c7f3a6ec8b9851063ca156dae493bf362a7faf146c78"

type Files []File

type Video struct {
   Links struct {
      Files struct {
         Href string
      }
   } `json:"_links"`
   Message string
   Name    string
}

type Token struct {
   AccessToken string `json:"access_token"`
}

func (t Token) Video(slug string) (*Video, error) {
   req, _ := http.NewRequest("", "https://api.vhx.com", nil)
   req.URL.Path = "/videos/" + slug
   req.URL.RawQuery = "url=" + url.QueryEscape(slug)
   req.Header.Set("authorization", "Bearer "+t.AccessToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var video1 Video
   err = json.NewDecoder(resp.Body).Decode(&video1)
   if err != nil {
      return nil, err
   }
   if video1.Message != "" {
      return nil, errors.New(video1.Message)
   }
   return &video1, nil
}

func NewToken(username, password string) (Byte[Token], error) {
   resp, err := http.PostForm("https://auth.vhx.com/v1/oauth/token", url.Values{
      "client_id":  {client_id},
      "grant_type": {"password"},
      "password":   {password},
      "username":   {username},
   })
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (t *Token) Unmarshal(data Byte[Token]) error {
   return json.Unmarshal(data, t)
}
