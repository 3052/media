package disney

import (
   "41.neocities.org/drm/widevine"
   "encoding/base64"
   "encoding/xml"
   "os"
   "testing"
   "time"
)

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func output(name string, arg ...string) (string, error) {
   data, err := exec.Command(name, arg...).Output()
   if err != nil {
      return "", err
   }
   return string(data), nil
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
