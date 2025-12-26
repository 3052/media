package itv

import (
   "net/http"
   "net/url"
   "os"
   "os/exec"
   "testing"
)

func output(name string, arg ...string) (string, error) {
   data, err := exec.Command(name, arg...).Output()
   if err != nil {
      return "", err
   }
   return string(data), nil
}

func TestPlayReady(t *testing.T) {
   user, err := output("credential", "-h=api.nordvpn.com", "-k=user")
   if err != nil {
      t.Fatal(err)
   }
   password, err := output("credential", "-h=api.nordvpn.com")
   if err != nil {
      t.Fatal(err)
   }
   http.DefaultTransport = &http.Transport{
      Proxy: http.ProxyURL(&url.URL{
         Scheme: "https",
         User:   url.UserPassword(user, password),
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

var watch_tests = []struct {
   category string
   watch    []string
}{
   {
      category: "DRAMA_AND_SOAPS",
      watch: []string{
         "https://itv.com/watch/joan/10a3918",
         "https://itv.com/watch/joan/10a3918/10a3918a0001",
      },
   },
   {
      category: "FILM",
      watch:    []string{"https://itv.com/watch/love-actually/27304"},
   },
}

func TestWatch(t *testing.T) {
   t.Log(watch_tests)
}
