package itv

import (
   "fmt"
   "net/http"
   "net/url"
   "os"
   "os/exec"
   "strings"
   "testing"
)

var watch_tests = []struct {
   category string
   id       string
   url      string
}{
   {
      category: "DRAMA_AND_SOAPS",
      url:      "https://itv.com/watch/grace/2a7610",
      id:       "2/7610",
   },
   {
      category: "DRAMA_AND_SOAPS",
      url:      "https://itv.com/watch/joan/10a3918",
      id:       "10/3918",
   },
   {
      category: "FILM",
      url:      "https://itv.com/watch/solo-a-star-wars-story/10a6201a0001B",
      id:       "10/6201/0001B",
   },
}

func TestPlayReady(t *testing.T) {
   data, err := exec.Command("password", "-i", "nordvpn.com").Output()
   if err != nil {
      t.Fatal(err)
   }
   username, password, _ := strings.Cut(string(data), ":")
   http.DefaultTransport = &http.Transport{
      Proxy: http.ProxyURL(&url.URL{
         Scheme: "https",
         User:   url.UserPassword(username, password),
         Host:   "uk871.nordvpn.com:89",
      }),
   }
   var play Playlist
   err = play.playReady("10_5503_0001.001")
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
