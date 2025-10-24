package rakuten

import (
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
      url:      "http://rakuten.tv/cz?content_type=movies&content_id=transvulcania-the-people-s-run",
   },
   {
      language: "ENG",
      url: "http://rakuten.tv/dk/movies/a-time-to-kill",
   },
   {
      language: "ENG",
      url:      "http://rakuten.tv/fr?content_type=movies&content_id=infidele",
   },
   {
      language: "ENG",
      url:      "http://rakuten.tv/pl?content_type=movies&content_id=ad-astra",
   },
   {
      language: "ENG",
      url:      "http://rakuten.tv/se?content_type=movies&content_id=i-heart-huckabees",
   },
   {
      language: "ENG",
      url:      "http://rakuten.tv/uk?content_type=tv_shows&tv_show_id=clink",
   },
}

var address_tests = []string{
   "https://www.rakuten.tv/fr/movies/michael-clayton",
   "https://www.rakuten.tv/fr/tv_shows/une-femme-d-honneur",
   "https://www.rakuten.tv/fr?content_type=movies&content_id=michael-clayton",
   "https://www.rakuten.tv/fr?content_type=tv_shows&tv_show_id=une-femme-d-honneur&content_id=une-femme-d-honneur-1",
}

func TestAddress(t *testing.T) {
   for _, test := range address_tests {
      var info Address
      err := info.Parse(test)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Printf("%+v\n", info)
   }
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
   err = web.Parse(test.url)
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
