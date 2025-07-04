package itv

import (
   "fmt"
   "testing"
)

func TestPlayReady(t *testing.T) {
   var value Title
   value.LatestAvailableVersion.PlaylistUrl = "https://magni.itv.com/playlist/itvonline/ITV/10_5503_0001.001"
   data, err := value.Playlist()
   if err != nil {
      t.Fatal(err)
   }
   var play Playlist
   err = play.Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   hd, ok := play.FullHd()
   if !ok {
      t.Fatal(".FullHd()")
   }
   fmt.Println(hd.KeyServiceUrl)
}

func TestWatch(t *testing.T) {
   fmt.Println(watch_tests)
}

var watch_tests = []struct {
   id  string
   url string
}{
   {
      id:  "10/5503/0001B",
      url: "itv.com/watch/gone-girl/10a5503a0001B",
   },
   {
      id:  "2/7610",
      url: "itv.com/watch/grace/2a7610",
   },
   {
      id:  "10/3918",
      url: "itv.com/watch/joan/10a3918",
   },
}
