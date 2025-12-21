package canal

import (
   "crypto/hmac"
   "crypto/sha256"
   "encoding/base64"
   "io"
   "net/url"
   "strconv"
   "strings"
   "time"
)

func GetClient(link *url.URL, body []byte) (string, error) {
   encoding := base64.RawURLEncoding
   // 1. base64 raw URL decode secret key
   decoded_key, err := encoding.DecodeString(secret_key)
   if err != nil {
      return "", err
   }
   // Prepare timestamp and hash the body
   timestamp := time.Now().Unix()
   body_checksum := sha256.Sum256(body)
   encoded_body_hash := encoding.EncodeToString(body_checksum[:])
   // 2. hmac.New(sha256.New, secret key)
   hash := hmac.New(sha256.New, decoded_key)
   // 3, 4, 5. Write components to the hasher
   // Instead of fmt.Fprint, write parts sequentially.
   io.WriteString(hash, link.String())
   io.WriteString(hash, encoded_body_hash)
   // Convert int64 timestamp to decimal string and write
   io.WriteString(hash, strconv.FormatInt(timestamp, 10))
   // 6. base64 raw URL encode the hmac sum
   signature := encoding.EncodeToString(hash.Sum(nil))
   // Construct final result string without "+" or fmt.Sprintf
   var sb strings.Builder
   sb.WriteString("Client key=")
   sb.WriteString(client_key)
   sb.WriteString(",time=")
   sb.WriteString(strconv.FormatInt(timestamp, 10))
   sb.WriteString(",sig=")
   sb.WriteString(signature)
   return sb.String(), nil
}

// Global variables for authentication
var (
   client_key = "web.NhFyz4KsZ54"
   secret_key = "OXh0-pIwu3gEXz1UiJtqLPscZQot3a0q"
)
