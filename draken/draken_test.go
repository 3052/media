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

func Test(t *testing.T) {
   for _, testVar := range tests {
      var movieVar Movie
      err := movieVar.Fetch(testVar.custom_id)
      if err != nil {
         t.Fatal(err)
      }
      time.Sleep(time.Second)
   }
}
