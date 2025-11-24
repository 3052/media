package rakuten

import (
   "net/http"
   "net/url"
   "os"
   "os/exec"
   "testing"
)

var classification_tests = []string{
   "https://rakuten.tv/pt/movies/bound",
   "https://rakuten.tv/dk/movies/a-time-to-kill",
   "https://rakuten.tv/es/movies/una-obra-maestra",
   "https://rakuten.tv/cz?content_type=movies&content_id=transvulcania-the-people-s-run",
   "https://rakuten.tv/fr?content_type=movies&content_id=michael-clayton",
   "https://rakuten.tv/nl?content_type=movies&content_id=made-in-america",
   "https://rakuten.tv/pl?content_type=movies&content_id=ad-astra",
   "https://rakuten.tv/se?content_type=movies&content_id=i-heart-huckabees",
   "https://rakuten.tv/uk?content_type=tv_shows&tv_show_id=clink",
}

var address_tests = []struct {
   format string
   url    string
}{
   {
      format: "/movies/",
      url:    "http://rakuten.tv/nl/movies/made-in-america",
   },
   {
      format: "/player/movies/stream/",
      url:    "http://rakuten.tv/nl/player/movies/stream/made-in-america",
   },
   {
      format: "/tv_shows/",
      url:    "http://rakuten.tv/fr/tv_shows/une-femme-d-honneur",
   },
   {
      format: "?content_id=",
      url:    "http://rakuten.tv/nl?content_type=movies&content_id=made-in-america",
   },
   {
      format: "?tv_show_id=",
      url:    "http://rakuten.tv/uk?content_type=tv_shows&tv_show_id=clink",
   },
}

func TestPlayReady(t *testing.T) {
   const (
      address  = "http://rakuten.tv/cz?content_type=movies&content_id=transvulcania-the-people-s-run"
      language = "SPA"
   )
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
         Host:   "cz103.nordvpn.com:89",
      }),
   }
   var mediaVar Media
   err = mediaVar.Parse(address)
   if err != nil {
      t.Fatal(err)
   }
   info, err := mediaVar.Pr(mediaVar.ContentId, language, Hd)
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

func TestAddress(t *testing.T) {
   t.Log(classification_tests)
   for _, test := range address_tests {
      var mediaVar Media
      err := mediaVar.Parse(test.url)
      if err != nil {
         t.Fatal(err)
      }
      t.Logf("%+v", mediaVar)
   }
}
