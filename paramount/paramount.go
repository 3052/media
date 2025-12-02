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

type Cache struct {
   Mpd      *url.URL
   MpdBody  []byte
}

func (i *Item) Mpd(storage *Cache) error {
   req, _ := http.NewRequest("", "https://link.theplatform.com", nil)
   req.URL.Path = func() string {
      data := []byte("/s/")
      data = append(data, i.CmsAccountId...)
      data = append(data, "/media/guid/"...)
      data = strconv.AppendInt(data, cms_account(i.CmsAccountId), 10)
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
      return err
   }
   defer resp.Body.Close()
   storage.MpdBody, err = io.ReadAll(resp.Body)
   if err != nil {
      return err
   }
   storage.Mpd = resp.Request.URL
   return nil
}

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

type At string

const secret_key = "302a6a0d70a7e9b967f91d39fef3e387816e3095925ae4537bce96063311f9c5"

// 16.0.0
const ComCbsCa AppSecret = "6c68178445de8138"

// 16.0.0
const ComCbsApp AppSecret = "9fc14cb03691c342"

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
   data = pkcs7_padding(data, aes.BlockSize)
   cipher.NewCBCEncrypter(block, iv[:]).CryptBlocks(data, data)
   data1 := []byte{0, aes.BlockSize}
   data1 = append(data1, iv[:]...)
   data1 = append(data1, data...)
   return At(base64.StdEncoding.EncodeToString(data1)), nil
}

type AppSecret string

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

func pkcs7_padding(data []byte, blockSize int) []byte {
   padLen := blockSize - (len(data) % blockSize)
   for i := 0; i < padLen; i++ {
      data = append(data, byte(padLen))
   }
   return data
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

type Item struct {
   AssetType    string
   CmsAccountId string
   ContentId    string
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
