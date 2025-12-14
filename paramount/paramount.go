package paramount

import (
   "bytes"
   "crypto/aes"
   "crypto/cipher"
   "encoding/base64"
   "encoding/binary"
   "encoding/hex"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "slices"
   "strings"
)

const (
   encoding = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
   secret_key = "302a6a0d70a7e9b967f91d39fef3e387816e3095925ae4537bce96063311f9c5"
)

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

func GetAt(appSecret string) (string, error) {
   // 1. Decode hex secret key
   key, err := hex.DecodeString(secret_key)
   if err != nil {
      return "", err
   }
   // 2. Create aes cipher with key
   block, err := aes.NewCipher(key)
   if err != nil {
      return "", err
   }
   // 3 & 4. Create payload: "|" + appSecret
   data := []byte{'|'}
   data = append(data, appSecret...)
   // 5. Apply PKCS7 Padding (Separate Function)
   data = pkcs7_pad(data, aes.BlockSize)
   // Prepare Empty IV (16 bytes of zeros)
   var iv [aes.BlockSize]byte
   // 6. CBC encrypt with empty IV
   // We encrypt 'data' in place
   cipher.NewCBCEncrypter(block, iv[:]).CryptBlocks(data, data)
   // 8. Create Header for block size (uint16)
   size := binary.BigEndian.AppendUint16(nil, aes.BlockSize)
   // 7 & 8. Combine [Size] + [IV] + [Encrypted Data]
   data = slices.Concat(size[:], iv[:], data)
   // 9. Return result base64 encoded
   return base64.StdEncoding.EncodeToString(data), nil
}

func pkcs7_pad(data []byte, blockSize int) []byte {
   // Calculate the number of padding bytes needed.
   // If data is already a multiple of blockSize, this results in a full block
   // of padding.
   paddingLen := blockSize - (len(data) % blockSize)
   // Create a padding byte (the value is the length of the padding)
   padByte := byte(paddingLen)
   // Append the padding byte 'paddingLen' times
   for i := 0; i < paddingLen; i++ {
      data = append(data, padByte)
   }
   return data
}

type Item struct {
   AssetType    string
   CmsAccountId string
   ContentId    string
}

func FetchItem(at, contentId string) (*Item, error) {
   req, _ := http.NewRequest("", "https://www.paramountplus.com", nil)
   req.URL.Path = func() string {
      var data strings.Builder
      data.WriteString("/apps-api/v2.0/androidphone/video/cid/")
      data.WriteString(contentId)
      data.WriteString(".json")
      return data.String()
   }()
   req.URL.RawQuery = url.Values{
      "at": {at},
   }.Encode()
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
   var result struct {
      ItemList []Item
   }
   err = json.Unmarshal(data, &result)
   if err != nil {
      return nil, err
   }
   if len(result.ItemList) == 0 { // error 200
      return nil, errors.New(string(data))
   }
   return &result.ItemList[0], nil
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

type Provider struct {
   AppSecret string
   Version   string
}

var ComCbsApp = Provider{
   AppSecret: "9fc14cb03691c342",
   Version:   "16.0.0",
}

var ComCbsCa = Provider{
   AppSecret: "6c68178445de8138",
   Version:   "16.0.0",
}

type SessionToken struct {
   LsSession string `json:"ls_session"`
   Url       string
}

func (s *SessionToken) Widevine(data []byte) ([]byte, error) {
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

func (s *SessionToken) playReady(at, contentId string) error {
   req, _ := http.NewRequest("", "https://www.paramountplus.com", nil)
   req.URL.Path = func() string {
      var data strings.Builder
      data.WriteString("/apps-api/v3.1/xboxone/irdeto-control")
      data.WriteString("/anonymous-session-token.json")
      return data.String()
   }()
   req.URL.RawQuery = url.Values{
      "at":        {at},
      "contentId": {contentId},
   }.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(s)
}

func (s *SessionToken) Fetch(at, contentId string) error {
   req, _ := http.NewRequest("", "https://www.paramountplus.com", nil)
   req.URL.Path = func() string {
      var data strings.Builder
      data.WriteString("/apps-api/v3.1/androidphone/irdeto-control")
      data.WriteString("/anonymous-session-token.json")
      return data.String()
   }()
   req.URL.RawQuery = url.Values{
      "at":        {at},
      "contentId": {contentId},
   }.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   if resp.StatusCode != http.StatusOK {
      var data strings.Builder
      err = resp.Write(&data)
      if err != nil {
         return err
      }
      return errors.New(data.String())
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(s)
}
