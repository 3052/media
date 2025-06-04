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
   for _, test1 := range tests {
      var movie1 Movie
      err := movie1.New(test1.custom_id)
      if err != nil {
         t.Fatal(err)
      }
      time.Sleep(time.Second)
   }
}
