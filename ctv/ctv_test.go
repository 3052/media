package ctv

import "testing"

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
   for _, testVar := range tests {
      _, err := Resolve(Path(testVar.url))
      if err != nil {
         t.Fatal(err)
      }
   }
}
