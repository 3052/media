package disney

import (
   "encoding/xml"
   "fmt"
   "os"
   "testing"
)

func TestPlayback(t *testing.T) {
   cache, err := os.UserCacheDir()
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
   play, err := token.playback()
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
