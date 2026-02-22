package peacock

import (
   "bytes"
   "crypto/hmac"
   "crypto/md5"
   "crypto/sha1"
   "encoding/base64"
   "fmt"
   "net/http"
   "slices"
   "strings"
   "time"
)

const (
   sky_client  = "NBCU-ANDROID-v3"
   sky_key     = "JuLQgyFz9n89D9pxcN6ZWZXKWfgj2PNBUb32zybj"
   sky_version = "1.0"
)

func generate_sky_ott(method, path string, headers http.Header, body []byte) string {
   // Sort headers by key.
   headerKeys := make([]string, 0, len(headers))
   for key := range headers {
      headerKeys = append(headerKeys, key)
   }
   slices.Sort(headerKeys)
   // Build the special headers string.
   var headersBuilder bytes.Buffer
   for _, key := range headerKeys {
      lowerKey := strings.ToLower(key)
      if strings.HasPrefix(lowerKey, "x-skyott-") {
         value := headers.Get(key)
         headersBuilder.WriteString(lowerKey)
         headersBuilder.WriteString(": ")
         headersBuilder.WriteString(value)
         headersBuilder.WriteByte('\n')
      }
   }
   // MD5 the headers string and request body.
   headersHash := md5.Sum(headersBuilder.Bytes())
   headersMD5 := fmt.Sprintf("%x", headersHash)
   bodyHash := md5.Sum(body)
   bodyMD5 := fmt.Sprintf("%x", bodyHash)
   // Get current timestamp string directly.
   timestampStr := fmt.Sprint(time.Now().Unix())
   // Construct the payload to be signed for the HMAC.
   var payload bytes.Buffer
   payload.WriteString(method)
   payload.WriteByte('\n')
   payload.WriteString(path)
   payload.WriteByte('\n')
   payload.WriteByte('\n')
   payload.WriteString(sky_client)
   payload.WriteByte('\n')
   payload.WriteString(sky_version)
   payload.WriteByte('\n')
   payload.WriteString(headersMD5)
   payload.WriteByte('\n')
   payload.WriteString(timestampStr)
   payload.WriteByte('\n')
   payload.WriteString(bodyMD5)
   payload.WriteByte('\n')
   // Calculate the HMAC signature.
   mac := hmac.New(sha1.New, []byte(sky_key))
   payload.WriteTo(mac)
   signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
   // Format the final output string.
   return fmt.Sprintf(
      "SkyOTT client=%q,signature=%q,timestamp=%q,version=%q",
      sky_client,
      signature,
      timestampStr,
      sky_version,
   )
}
