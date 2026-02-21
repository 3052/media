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
   "io"
   "net/http"
   "net/url"
   "slices"
   "strconv"
   "strings"
)

type Dash struct {
   Body []byte
   Url  *url.URL
}

type Item struct {
   CmsAccountId string
   ContentId    string
}

// WARNING IF YOU RUN THIS TOO MANY TIMES YOU WILL GET AN IP BAN. HOWEVER THE BAN
// IS ONLY FOR THE ANDROID CLIENT NOT WEB CLIENT
func Login(at, username, password string) (*http.Cookie, error) {
   data := url.Values{
      "j_username": {username},
      "j_password": {password},
   }.Encode()
   var req http.Request
   req.Method = "POST"
   req.URL = &url.URL{
      Scheme:   "https",
      Host:     "www.paramountplus.com",
      Path:     "/apps-api/v2.0/androidphone/auth/login.json",
      RawQuery: url.Values{"at": {at}}.Encode(),
   }
   req.Header = http.Header{}
   req.Header.Set("content-type", "application/x-www-form-urlencoded")
   // randomly fails if this is missing
   req.Header.Set("user-agent", "!")
   req.Body = io.NopCloser(strings.NewReader(data))
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   _, err = io.Copy(io.Discard, resp.Body)
   if err != nil {
      return nil, err
   }
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   for _, cookie := range resp.Cookies() {
      if cookie.Name == "CBS_COM" {
         return cookie, nil
      }
   }
   return nil, http.ErrNoCookie
}

func FetchItem(at, cId string) (*Item, error) {
   var req http.Request
   req.Header = http.Header{}
   req.URL = &url.URL{
      Scheme:   "https",
      Host:     "www.paramountplus.com",
      Path:     join("/apps-api/v2.0/androidphone/video/cid/", cId, ".json"),
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

// 1080p SL2000
// 1440p SL2000 + cookie
func PlayReady(at, contentId string, cookie *http.Cookie) (*SessionToken, error) {
   var req http.Request
   req.Header = http.Header{}
   req.URL = &url.URL{}
   req.URL.Scheme = "https"
   req.URL.Host = "www.paramountplus.com"
   req.URL.RawQuery = url.Values{
      "at":        {at},
      "contentId": {contentId},
   }.Encode()
   if cookie != nil {
      req.AddCookie(cookie)
      req.URL.Path = "/apps-api/v3.1/xboxone/irdeto-control/session-token.json"
   } else {
      req.URL.Path = "/apps-api/v3.1/xboxone/irdeto-control/anonymous-session-token.json"
   }
   resp, err := http.DefaultClient.Do(&req)
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
   var result SessionToken
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result, nil
}

// 576p L3
func Widevine(at, contentId string) (*SessionToken, error) {
   var req http.Request
   req.URL = &url.URL{}
   req.URL.Scheme = "https"
   req.URL.Host = "www.paramountplus.com"
   req.URL.Path = "/apps-api/v3.1/androidphone/irdeto-control/anonymous-session-token.json"
   req.URL.RawQuery = url.Values{
      "at":        {at},
      "contentId": {contentId},
   }.Encode()
   req.Header = http.Header{}
   resp, err := http.DefaultClient.Do(&req)
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
   var result SessionToken
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result, nil
}

func (s *SessionToken) Send(data []byte) ([]byte, error) {
   req, err := http.NewRequest("POST", s.Url, bytes.NewReader(data))
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+s.LsSession)
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

func join(items ...string) string {
   return strings.Join(items, "")
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

type SessionToken struct {
   LsSession string `json:"ls_session"`
   Url       string
}

func (i *Item) Dash() (*Dash, error) {
   var req http.Request
   req.Header = http.Header{}
   req.URL = &url.URL{
      Scheme: "https",
      Host:   "link.theplatform.com",
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
   if resp.StatusCode != http.StatusOK {
      var value struct {
         Description string
      }
      err = json.Unmarshal(data, &value)
      if err != nil {
         return nil, err
      }
      return nil, errors.New(value.Description)
   }
   return &Dash{Body: data, Url: resp.Request.URL}, nil
}
