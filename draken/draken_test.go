package draken

import (
   "testing"
   "time"
)

var tests = []struct {
   content_id string
   custom_id  string
   key_id     string
   url        string
}{
   {
      content_id: "MjNkM2MxYjYtZTA0ZC00ZjMyLWIwYTYtOTgxYzU2MTliNGI0",
      custom_id:  "moon",
      key_id:     "74/ZQoQJukeOkUjy76DE+Q==",
      url:        "drakenfilm.se/film/moon",
   },
   {
      content_id: "MTcxMzkzNTctZWQwYi00YTE2LThiZTYtNjllNDE4YzRiYTQw",
      key_id:     "ToV4wH2nlVZE8QYLmLywDg==",
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
