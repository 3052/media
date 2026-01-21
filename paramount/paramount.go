package paramount

import (
   "crypto/aes"
   "crypto/cipher"
   "encoding/base64"
   "encoding/binary"
   "encoding/hex"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "slices"
   "strconv"
   "strings"
)

func (i *Item) Mpd() (*Mpd, error) {
   var req http.Request
   req.Header = http.Header{}
   req.URL = &url.URL{
      Scheme: "https",
      Host: "link.theplatform.com",
      Path: join(
         "/s/", i.CmsAccountId,
         "/media/guid/", strconv.Itoa(cms_account(i.CmsAccountId)),
         "/", i.ContentId,
      ),
      RawQuery: "formats=MPEG-DASH",
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   data, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   return &Mpd{data, resp.Request.URL}, nil
}

type Item struct {
   CmsAccountId string
   ContentId    string
}

func join(items ...string) string {
   return strings.Join(items, "")
}

type Mpd struct {
   Body []byte
   Url  *url.URL
}

const (
   encoding   = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
   secret_key = "302a6a0d70a7e9b967f91d39fef3e387816e3095925ae4537bce96063311f9c5"
)

func cms_account(id string) int {
   var (
      result     = 0
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

type Provider struct {
   AppSecret string
   Version   string
}

type SessionToken struct {
   LsSession string `json:"ls_session"`
   Url       string
}
var ComCbsApp = Provider{
   AppSecret: "9fc14cb03691c342",
   Version:   "16.0.0",
}

var ComCbsCa = Provider{
   AppSecret: "6c68178445de8138",
   Version:   "16.0.0",
}

func FetchItem(at, cId string) (*Item, error) {
   var req http.Request
   req.Header = http.Header{}
   req.URL = &url.URL{
      Scheme: "https",
      Host: "www.paramountplus.com",
      Path: join("/apps-api/v2.0/androidphone/video/cid/", cId, ".json"),
      RawQuery: url.Values{"at": {at}}.Encode(),
   }
   resp, err := http.DefaultClient.Do(&req)
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

