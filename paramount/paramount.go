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

var AppSecrets = []struct {
   Version   string
   ComCbsApp string
   ComCbsCa  string
}{
   {
      Version: "16.4.1",
      ComCbsApp: "7cd07f93a6e44cf7",
      ComCbsCa: "68b4475a49bed95a",
   },
   {
      Version:   "16.0.0",
      ComCbsApp: "9fc14cb03691c342",
      ComCbsCa:  "6c68178445de8138",
   },
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

type Dash struct {
   Body []byte
   Url  *url.URL
}

type Item struct {
   CmsAccountId string
   ContentId    string
}

func (i *Item) Dash() (*Dash, error) {
   var req http.Request
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
   req.Header = http.Header{}
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Dash
   result.Body, err = io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   result.Url = resp.Request.URL
   return &result, nil
}

// https://paramountplus.com
// https://paramountplus.com/movies/video/wjQ4RChi6BHHu4MVTncppVuCwu44uq2Q
func FetchPath(address string) (string, error) {
   resp, err := http.Head(address)
   if err != nil {
      return "", err
   }
   return resp.Request.URL.Path, nil
}

// WARNING IF YOU RUN THIS TOO MANY TIMES YOU WILL GET AN IP BAN. HOWEVER THE BAN
// IS ONLY FOR THE ANDROID CLIENT NOT WEB CLIENT
func (c *Content) Login(username, password string) (*http.Cookie, error) {
   at, err := GetAt(c.AppSecret())
   if err != nil {
      return nil, err
   }
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
   for _, cookie := range resp.Cookies() {
      if cookie.Name == "CBS_COM" {
         return cookie, nil
      }
   }
   return nil, http.ErrNoCookie
}

// Content holds extracted data from a Paramount+ path.
type Content struct {
   ID          string
   CountryCode string
   Type        string // "movies" or "shows"
}

// Parse populates the Content struct by parsing the given Paramount+ path.
// It assumes the receiver (c) is a zero-value struct.
func (c *Content) Parse(path string) error {
   // 1. Trim both leading and trailing slashes for clean processing.
   cleanPath := strings.Trim(path, "/")
   // 2. Handle the root path case (e.g., "/").
   // An empty path after trimming is valid and results in a zero-value struct.
   if cleanPath == "" {
      return nil
   }
   parts := strings.Split(cleanPath, "/")
   // 3. Handle the region-only case (e.g., "/ie/").
   // This is a single path component that is two letters long.
   if len(parts) == 1 && len(parts[0]) == 2 {
      c.CountryCode = parts[0]
      // ID and Type correctly remain empty for this case.
      return nil
   }
   // 4. Handle paths that must contain a content ID.
   // These paths must have at least 3 components.
   if len(parts) < 3 {
      return errors.New("invalid path: not enough components for a content path")
   }
   // 5. The ID is always the last part.
   c.ID = parts[len(parts)-1]
   if c.ID == "" {
      return errors.New("invalid path: ID is missing")
   }
   // 6. Determine the type and country code based on path structure.
   if len(parts) >= 4 && len(parts[0]) == 2 {
      // Structure with country code: [cc, type, slug, id]
      c.CountryCode = parts[0]
      c.Type = parts[1]
   } else {
      // Structure without country code: [type, slug, id]
      c.Type = parts[0]
   }
   // 7. Validate the assigned type.
   if c.Type != "movies" && c.Type != "shows" {
      return errors.New("invalid content type")
   }
   return nil // Success
}

func (c *Content) AppSecret() string {
   if c.CountryCode != "" {
      return AppSecrets[0].ComCbsCa
   }
   return AppSecrets[0].ComCbsApp
}

///

func (c *Content) Item() (*Item, error) {
   at, err := GetAt(c.AppSecret())
   if err != nil {
      return nil, err
   }
   var req http.Request
   req.URL = &url.URL{
      Scheme:   "https",
      Host:     "www.paramountplus.com",
      Path:     join("/apps-api/v2.0/androidphone/video/cid/", c.Id, ".json"),
      RawQuery: url.Values{"at": {at}}.Encode(),
   }
   req.Header = http.Header{}
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK { // error 403 406
      return nil, errors.New(resp.Status)
   }
   var result struct {
      ItemList []Item
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.ItemList) == 0 { // error 200
      return nil, errors.New("item list zero length")
   }
   return &result.ItemList[0], nil
}

// 576p L3
func (c *Content) Widevine() (*SessionToken, error) {
   at, err := GetAt(c.AppSecret())
   if err != nil {
      return nil, err
   }
   var req http.Request
   req.URL = &url.URL{}
   req.URL.Scheme = "https"
   req.URL.Host = "www.paramountplus.com"
   req.URL.Path = "/apps-api/v3.1/androidphone/irdeto-control/anonymous-session-token.json"
   req.URL.RawQuery = url.Values{
      "at":        {at},
      "contentId": {c.Id},
   }.Encode()
   req.Header = http.Header{}
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result SessionToken
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result, nil
}

// 1080p SL2000
// 1440p SL2000 + cookie
func (c *Content) PlayReady(cookie *http.Cookie) (*SessionToken, error) {
   at, err := GetAt(c.AppSecret())
   if err != nil {
      return nil, err
   }
   var req http.Request
   req.URL = &url.URL{}
   req.URL.Scheme = "https"
   req.URL.Host = "www.paramountplus.com"
   req.URL.RawQuery = url.Values{
      "at":        {at},
      "contentId": {c.Id},
   }.Encode()
   if cookie != nil {
      req.URL.Path = "/apps-api/v3.1/xboxone/irdeto-control/session-token.json"
      req.AddCookie(cookie)
   } else {
      req.URL.Path = "/apps-api/v3.1/xboxone/irdeto-control/anonymous-session-token.json"
   }
   req.Header = http.Header{}
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var result SessionToken
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result, nil
}
