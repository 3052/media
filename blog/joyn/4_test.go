package joyn

import (
   "154.pages.dev/widevine"
   "encoding/base64"
   "fmt"
   "os"
   "testing"
)

// joyn.de/filme/barry-seal-only-in-america
const (
   raw_content_id = "YV9wNHN2bjRhMjhmcQ=="
   raw_key_id = "e+os9wvbQLpkvIFRuG3exA=="
   raw_pssh = "AAAAQ3Bzc2gAAAAA7e+LqXnWSs6jyCfc1R0h7QAAACMIARIQe+os9wvbQLpkvIFRuG3exCINYV9wNHN2bjRhMjhmcQ=="
)

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
   pssh, err := base64.StdEncoding.DecodeString(raw_pssh)
   if err != nil {
      t.Fatal(err)
   }
   var module widevine.CDM
   err = module.New(private_key, client_id, pssh)
   if err != nil {
      t.Fatal(err)
   }
   var anon anonymous
   err = anon.New()
   if err != nil {
      t.Fatal(err)
   }
   var movie movie_detail
   err = movie.New(barry_seal)
   if err != nil {
      t.Fatal(err)
   }
   title, err := anon.entitlement(movie)
   if err != nil {
      t.Fatal(err)
   }
   play, err := title.playlist(movie)
   if err != nil {
      t.Fatal(err)
   }
   key_id, err := base64.StdEncoding.DecodeString(raw_key_id)
   if err != nil {
      t.Fatal(err)
   }
   key, err := module.Key(play, key_id)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%x\n", key)
}

func TestPlaylist(t *testing.T) {
   var anon anonymous
   err := anon.New()
   if err != nil {
      t.Fatal(err)
   }
   var movie movie_detail
   err = movie.New(barry_seal)
   if err != nil {
      t.Fatal(err)
   }
   title, err := anon.entitlement(movie)
   if err != nil {
      t.Fatal(err)
   }
   play, err := title.playlist(movie)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", play)
}
