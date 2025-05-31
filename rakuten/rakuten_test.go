package rakuten

import "testing"

type web_test struct {
   language string
   location string
   url  []string
}

var web_tests = []web_test{
   {
      language: "ENG",
      location: "gb",
      url: []string{
         "rakuten.tv/uk/tv_shows/hell-s-kitchen-usa",
         "rakuten.tv/uk?content_type=tv_shows&tv_show_id=hell-s-kitchen-usa",
      },
   },
   {
      language: "ENG",
      location: "at",
      url: []string{
         "rakuten.tv/at/movies/ricky-bobby-koenig-der-rennfahrer",
         "rakuten.tv/at?content_type=movies&content_id=ricky-bobby-koenig-der-rennfahrer",
      },
   },
   {
      language: "ENG",
      location: "ch",
      url: []string{
         "rakuten.tv/ch/movies/ricky-bobby-koenig-der-rennfahrer",
         "rakuten.tv/ch?content_type=movies&content_id=ricky-bobby-koenig-der-rennfahrer",
      },
   },
   {
      language: "SPA",
      location: "cz",
      url: []string{
         "rakuten.tv/cz/movies/transvulcania-the-people-s-run",
         "rakuten.tv/cz?content_type=movies&content_id=transvulcania-the-people-s-run",
      },
   },
   {
      language: "ENG",
      location: "de",
      url: []string{
         "rakuten.tv/de/movies/ricky-bobby-konig-der-rennfahrer",
         "rakuten.tv/de?content_type=movies&content_id=ricky-bobby-konig-der-rennfahrer",
      },
   },
   {
      language: "ENG",
      location: "fr",
      url: []string{
         "rakuten.tv/fr/movies/infidele",
         "rakuten.tv/fr?content_type=movies&content_id=infidele",
      },
   },
   {
      language: "ENG",
      location: "ie",
      url: []string{
         "rakuten.tv/ie/movies/talladega-nights-the-ballad-of-ricky-bobby",
         "rakuten.tv/ie?content_type=movies&content_id=talladega-nights-the-ballad-of-ricky-bobby",
      },
   },
   {
      language: "ENG",
      location: "nl",
      url: []string{
         "rakuten.tv/nl/movies/a-knight-s-tale",
         "rakuten.tv/nl?content_type=movies&content_id=a-knight-s-tale",
      },
   },
   {
      language: "ENG",
      location: "pl",
      url: []string{
         "rakuten.tv/pl/movies/ad-astra",
         "rakuten.tv/pl?content_type=movies&content_id=ad-astra",
      },
   },
   {
      language: "ENG",
      location: "se",
      url: []string{
         "rakuten.tv/se/movies/i-heart-huckabees",
         "rakuten.tv/se?content_type=movies&content_id=i-heart-huckabees",
      },
   },
}

func Test(t *testing.T) {
   for _, test1 := range web_tests {
      fmt.Println(test1)
   }
}
