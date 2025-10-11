package rakuten

import (
   "bytes"
   "encoding/hex"
   "fmt"
   "net/http"
   "net/url"
   "os"
   "os/exec"
   "strings"
   "testing"
)

var web_tests = []struct {
   language string
   url      string
}{
   {
      language: "SPA",
      url:      "//rakuten.tv/cz?content_type=movies&content_id=transvulcania-the-people-s-run",
   },
   {
      language: "ENG",
      url:      "//rakuten.tv/fr?content_type=movies&content_id=infidele",
   },
   {
      language: "ENG",
      url:      "//rakuten.tv/pl?content_type=movies&content_id=ad-astra",
   },
   {
      language: "ENG",
      url:      "//rakuten.tv/uk?content_type=tv_shows&tv_show_id=clink",
   },
   {
      language: "ENG",
      url:      "//rakuten.tv/se?content_type=movies&content_id=i-heart-huckabees",
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
         Host:   "cz103.nordvpn.com:89",
      }),
   }
   test := web_tests[0]
   var web Address
   err = web.Set(test.url)
   if err != nil {
      t.Fatal(err)
   }
   info, err := web.Pr(web.ContentId, test.language, Hd)
   if err != nil {
      t.Fatal(err)
   }
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile(
      cache+"/rakuten/PlayReady", []byte(info.LicenseUrl), os.ModePerm,
   )
   if err != nil {
      t.Fatal(err)
   }
}
