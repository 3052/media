package disney

import (
   "41.neocities.org/drm/widevine"
   "encoding/base64"
   "encoding/xml"
   "os"
   "testing"
   "time"
)

func TestExplore(t *testing.T) {
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(cache + "/disney/account.xml")
   if err != nil {
      t.Fatal(err)
   }
   var account_with Account
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
      t.Logf("%+v", explore_item)
   }
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
   var account_with Account
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
      resource_id, ok := explore_item.play_restart()
      if !ok {
         t.Fatal(".play_restart()")
      }
      play, err := account_with.Playback(resource_id)
      if err != nil {
         t.Fatal(err)
      }
      t.Logf("%+v", play.Stream.Sources[0])
   }
}

var tests = []struct {
   entity   string
   key_ids  []string
   location string
   url      string
}{
   {
      entity: "7df81cf5-6be5-4e05-9ff6-da33baf0b94d",
      key_ids: []string{
         "vDaf4HL8QKG42hzb6/0AaQ==", // L3
         "GzlnsDJgQ82Du9PjayM4IQ==", // L1
      },
      location: "US",
      url:      "https://disneyplus.com/play/7df81cf5-6be5-4e05-9ff6-da33baf0b94d",
   },
   {
      entity: "d0d0796c-a144-42fa-a730-4cbd1014ef1f",
      key_ids: []string{
         "UKmxApOuRGKsyP+co9ABog==",
         "nYuSQdxPSJ67kNT6coO3cA==",
      },
      location: "KR", // MUST DO KR LOGIN FIRST
      url:      "https://disneyplus.com/browse/entity-d0d0796c-a144-42fa-a730-4cbd1014ef1f",
   },
}

func TestWidevine(t *testing.T) {
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(cache + "/disney/account.xml")
   if err != nil {
      t.Fatal(err)
   }
   var account_with Account
   err = xml.Unmarshal(data, &account_with)
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
   original_request, err := pssh.BuildLicenseRequest(client_id)
   if err != nil {
      t.Fatal(err)
   }
   data, err = widevine.BuildSignedMessage(original_request, private_key)
   if err != nil {
      t.Fatal(err)
   }
   data, err = account_with.widevine(data)
   if err != nil {
      t.Fatal(err)
   }
   keys, err := widevine.ParseLicenseResponse(
      data, original_request, private_key,
   )
   if err != nil {
      t.Fatal(err)
   }
   for _, key := range keys {
      t.Logf("Id:%x", key.Id)
      t.Logf("Key:%x", key.Key)
   }
}
