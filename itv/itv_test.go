package itv

import (
   "net/http"
   "net/url"
   "os"
   "os/exec"
   "testing"
)

// props.pageProps.seriesList[0].titles[0].playlistUrl
var watch_tests = []struct {
   category string
   watch      string
}{
   {
      category: "ENTERTAINMENT",
      watch: "https://itv.com/watch/im-a-celebrity-get-me-out-of-here/L2649/L2649a0039",
   },
   {
      category: "FILM",
      watch: "https://itv.com/watch/love-actually/27304",
   },
   {
      category: "DRAMA_AND_SOAPS",
      watch: "https://itv.com/watch/joan/10a3918",
   },
}

func TestWatch(t *testing.T) {
   t.Log(watch_tests)
}

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
