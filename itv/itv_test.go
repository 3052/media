package itv

import (
   "fmt"
   "net/http"
   "net/url"
   "os"
   "os/exec"
   "testing"
)

func TestPlayReady(t *testing.T) {
   user, err := exec.Command(
      "credential", "-h", "api.nordvpn.com", "-k", "user",
   ).Output()
   if err != nil {
      t.Fatal(err)
   }
   password, err := exec.Command(
      "credential", "-h", "api.nordvpn.com",
   ).Output()
   if err != nil {
      t.Fatal(err)
   }
   http.DefaultTransport = &http.Transport{
      Proxy: http.ProxyURL(&url.URL{
         Scheme: "https",
         User:   url.UserPassword(string(user), string(password)),
         Host:   "uk813.nordvpn.com:89",
      }),
   }
   var play Playlist
   err = play.playReady("10_6201_0001.002")
   if err != nil {
      t.Fatal(err)
   }
   hd, ok := play.FullHd()
   if !ok {
      t.Fatal(".FullHd()")
   }
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile(
      cache+"/itv/PlayReady", []byte(hd.KeyServiceUrl), os.ModePerm,
   )
   if err != nil {
      t.Fatal(err)
   }
}

func TestWatch(t *testing.T) {
   fmt.Println(watch_tests)
}

var watch_tests = []struct {
   category string
   id       string
   url      string
}{
   {
      category: "FILM",
      url:      "http://itv.com/watch/solo-a-star-wars-story/10a6201a0001B",
      id:       "10/6201/0001B",
   },
   {
      category: "DRAMA_AND_SOAPS",
      url:      "http://itv.com/watch/grace/2a7610",
      id:       "2/7610",
   },
   {
      category: "DRAMA_AND_SOAPS",
      url:      "http://itv.com/watch/joan/10a3918",
      id:       "10/3918",
   },
}
