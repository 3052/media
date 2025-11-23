package ctv

import "testing"

var test_address = []string{
   "https://ctv.ca/movies/heathers",
   "https://ctv.ca/shows/greys-anatomy/we-built-this-city-s22e2",
}

func Test(t *testing.T) {
   for _, address := range test_address {
      t.Log(address)
   }
}
