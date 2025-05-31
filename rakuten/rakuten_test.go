package rakuten

import (
   "41.neocities.org/net"
   "fmt"
   "net/http"
   "net/url"
   "os"
   "os/exec"
   "strings"
   "testing"
)

var test = struct {
   url      string
   season   string
   episode  string
   language string
}{
   url:      "//rakuten.tv/uk?content_type=tv_shows&tv_show_id=clink",
   season:   "clink-1",
   episode:  "clink-1-1",
   language: "ENG",
}

func TestInfo(t *testing.T) {
   data, err := exec.Command("password", "-i", "nordvpn.com").Output()
   if err != nil {
      t.Fatal(err)
   }
   user, password, _ := strings.Cut(string(data), ":")
   http.DefaultTransport = net.Transport(&url.URL{
      Scheme: "https",
      User:   url.UserPassword(user, password),
      Host:   "uk812.nordvpn.com:89",
   })
   data, err = os.ReadFile("tv_show")
   if err != nil {
      t.Fatal(err)
   }
   var show tv_show
   err = show.Set(string(data))
   if err != nil {
      t.Fatal(err)
   }
   info, err := show.info(test.episode, test.language, fhd)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", info)
}

func TestEpisodes(t *testing.T) {
   data, err := os.ReadFile("tv_show")
   if err != nil {
      t.Fatal(err)
   }
   var show tv_show
   err = show.Set(string(data))
   if err != nil {
      t.Fatal(err)
   }
   episodes, err := show.episodes(test.season)
   if err != nil {
      t.Fatal(err)
   }
   for i, episode1 := range episodes {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(&episode1)
   }
}

func TestSeasons(t *testing.T) {
   var show tv_show
   err := show.Set(test.url)
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile("tv_show", []byte(test.url), os.ModePerm)
   if err != nil {
      t.Fatal(err)
   }
   seasons, err := show.seasons()
   if err != nil {
      t.Fatal(err)
   }
   for _, season1 := range seasons {
      fmt.Println(&season1)
   }
}
type web_test struct {
   language string
   url      string
}

var web_tests = []web_test{
   {
      language: "ENG",
      url:      "//rakuten.tv/at?content_type=movies&content_id=ricky-bobby-koenig-der-rennfahrer",
   },
   {
      language: "ENG",
      url:      "//rakuten.tv/ch?content_type=movies&content_id=ricky-bobby-koenig-der-rennfahrer",
   },
   {
      language: "SPA",
      url:      "//rakuten.tv/cz?content_type=movies&content_id=transvulcania-the-people-s-run",
   },
   {
      language: "ENG",
      url:      "//rakuten.tv/de?content_type=movies&content_id=ricky-bobby-konig-der-rennfahrer",
   },
   {
      language: "ENG",
      url:      "//rakuten.tv/fr?content_type=movies&content_id=infidele",
   },
   {
      language: "ENG",
      url:      "//rakuten.tv/ie?content_type=movies&content_id=talladega-nights-the-ballad-of-ricky-bobby",
   },
   {
      language: "ENG",
      url:      "//rakuten.tv/nl?content_type=movies&content_id=a-knight-s-tale",
   },
   {
      language: "ENG",
      url:      "//rakuten.tv/pl?content_type=movies&content_id=ad-astra",
   },
   {
      language: "ENG",
      url:      "//rakuten.tv/se?content_type=movies&content_id=i-heart-huckabees",
   },
   {
      url: "//rakuten.tv/uk?content_type=tv_shows&tv_show_id=clink",
   },
}

func Test(t *testing.T) {
   for _, test1 := range web_tests {
      fmt.Println(test1)
   }
}
