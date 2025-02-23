package paramount

import (
   "bytes"
   "crypto/aes"
   "crypto/cipher"
   "encoding/base64"
   "encoding/hex"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

type AppSecret string

// 15.0.52
const ComCbsApp AppSecret = "4fb47ec1f5c17caa"

// 15.0.52
const ComCbsCa AppSecret = "e55edaeb8451f737"

func (a At) Session(content_id string) (*Session, error) {
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
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      var b strings.Builder
      resp.Write(&b)
      return nil, errors.New(b.String())
   }
   session1 := &Session{}
   err = json.NewDecoder(resp.Body).Decode(session1)
   if err != nil {
      return nil, err
   }
   return session1, nil
}

const secret_key = "302a6a0d70a7e9b967f91d39fef3e387816e3095925ae4537bce96063311f9c5"

func pad(data []byte) []byte {
   length := aes.BlockSize - len(data)%aes.BlockSize
   for high := byte(length); length >= 1; length-- {
      data = append(data, high)
   }
   return data
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

type Session struct {
   LsSession string `json:"ls_session"`
   Url       string
}

type Item struct {
   AssetType    string
   CmsAccountId string
   ContentId    string
}

func (s *Session) Widevine() func([]byte) ([]byte, error) {
   return func(data []byte) ([]byte, error) {
      req, err := http.NewRequest("POST", s.Url, bytes.NewReader(data))
      if err != nil {
         return nil, err
      }
      req.Header = http.Header{
         "authorization": {"Bearer " + s.LsSession},
         "content-type":  {"application/x-protobuf"},
      }
      resp, err := http.DefaultClient.Do(req)
      if err != nil {
         return nil, err
      }
      defer resp.Body.Close()
      return io.ReadAll(resp.Body)
   }
}

type At string

func (a AppSecret) At() (At, error) {
   key, err := hex.DecodeString(secret_key)
   if err != nil {
      return "", err
   }
   block, err := aes.NewCipher(key)
   if err != nil {
      return "", err
   }
   var iv [aes.BlockSize]byte
   data := []byte{'|'}
   data = append(data, a...)
   data = pad(data)
   cipher.NewCBCEncrypter(block, iv[:]).CryptBlocks(data, data)
   data1 := []byte{0, aes.BlockSize}
   data1 = append(data1, iv[:]...)
   data1 = append(data1, data...)
   return At(base64.StdEncoding.EncodeToString(data1)), nil
}

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
   req.Header.Set("vpn", "true")
   return http.DefaultClient.Do(req)
}

func (a At) Item(cid string) (*Item, error) {
   req, _ := http.NewRequest("", "https://www.paramountplus.com", nil)
   req.URL.Path = func() string {
      var b strings.Builder
      b.WriteString("/apps-api/v2.0/androidphone/video/cid/")
      b.WriteString(cid)
      b.WriteString(".json")
      return b.String()
   }()
   req.URL.RawQuery = "at=" + string(a)
   req.Header.Set("vpn", "true")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   data, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   var value struct {
      ItemList []Item
   }
   err = json.Unmarshal(data, &value)
   if err != nil {
      return nil, err
   }
   if len(value.ItemList) == 0 {
      return nil, errors.New(string(data))
   }
   return &value.ItemList[0], nil
}
