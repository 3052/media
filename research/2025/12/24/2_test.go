package disney

import (
   "encoding/xml"
   "fmt"
   "os"
   "testing"
)

var test = struct {
   entity string
   url    string
}{
   entity: "7df81cf5-6be5-4e05-9ff6-da33baf0b94d",
   url:    "https://disneyplus.com/play/7df81cf5-6be5-4e05-9ff6-da33baf0b94d",
}

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
   explore, err := token.explore(test.entity)
   if err != nil {
      t.Fatal(err)
   }
   resource_id, ok := explore.restart()
   if !ok {
      t.Fatal(".restart()")
   }
   play, err := token.playback(resource_id)
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
