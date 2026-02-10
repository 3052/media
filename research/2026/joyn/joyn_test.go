package joyn

import (
   "154.pages.dev/widevine"
   "encoding/base64"
   "fmt"
   "os"
   "strings"
   "testing"
)

func TestLicense(t *testing.T) {
   test := tests[0]
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   private_key, err := os.ReadFile(home + "/widevine/private_key.pem")
   if err != nil {
      t.Fatal(err)
   }
   client_id, err := os.ReadFile(home + "/widevine/client_id.bin")
   if err != nil {
      t.Fatal(err)
   }
   key_id, err := base64.StdEncoding.DecodeString(test.key_id)
   if err != nil {
      t.Fatal(err)
   }
   var module widevine.CDM
   err = module.New(private_key, client_id, widevine.PSSH(key_id, nil))
   if err != nil {
      t.Fatal(err)
   }
   var anon Anonymous
   err = anon.New()
   if err != nil {
      t.Fatal(err)
   }
   detail, err := Path(test.path).Detail()
   if err != nil {
      t.Fatal(err)
   }
   content_id, ok := detail.ContentId()
   if !ok {
      t.Fatal("detail_page.content_id")
   }
   title, err := anon.Entitlement(content_id)
   if err != nil {
      t.Fatal(err)
   }
   play, err := title.Playlist(content_id)
   if err != nil {
      t.Fatal(err)
   }
   key, err := module.Key(play, key_id)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%x\n", key)
}

func TestEntitlement(t *testing.T) {
   test := tests[0]
   var anon Anonymous
   err := anon.New()
   if err != nil {
      t.Fatal(err)
   }
   detail, err := Path(test.path).Detail()
   if err != nil {
      t.Fatal(err)
   }
   content_id, ok := detail.ContentId()
   if !ok {
      t.Fatal("detail_page.content_id")
   }
   title, err := anon.Entitlement(content_id)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(title.Entitlement_Token)
   filter := strings.Split(title.Entitlement_Token, ".")[1]
   data, err := base64.RawStdEncoding.DecodeString(filter)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(string(data))
}

func TestPlaylist(t *testing.T) {
   test := tests[0]
   var anon Anonymous
   err := anon.New()
   if err != nil {
      t.Fatal(err)
   }
   detail, err := Path(test.path).Detail()
   if err != nil {
      t.Fatal(err)
   }
   content_id, ok := detail.ContentId()
   if !ok {
      t.Fatal("detail_page.content_id")
   }
   title, err := anon.Entitlement(content_id)
   if err != nil {
      t.Fatal(err)
   }
   play, err := title.Playlist(content_id)
   if err != nil {
      t.Fatal(err)
   }
   _, filter, _ := strings.Cut(play.ManifestUrl, "=")
   filter = strings.Split(filter, ".")[1]
   fmt.Printf("%q\n", filter)
   data, err := base64.RawStdEncoding.DecodeString(filter)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%q\n", data)
}
