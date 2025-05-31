package rakuten

import "testing"

type web_test struct {
   language string
   location string
   url      string
}

var web_tests = []web_test{
   {
      language: "ENG",
      location: "gb",
      url:      "//rakuten.tv/uk?content_type=tv_shows&tv_show_id=hell-s-kitchen-usa",
   },
   {
      language: "ENG",
      location: "at",
      url:      "//rakuten.tv/at?content_type=movies&content_id=ricky-bobby-koenig-der-rennfahrer",
   },
   {
      language: "ENG",
      location: "ch",
      url:      "//rakuten.tv/ch?content_type=movies&content_id=ricky-bobby-koenig-der-rennfahrer",
   },
   {
      language: "SPA",
      location: "cz",
      url:      "//rakuten.tv/cz?content_type=movies&content_id=transvulcania-the-people-s-run",
   },
   {
      language: "ENG",
      location: "de",
      url:      "//rakuten.tv/de?content_type=movies&content_id=ricky-bobby-konig-der-rennfahrer",
   },
   {
      language: "ENG",
      location: "fr",
      url:      "//rakuten.tv/fr?content_type=movies&content_id=infidele",
   },
   {
      language: "ENG",
      location: "ie",
      url:      "//rakuten.tv/ie?content_type=movies&content_id=talladega-nights-the-ballad-of-ricky-bobby",
   },
   {
      language: "ENG",
      location: "nl",
      url:      "//rakuten.tv/nl?content_type=movies&content_id=a-knight-s-tale",
   },
   {
      language: "ENG",
      location: "pl",
      url:      "//rakuten.tv/pl?content_type=movies&content_id=ad-astra",
   },
   {
      language: "ENG",
      location: "se",
      url:      "//rakuten.tv/se?content_type=movies&content_id=i-heart-huckabees",
   },
}

func Test(t *testing.T) {
   for _, test1 := range web_tests {
      fmt.Println(test1)
   }
}
