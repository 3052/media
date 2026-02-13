package peacock

import (
   "encoding/hex"
   "fmt"
   "os"
   "testing"
)

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
