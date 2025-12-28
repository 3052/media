package disney

import (
   "41.neocities.org/drm/widevine"
   "encoding/base64"
   "encoding/xml"
   "os"
   "testing"
   "time"
)

var tests = []struct {
   entity string
   key_ids []string
   url    string
}{
   {
      entity: "7df81cf5-6be5-4e05-9ff6-da33baf0b94d",
      key_ids: []string{
         "vDaf4HL8QKG42hzb6/0AaQ==", // L3
         "GzlnsDJgQ82Du9PjayM4IQ==", // L1
      },
      url: "https://disneyplus.com/play/7df81cf5-6be5-4e05-9ff6-da33baf0b94d",
   },
   {
      entity: "d0d0796c-a144-42fa-a730-4cbd1014ef1f",
      url: "https://disneyplus.com/browse/entity-d0d0796c-a144-42fa-a730-4cbd1014ef1f",
   },
}

func TestPlayback(t *testing.T) {
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
   for i, test := range tests {
      if i >= 1 {
         time.Sleep(time.Second)
      }
      explore_item, err := account_with.explore(test.entity)
      if err != nil {
         t.Fatal(err)
      }
      resource_id, ok := explore_item.restart()
      if !ok {
         t.Fatal(".restart()")
      }
      play, err := account_with.playback(resource_id)
      if err != nil {
         t.Fatal(err)
      }
      t.Logf("%+v", play.Stream.Sources[0])
   }
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
   key_id, err := base64.StdEncoding.DecodeString(tests[0].key_ids[0])
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
