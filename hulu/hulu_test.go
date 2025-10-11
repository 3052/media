package hulu

import (
   "fmt"
   "os"
   "testing"
)

func TestPlayReady(t *testing.T) {
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(cache + "/hulu/Authenticate")
   if err != nil {
      t.Fatal(err)
   }
   var auth Authenticate
   err = auth.Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   err = auth.Refresh()
   if err != nil {
      t.Fatal(err)
   }
   test := tests[0]
   data, err = auth.Playlist(&DeepLink{EabId: test.id})
   if err != nil {
      t.Fatal(err)
   }
   var play Playlist
   err = play.Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile(
      cache + "/hulu/PlayReady", []byte(play.DashPrServer), os.ModePerm,
   )
   if err != nil {
      t.Fatal(err)
   }
}

func TestWatch(t *testing.T) {
   fmt.Println(tests)
}
var tests = []struct {
   content string
   id      string
   url     string
}{
   {
      content: "film",
      id:      "EAB::f70dfd4d-dbfb-46b8-abb3-136c841bba11::61556664::101167038",
      url:     "hulu.com/watch/f70dfd4d-dbfb-46b8-abb3-136c841bba11",
   },
   {
      content: "episode",
      url:     "hulu.com/watch/023c49bf-6a99-4c67-851c-4c9e7609cc1d",
   },
}
