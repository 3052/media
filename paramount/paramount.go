package paramount

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

func (i *Item) Mpd() (*http.Response, error) {
   req, _ := http.NewRequest("", "https://link.theplatform.com", nil)
   req.URL.Path = func() string {
      b := []byte("/s/")
      b = append(b, i.CmsAccountId...)
      b = append(b, "/media/guid/"...)
      b = strconv.AppendInt(b, cms_account(i.CmsAccountId), 10)
      b = append(b, '/')
      b = append(b, i.ContentId...)
      return string(b)
   }()
   req.URL.RawQuery = url.Values{
      "assetTypes": {i.AssetType},
      "formats":    {"MPEG-DASH"},
   }.Encode()
   return http.DefaultClient.Do(req)
}

type Item struct {
   AssetType    string
   CmsAccountId string
   ContentId    string
}

const encoding = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

func cms_account(id string) int64 {
   var (
      i = 0
      j = 1
   )
   for _, value := range id {
      i += strings.IndexRune(encoding, value) * j
      j *= len(encoding)
   }
   return int64(i)
}

func (s *Session) Widevine(data []byte) ([]byte, error) {
   req, err := http.NewRequest("POST", s.Url, bytes.NewReader(data))
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+s.LsSession)
   req.Header.Set("content-type", "application/x-protobuf")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   data, err = io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(string(data))
   }
   return data, nil
}

func (a At) playReady(content_id string) (*Session, error) {
   req, _ := http.NewRequest("", "https://www.paramountplus.com", nil)
   req.URL.Path = func() string {
      var b strings.Builder
      b.WriteString("/apps-api/v3.1/xboxone/irdeto-control")
      b.WriteString("/anonymous-session-token.json")
      return b.String()
   }()
   req.URL.RawQuery = url.Values{
      "at":        {string(a)},
      "contentId": {content_id},
   }.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   sessionVar := &Session{}
   err = json.NewDecoder(resp.Body).Decode(sessionVar)
   if err != nil {
      return nil, err
   }
   return sessionVar, nil
}

type Session struct {
   LsSession string `json:"ls_session"`
   Url       string
}

// proxy
func (a At) Item(cid string) (*Item, error) {
   req, _ := http.NewRequest("", "https://www.paramountplus.com", nil)
   req.URL.Path = func() string {
      var b strings.Builder
      b.WriteString("/apps-api/v2.0/androidphone/video/cid/")
      b.WriteString(cid)
      b.WriteString(".json")
      return b.String()
   }()
   req.URL.RawQuery = "at=" + url.QueryEscape(string(a))
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   data, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   if resp.StatusCode != http.StatusOK { // error 403 406
      if len(data) >= 1 {
         return nil, errors.New(string(data))
      }
      return nil, errors.New(resp.Status)
   }
   var value struct {
      ItemList []Item
   }
   err = json.Unmarshal(data, &value)
   if err != nil {
      return nil, err
   }
   if len(value.ItemList) == 0 { // error 200
      return nil, errors.New(string(data))
   }
   return &value.ItemList[0], nil
}

func (a At) Session(content_id string) (Byte[Session], error) {
   req, _ := http.NewRequest("", "https://www.paramountplus.com", nil)
   req.URL.Path = func() string {
      var b strings.Builder
      b.WriteString("/apps-api/v3.1/androidphone/irdeto-control")
      b.WriteString("/anonymous-session-token.json")
      return b.String()
   }()
   req.URL.RawQuery = url.Values{
      "at":        {string(a)},
      "contentId": {content_id},
   }.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   if resp.StatusCode != http.StatusOK {
      var b strings.Builder
      err = resp.Write(&b)
      if err != nil {
         return nil, err
      }
      return nil, errors.New(b.String())
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

type Byte[T any] []byte

func (s *Session) Unmarshal(data Byte[Session]) error {
   return json.Unmarshal(data, s)
}
