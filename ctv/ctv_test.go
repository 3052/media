package ctv

import (
   "testing"
   "time"
)

var tests = []struct {
   path string
   url  string
}{
   {
      url: "ctv.ca/movies/x-men-days-of-future-past",
   },
   {
      path: "/shows/friends/the-one-with-the-chicken-pox-s2e23",
      url:  "ctv.ca/shows/friends/the-one-with-the-chicken-pox-s2e23",
   },
}

func Test(t *testing.T) {
   for _, test1 := range tests {
      _, err := Address(test1.path).Resolve()
      if err != nil {
         t.Fatal(err)
      }
      time.Sleep(time.Second)
   }
}
