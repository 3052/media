package ctv

import "testing"

var tests = []struct {
   path string
   url  string
}{
   {
      url: "https://ctv.ca/movies/the-hurt-locker",
   },
   {
      path: "/shows/friends/the-one-with-the-chicken-pox-s2e23",
      url:  "ctv.ca/shows/friends/the-one-with-the-chicken-pox-s2e23",
   },
}

func Test(t *testing.T) {
   for _, testVar := range tests {
      _, err := GetPath(testVar.url)
      if err != nil {
         t.Fatal(err)
      }
   }
}
