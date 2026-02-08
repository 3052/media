package draken

import "testing"

var tests = []string{
   "https://drakenfilm.se/film/moon",
   "https://drakenfilm.se/film/the-card-counter",
}

func Test(t *testing.T) {
   t.Log(tests)
}
