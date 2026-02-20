package rakuten

import "testing"

var classification_tests = []string{
   "https://rakuten.tv/pt/movies/bound",
   "https://rakuten.tv/ie/movies/miss-sloane",
   "https://rakuten.tv/dk/movies/a-time-to-kill",
   "https://rakuten.tv/cz?content_type=movies&content_id=transvulcania-the-people-s-run",
   "https://rakuten.tv/es/movies/una-obra-maestra",
   "https://rakuten.tv/fr?content_type=movies&content_id=michael-clayton",
   "https://rakuten.tv/nl?content_type=movies&content_id=made-in-america",
   "https://rakuten.tv/pl?content_type=movies&content_id=ad-astra",
   "https://rakuten.tv/se?content_type=movies&content_id=i-heart-huckabees",
   "https://rakuten.tv/uk?content_type=tv_shows&tv_show_id=clink",
}

func TestLog(t *testing.T) {
   t.Log(address_tests, classification_tests)
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
