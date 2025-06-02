package rakuten

import (
   "fmt"
   "testing"
)

func Test(t *testing.T) {
   test := web_tests[0]
   var web Address
   err := web.Set(test.url)
   if err != nil {
      t.Fatal(err)
   }
   info, err := web.Info(web.ContentId, test.language, Pr, Hd)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(info.LicenseUrl)
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
