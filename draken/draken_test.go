package draken

import (
   "testing"
   "time"
)

var tests = []struct {
   custom_id  string
   url        string
}{
   {
      custom_id:  "moon",
      url:        "drakenfilm.se/film/moon",
   },
   {
      custom_id:  "the-card-counter",
      url:        "drakenfilm.se/film/the-card-counter",
   },
}

func TestDraken(t *testing.T) {
   for _, test := range tests {
      var movie_var Movie
      err := movie_var.Fetch(test.custom_id)
      if err != nil {
         t.Fatal(err)
      }
      time.Sleep(time.Second)
   }
}
