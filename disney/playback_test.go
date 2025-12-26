package disney

import (
   "41.neocities.org/drm/widevine"
   "encoding/xml"
   "fmt"
   "os"
   "testing"
)

func TestStream(t *testing.T) {
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(cache + "/disney/account.xml")
   if err != nil {
      t.Fatal(err)
   }
   var account_with account
   err = xml.Unmarshal(data, &account_with)
   if err != nil {
      t.Fatal(err)
   }
   explore, err := account_with.explore(test.entity)
   if err != nil {
      t.Fatal(err)
   }
   resource_id, ok := explore.restart()
   if !ok {
      t.Fatal(".restart()")
   }
   play, err := account_with.playback(resource_id)
   if err != nil {
      t.Fatal(err)
   }
   for i, source := range play.Stream.Sources {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Printf("%+v\n", source)
   }
}

var key_id = []byte{
   188, 54, 159, 224, 114, 252, 64, 161, 184, 218, 28, 219, 235, 253, 0, 105,
}

func TestWidevine(t *testing.T) {
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
   data, err := os.ReadFile(cache + "/disney/account.xml")
   if err != nil {
      t.Fatal(err)
   }
   var account_with account
   err = xml.Unmarshal(data, &account_with)
   if err != nil {
      t.Fatal(err)
   }
   _, err = account_with.widevine(msg)
   if err != nil {
      t.Fatal(err)
   }
}

var test = struct {
   entity string
   url    string
}{
   entity: "7df81cf5-6be5-4e05-9ff6-da33baf0b94d",
   url:    "https://disneyplus.com/play/7df81cf5-6be5-4e05-9ff6-da33baf0b94d",
}
