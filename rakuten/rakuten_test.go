package rakuten

import "testing"

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
      url: "rakuten.tv/uk?content_type=tv_shows&tv_show_id=clink",
   },
}

func Test(t *testing.T) {
   for _, test1 := range web_tests {
      fmt.Println(test1)
   }
}
