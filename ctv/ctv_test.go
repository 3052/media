package ctv

import (
   "41.neocities.org/text"
   "41.neocities.org/widevine"
   "encoding/base64"
   "fmt"
   "os"
   "testing"
   "time"
)

func TestMedia(t *testing.T) {
   for _, test_path := range test_paths {
      resolve, err := Address{test_path}.Resolve()
      if err != nil {
         t.Fatal(err)
      }
      axis, err := resolve.Axis()
      if err != nil {
         t.Fatal(err)
      }
      var media MediaContent
      data, err := media.Marshal(axis)
      if err != nil {
         t.Fatal(err)
      }
      err = media.Unmarshal(data)
      if err != nil {
         t.Fatal(err)
      }
      name := text.Name(&Namer{media})
      fmt.Printf("%q\n", name)
      time.Sleep(time.Second)
   }
}

func TestManifest(t *testing.T) {
   for _, test_path := range test_paths {
      resolve, err := Address{test_path}.Resolve()
      if err != nil {
         t.Fatal(err)
      }
      axis, err := resolve.Axis()
      if err != nil {
         t.Fatal(err)
      }
      var media MediaContent
      data, err := media.Marshal(axis)
      if err != nil {
         t.Fatal(err)
      }
      err = media.Unmarshal(data)
      if err != nil {
         t.Fatal(err)
      }
      manifest, err := axis.Manifest(&media)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(string(manifest))
      time.Sleep(time.Second)
   }
}

func TestLicense(t *testing.T) {
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
   var pssh widevine.Pssh
   pssh.ContentId, err = base64.StdEncoding.DecodeString(content_id)
   if err != nil {
      t.Fatal(err)
   }
   var module widevine.Module
   err = module.New(private_key, client_id, pssh.Marshal())
   if err != nil {
      t.Fatal(err)
   }
   key_id, err := base64.StdEncoding.DecodeString(raw_key_id)
   if err != nil {
      t.Fatal(err)
   }
   key, err := module.Key(Client{}, key_id)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%x\n", key)
}

// ctv.ca/movies/the-girl-with-the-dragon-tattoo-2011
const (
   content_id = "ZmYtZDAxM2NhN2EtMjY0MjY1"
   raw_key_id = "ywlXHuvLP3KHICZX9rn3pg=="
)

var test_paths = []string{
   // ctv.ca/shows/friends/the-one-with-the-chicken-pox-s2e23
   "/shows/friends/the-one-with-the-chicken-pox-s2e23",
   // ctv.ca/movies/the-girl-with-the-dragon-tattoo-2011
   "/movies/the-girl-with-the-dragon-tattoo-2011",
   // ctv.ca/movies/baby-driver
   "/movies/baby-driver",
}
