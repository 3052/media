package rakuten

import "testing"

var web_tests = []web_test{
   {
      address:  "rakuten.tv/at/movies/ricky-bobby-koenig-der-rennfahrer",
      language: "ENG",
      location: "at",
   },
   {
      address:  "rakuten.tv/ch/movies/ricky-bobby-koenig-der-rennfahrer",
      language: "ENG",
      location: "ch",
   },
   {
      address:  "rakuten.tv/cz/movies/transvulcania-the-people-s-run",
      language: "SPA",
      location: "cz",
   },
   {
      address:  "rakuten.tv/de/movies/ricky-bobby-konig-der-rennfahrer",
      language: "ENG",
      location: "de",
   },
   {
      address:  "rakuten.tv/fr/movies/infidele",
      language: "ENG",
      location: "fr",
   },
   {
      address:  "rakuten.tv/ie/movies/talladega-nights-the-ballad-of-ricky-bobby",
      language: "ENG",
      location: "ie",
   },
   {
      address:  "rakuten.tv/nl/movies/a-knight-s-tale",
      language: "ENG",
      location: "nl",
   },
   {
      address:  "rakuten.tv/pl/movies/ad-astra",
      language: "ENG",
      location: "pl",
   },
   {
      address:  "rakuten.tv/se/movies/i-heart-huckabees",
      language: "ENG",
      location: "se",
   },
   {
      address:  "rakuten.tv/uk/player/episodes/stream/hell-s-kitchen-usa-15/hell-s-kitchen-usa-15-1",
      language: "ENG",
      location: "gb",
   },
}

type web_test struct {
   address  string
   language string
   location string
}

func Test(t *testing.T) {
   for _, test1 := range web_tests {
      var web Address
      web.Set(test1.address)
      class, ok := web.ClassificationId()
      if !ok {
         t.Fatal(".ClassificationId()")
      }
      if web.SeasonId != "" {
         data, err := web.Season(class)
         if err != nil {
            t.Fatal(err)
         }
         var season1 Season
         err = season1.Unmarshal(data)
         if err != nil {
            t.Fatal(err)
         }
      } else {
         _, err := web.Movie(class)
         if err != nil {
            t.Fatal(err)
         }
      }
   }
}
