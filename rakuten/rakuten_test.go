package rakuten

import (
   "net/http"
   "net/url"
   "os"
   "os/exec"
   "strings"
   "testing"
)

func Test(t *testing.T) {
   data, err := exec.Command("password", "-i", "nordvpn.com").Output()
   if err != nil {
      t.Fatal(err)
   }
   username, password, _ := strings.Cut(string(data), ":")
   http.DefaultTransport = &http.Transport{
      Proxy: http.ProxyURL(&url.URL{
         Scheme: "https",
         User: url.UserPassword(username, password),
         Host: "cz103.nordvpn.com:89",
      }),
   }
   test := web_tests[0]
   var web Address
   err = web.Set(test.url)
   if err != nil {
      t.Fatal(err)
   }
   info, err := web.Info(web.ContentId, test.language, Pr, Hd)
   if err != nil {
      t.Fatal(err)
   }
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile(
      home + "/media/rakuten/PlayReady",
      []byte(info.LicenseUrl), os.ModePerm,
   )
   if err != nil {
      t.Fatal(err)
   }
}

var web_tests = []web_test{
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
      url:      "//rakuten.tv/se?content_type=movies&content_id=i-heart-huckabees",
   },
   {
      language: "ENG",
      url:      "//rakuten.tv/uk?content_type=tv_shows&tv_show_id=clink",
   },
}

type web_test struct {
   language string
   url      string
}
