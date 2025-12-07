package paramount

import (
   "bytes"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strings"
)

func (a At) Token(content_id string) (*Token, error) {
   req, _ := http.NewRequest("", "https://www.paramountplus.com", nil)
   req.URL.Path = func() string {
      var data strings.Builder
      data.WriteString("/apps-api/v3.1/androidphone/irdeto-control")
      data.WriteString("/anonymous-session-token.json")
      return data.String()
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
      var data strings.Builder
      err = resp.Write(&data)
      if err != nil {
         return nil, err
      }
      return nil, errors.New(data.String())
   }
   defer resp.Body.Close()
   session := &Token{}
   err = json.NewDecoder(resp.Body).Decode(session)
   if err != nil {
      return nil, err
   }
   return session, nil
}

func (t *Token) Widevine(data []byte) ([]byte, error) {
   req, err := http.NewRequest("POST", t.Url, bytes.NewReader(data))
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+t.LsSession)
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

func (a At) playReady(content_id string) (*Token, error) {
   req, _ := http.NewRequest("", "https://www.paramountplus.com", nil)
   req.URL.Path = func() string {
      var data strings.Builder
      data.WriteString("/apps-api/v3.1/xboxone/irdeto-control")
      data.WriteString("/anonymous-session-token.json")
      return data.String()
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
   session_var := &Token{}
   err = json.NewDecoder(resp.Body).Decode(session_var)
   if err != nil {
      return nil, err
   }
   return session_var, nil
}

type Token struct {
   LsSession string `json:"ls_session"`
   Url       string
}

func (a At) Item(cid string) (*Item, error) {
   req, _ := http.NewRequest("", "https://www.paramountplus.com", nil)
   req.URL.Path = func() string {
      var data strings.Builder
      data.WriteString("/apps-api/v2.0/androidphone/video/cid/")
      data.WriteString(cid)
      data.WriteString(".json")
      return data.String()
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
type Item struct {
   AssetType    string
   CmsAccountId string
   ContentId    string
}

const encoding = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

func cms_account(id string) int {
   var (
      result = 0
      multiplier = 1
   )
   for _, digit := range id {
      result += strings.IndexRune(encoding, digit) * multiplier
      multiplier *= len(encoding)
   }
   return result
}
func (i *Item) Mpd() (*url.URL, []byte, error) {
   req, _ := http.NewRequest("", "https://link.theplatform.com", nil)
   req.URL.Path = func() string {
      data := []byte("/s/")
      data = append(data, i.CmsAccountId...)
      data = append(data, "/media/guid/"...)
      data = fmt.Append(data, cms_account(i.CmsAccountId))
      data = append(data, '/')
      data = append(data, i.ContentId...)
      return string(data)
   }()
   req.URL.RawQuery = url.Values{
      "assetTypes": {i.AssetType},
      "formats":    {"MPEG-DASH"},
   }.Encode()
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
