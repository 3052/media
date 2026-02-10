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
      url:        "ctv.ca/movies/fools-rush-in-57470",
      key_id:     "A98dtspZsb9/z++3IHp0Dw==",
      content_id: "ZmYtOGYyNjEzYWUtNTIxNTAx",
      path:       "/movies/fools-rush-in-57470",
   },
   {
      url:  "ctv.ca/shows/friends/the-one-with-the-chicken-pox-s2e23",
      path: "/shows/friends/the-one-with-the-chicken-pox-s2e23",
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
