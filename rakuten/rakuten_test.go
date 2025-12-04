package rakuten

import (
   "net/http"
   "net/url"
   "os"
   "os/exec"
   "testing"
)

var classification_tests = []string{
   "https://rakuten.tv/uk?content_type=tv_shows&tv_show_id=clink",
   "https://rakuten.tv/pt/movies/bound",
   "https://rakuten.tv/dk/movies/a-time-to-kill",
   "https://rakuten.tv/es/movies/una-obra-maestra",
   "https://rakuten.tv/cz?content_type=movies&content_id=transvulcania-the-people-s-run",
   "https://rakuten.tv/fr?content_type=movies&content_id=michael-clayton",
   "https://rakuten.tv/nl?content_type=movies&content_id=made-in-america",
   "https://rakuten.tv/pl?content_type=movies&content_id=ad-astra",
   "https://rakuten.tv/se?content_type=movies&content_id=i-heart-huckabees",
}

var address_tests = []struct {
   format string
   url    string
}{
   {
      format: "/movies/",
      url:    "https://rakuten.tv/nl/movies/made-in-america",
   },
   {
      format: "/player/movies/stream/",
      url:    "https://rakuten.tv/nl/player/movies/stream/made-in-america",
   },
   {
      format: "/tv_shows/",
      url:    "https://rakuten.tv/fr/tv_shows/une-femme-d-honneur",
   },
   {
      format: "?content_id=",
      url:    "https://rakuten.tv/nl?content_type=movies&content_id=made-in-america",
   },
   {
      format: "?tv_show_id=",
      url:    "https://rakuten.tv/uk?content_type=tv_shows&tv_show_id=clink",
   },
}

func TestPlayReady(t *testing.T) {
   const (
      address  = "https://rakuten.tv/cz?content_type=movies&content_id=transvulcania-the-people-s-run"
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
   var movie_var Movie
   err = movie_var.ParseURL(address)
   if err != nil {
      t.Fatal(err)
   }
   stream, err := movie_var.RequestStream(language, Player.PlayReady, Quality.HD)
   if err != nil {
      t.Fatal(err)
   }
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile(
      cache+"/rakuten/PlayReady",
      []byte(stream.StreamInfos[0].LicenseUrl), os.ModePerm,
   )
   if err != nil {
      t.Fatal(err)
   }
}

func TestLog(t *testing.T) {
   t.Log(address_tests, classification_tests)
}
