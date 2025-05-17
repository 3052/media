package pluto

import (
   "fmt"
   "testing"
)

var tests = []struct {
   id  string
   url string
}{
   {
      id:  "623a01faef11000014cf41f7",
      url: "pluto.tv/on-demand/movies/623a01faef11000014cf41f7",
   },
   {
      id:  "66d0bb64a1c89200137fb0e6",
      url: "pluto.tv/on-demand/series/66d0bb64a1c89200137fb0e6",
   },
}

func Test(t *testing.T) {
   for _, test1 := range tests {
      var web Address
      err := web.Set(test1.url)
      if err != nil {
         t.Fatal(err)
      }
      video, err := web.Vod()
      if err != nil {
         t.Fatal(err)
      }
      fmt.Printf("%+v\n", video)
   }
}
