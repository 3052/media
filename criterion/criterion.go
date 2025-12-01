package criterion

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
)

const client_id = "9a87f110f79cd25250f6c7f3a6ec8b9851063ca156dae493bf362a7faf146c78"

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
   AccessToken      string `json:"access_token"`
   ErrorDescription string `json:"error_description"`
   RefreshToken     string `json:"refresh_token"`
}

func (t *Token) Unmarshal(data TokenData) error {
   err := json.Unmarshal(data, t)
   if err != nil {
      return err
   }
   if t.ErrorDescription != "" {
      return errors.New(t.ErrorDescription)
   }
   return nil
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
   var video_var Video
   err = json.NewDecoder(resp.Body).Decode(&video_var)
   if err != nil {
      return nil, err
   }
   if video_var.Message != "" {
      return nil, errors.New(video_var.Message)
   }
   return &video_var, nil
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
   var files_var Files
   err = json.NewDecoder(resp.Body).Decode(&files_var)
   if err != nil {
      return nil, err
   }
   return files_var, nil
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
   for _, file_var := range f {
      if file_var.Method == "dash" {
         return &file_var, true
      }
   }
   return nil, false
}
