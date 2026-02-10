package peacock

import (
   "154.pages.dev/media/internal"
   "154.pages.dev/widevine"
   "encoding/hex"
   "fmt"
   "os"
   "testing"
)

func TestQuery(t *testing.T) {
   var node QueryNode
   err := node.New(content_id)
   if err != nil {
      t.Fatal(err)
   }
   name, err := internal.Name(node)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%q\n", name)
}

func TestVideo(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   text, err := os.ReadFile(home + "/peacock.json")
   if err != nil {
      t.Fatal(err)
   }
   var sign SignIn
   sign.Unmarshal(text)
   auth, err := sign.Auth()
   if err != nil {
      t.Fatal(err)
   }
   video, err := auth.Video(content_id)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", video)
}

// peacocktv.com/watch/playback/vod/GMO_00000000224510_02_HDSDR
const (
   content_id = "GMO_00000000224510_02_HDSDR"
   raw_key_id = "0016e23473ebe77d93d8d1a72dc690d7"
)

func TestLicense(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(home + "/peacock.json")
   if err != nil {
      t.Fatal(err)
   }
   var sign SignIn
   sign.Unmarshal(data)
   auth, err := sign.Auth()
   if err != nil {
      t.Fatal(err)
   }
   video, err := auth.Video(content_id)
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
   key_id, err := hex.DecodeString(raw_key_id)
   if err != nil {
      t.Fatal(err)
   }
   var module widevine.CDM
   if err := module.New(private_key, client_id, key_id); err != nil {
      t.Fatal(err)
   }
   license, err := module.License(video)
   if err != nil {
      t.Fatal(err)
   }
   key, ok := module.Key(license)
   fmt.Printf("%x %v\n", key, ok)
}

