package rakuten

import "testing"

var web_tests = []web_test{
   {
      address:    "rakuten.tv/at/movies/ricky-bobby-koenig-der-rennfahrer",
      language:   "ENG",
      location:   "at",
      key_id:     "OsBLtLhCGMexX+THBcRobw==",
      content_id: "M2FjMDRiYjRiODQyMThjN2IxNWZlNGM3MDVjNDY4NmYtbWMtMC0xMzktMC0w",
   },
   {
      address:    "rakuten.tv/ch/movies/ricky-bobby-koenig-der-rennfahrer",
      language:   "ENG",
      location:   "ch",
      key_id:     "OsBLtLhCGMexX+THBcRobw==",
      content_id: "M2FjMDRiYjRiODQyMThjN2IxNWZlNGM3MDVjNDY4NmYtbWMtMC0xMzktMC0w",
   },
   {
      address:    "rakuten.tv/cz/movies/transvulcania-the-people-s-run",
      content_id: "MzE4ZjdlY2U2OWFmY2ZlM2U5NmRlMzFiZTZiNzcyNzItbWMtMC0xNjQtMC0w",
      key_id:     "MY9+zmmvz+PpbeMb5rdycg==",
      language:   "SPA",
      location:   "cz",
   },
   {
      address:    "rakuten.tv/de/movies/ricky-bobby-konig-der-rennfahrer",
      language:   "ENG",
      location:   "de",
      key_id:     "OsBLtLhCGMexX+THBcRobw==",
      content_id: "M2FjMDRiYjRiODQyMThjN2IxNWZlNGM3MDVjNDY4NmYtbWMtMC0xMzktMC0w",
   },
   {
      address:    "rakuten.tv/fr/movies/infidele",
      content_id: "MGU1MTgwMDA2Y2Q1MDhlZWMwMGQ1MzVmZWM2YzQyMGQtbWMtMC0xNDEtMC0w",
      key_id:     "DlGAAGzVCO7ADVNf7GxCDQ==",
      language:   "ENG",
      location:   "fr",
   },
   {
      address:    "rakuten.tv/ie/movies/talladega-nights-the-ballad-of-ricky-bobby",
      language:   "ENG",
      location:   "ie",
      key_id:     "r+ROUU1Y1yEFHQKKKSmwkg==",
      content_id: "YWZlNDRlNTE0ZDU4ZDcyMTA1MWQwMjhhMjkyOWIwOTItbWMtMC0xNDMtMC0w",
   },
   {
      address:    "rakuten.tv/nl/movies/a-knight-s-tale",
      content_id: "MGJlNmZmYWRhMzY2NjNhMGExNzMwODYwN2U3Y2ZjYzYtbWMtMC0xMzctMC0w",
      key_id:     "C+b/raNmY6Chcwhgfnz8xg==",
      language:   "ENG",
      location:   "nl",
   },
   {
      address:    "rakuten.tv/pl/movies/ad-astra",
      content_id: "YTk1MjMzMDI1NWFiOWJmZmIxYTAwZTk3ZDA1ZTBhZjItbWMtMC0xMzctMC0w",
      key_id:     "qVIzAlWrm/+xoA6X0F4K8g==",
      language:   "ENG",
      location:   "pl",
   },
   {
      address:    "rakuten.tv/se/movies/i-heart-huckabees",
      content_id: "OWE1MzRhMWYxMmQ2OGUxYTIzNTlmMzg3MTBmZGRiNjUtbWMtMC0xNDctMC0w",
      key_id:     "mlNKHxLWjhojWfOHEP3bZQ==",
      language:   "ENG",
      location:   "se",
   },
   {
      address:    "rakuten.tv/uk/player/episodes/stream/hell-s-kitchen-usa-15/hell-s-kitchen-usa-15-1",
      content_id: "YmI5NGE0YTA0MTdkMjYyY2MzMGMyZjIzODExNmQ2NzktbWMtMC0xMzktMC0w",
      key_id:     "u5SkoEF9JizDDC8jgRbWeQ==",
      language:   "ENG",
      location:   "gb",
   },
}

type web_test struct {
   address    string
   content_id string
   key_id     string
   language   string
   location   string
}

func Test(t *testing.T) {
   for _, test1 := range web_tests {
      var web Address
      web.New(test1.address)
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
