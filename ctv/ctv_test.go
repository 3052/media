package ctv

import (
   "testing"
   "time"
)

var tests = []struct {
   content_id string
   key_id     string
   path       string
   url        string
}{
   {
      url:  "ctv.ca/shows/friends/the-one-with-the-chicken-pox-s2e23",
      path: "/shows/friends/the-one-with-the-chicken-pox-s2e23",
   },
   {
      url:        "ctv.ca/movies/the-transporter",
      path:       "/movies/the-transporter",
   },
}

func Test(t *testing.T) {
   for _, test1 := range tests {
      _, err := Address{test1.path}.Resolve()
      if err != nil {
         t.Fatal(err)
      }
      time.Sleep(time.Second)
   }
}
