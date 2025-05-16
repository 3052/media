package ctv

import (
   "testing"
   "time"
)

var tests = []struct {
   path       string
   url        string
}{
   {
      path: "/shows/friends/the-one-with-the-chicken-pox-s2e23",
      url:  "ctv.ca/shows/friends/the-one-with-the-chicken-pox-s2e23",
   },
   {
      path: "/movies/the-transporter",
      url:  "ctv.ca/movies/the-transporter",
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
