package criterion

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "log"
   "net/http"
   "net/url"
)

func FetchToken(username, password string) (TokenData, error) {
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

type Token struct {
   AccessToken  string `json:"access_token"`
   RefreshToken string `json:"refresh_token"`
}

func (t *Token) Refresh() (TokenData, error) {
   resp, err := http.PostForm("https://auth.vhx.com/v1/oauth/token", url.Values{
      "client_id":     {client_id},
      "grant_type":    {"refresh_token"},
      "refresh_token": {t.RefreshToken},
   })
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

type TokenData []byte

func (t *Token) Unmarshal(data TokenData) error {
   return json.Unmarshal(data, t)
}

var Transport = http.Transport{
   Proxy: func(req *http.Request) (*url.URL, error) {
      log.Println(req.Method, req.URL)
      return http.ProxyFromEnvironment(req)
   },
}

const client_id = "9a87f110f79cd25250f6c7f3a6ec8b9851063ca156dae493bf362a7faf146c78"

func (t *Token) Video(slug string) (*Video, error) {
   req, _ := http.NewRequest("", "https://api.vhx.com", nil)
   req.URL.Path = "/videos/" + slug
   req.URL.RawQuery = "url=" + url.QueryEscape(slug)
   req.Header.Set("authorization", "Bearer "+t.AccessToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var videoVar Video
   err = json.NewDecoder(resp.Body).Decode(&videoVar)
   if err != nil {
      return nil, err
   }
   if videoVar.Message != "" {
      return nil, errors.New(videoVar.Message)
   }
   return &videoVar, nil
}

type Video struct {
   Links struct {
      Files struct {
         Href string
      }
   } `json:"_links"`
   Message string
   Name    string
}

func (t *Token) Files(videoVar *Video) (Files, error) {
   req, err := http.NewRequest("", videoVar.Links.Files.Href, nil)
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+t.AccessToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var filesVar Files
   err = json.NewDecoder(resp.Body).Decode(&filesVar)
   if err != nil {
      return nil, err
   }
   return filesVar, nil
}

type Files []File

type File struct {
   DrmAuthorizationToken string `json:"drm_authorization_token"`
   Links                 struct {
      Source struct {
         Href string // MPD
      }
   } `json:"_links"`
   Method string
}

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

func (f Files) Dash() (*File, bool) {
   for _, fileVar := range f {
      if fileVar.Method == "dash" {
         return &fileVar, true
      }
   }
   return nil, false
}
