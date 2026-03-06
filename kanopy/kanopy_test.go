package kanopy

import (
   "encoding/xml"
   "fmt"
   "os"
   "testing"
)

func TestAlias(t *testing.T) {
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(cache + "/rosso/kanopy.xml")
   if err != nil {
      t.Fatal(err)
   }
   var state struct {
      Login Login
   }
   err = xml.Unmarshal(data, &state)
   if err != nil {
      t.Fatal(err)
   }
   // https://kanopy.com/video/genius-party
   result, err := state.Login.Video("genius-party")
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", result)
}

var tests = []struct {
   genre string
   url   string
}{
   {
      genre: "Movies",
      url:   "https://kanopy.com/video/13808102",
   },
   {
      genre: "TV Series",
      url:   "https://kanopy.com/video/14098194",
   },
}

func TestLog(t *testing.T) {
   t.Log(tests)
}
