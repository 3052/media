package rakuten

import (
   "41.neocities.org/widevine"
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

func TestWidevine(t *testing.T) {
   test := web_tests[0]
   var web Address
   err := web.Set(test.url)
   if err != nil {
      t.Fatal(err)
   }
   info, err := web.Wvm(web.ContentId, test.language, Hd)
   if err != nil {
      t.Fatal(err)
   }
   var pssh widevine.Pssh
   key_id, err := hex.DecodeString("318f7ece69afcfe3e96de31be6b77272")
   if err != nil {
      t.Fatal(err)
   }
   // need both
   pssh.KeyIds = [][]byte{key_id}
   pssh.ContentId = []byte("318f7ece69afcfe3e96de31be6b77272-mc-0-164-0-0")
   var module widevine.Cdm
   private_key, err := os.ReadFile("C:/Users/Steven/media/private_key.pem")
   if err != nil {
      t.Fatal(err)
   }
   client_id, err := os.ReadFile("C:/Users/Steven/media/client_id.bin")
   if err != nil {
      t.Fatal(err)
   }
   err = module.New(private_key, client_id, pssh.Marshal())
   if err != nil {
      t.Fatal(err)
   }
   data, err := module.RequestBody()
   if err != nil {
      t.Fatal(err)
   }
   resp, err := http.Post(
      info.LicenseUrl, "application/x-protobuf", bytes.NewReader(data),
   )
   if err != nil {
      t.Fatal(err)
   }
   defer resp.Body.Close()
   fmt.Printf("%+v\n", resp)
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
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile(
      home+"/media/rakuten/PlayReady",
      []byte(info.LicenseUrl), os.ModePerm,
   )
   if err != nil {
      t.Fatal(err)
   }
}

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
      url:      "//rakuten.tv/uk?content_type=tv_shows&tv_show_id=clink",
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
}
