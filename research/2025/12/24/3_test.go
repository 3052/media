package disney

import (
   "41.neocities.org/drm/widevine"
   "encoding/xml"
   "os"
   "testing"
)

var key_id = []byte{188, 54, 159, 224, 114, 252, 64, 161, 184, 218, 28, 219, 235, 253, 0, 105}

func TestObtainLicense(t *testing.T) {
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   client_id, err := os.ReadFile(cache + "/L3/client_id.bin")
   if err != nil {
      t.Fatal(err)
   }
   pem_bytes, err := os.ReadFile(cache + "/L3/private_key.pem")
   if err != nil {
      t.Fatal(err)
   }
   private_key, err := widevine.ParsePrivateKey(pem_bytes)
   if err != nil {
      t.Fatal(err)
   }
   var pssh widevine.PsshData
   pssh.KeyIds = [][]byte{key_id}
   msg, err := pssh.BuildLicenseRequest(client_id)
   if err != nil {
      t.Fatal(err)
   }
   msg, err = widevine.BuildSignedMessage(msg, private_key)
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(cache + "/disney/refresh_token.xml")
   if err != nil {
      t.Fatal(err)
   }
   var token refresh_token
   err = xml.Unmarshal(data, &token)
   if err != nil {
      t.Fatal(err)
   }
   _, err = token.obtain_license(msg)
   if err != nil {
      t.Fatal(err)
   }
}
