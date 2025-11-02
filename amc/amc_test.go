package amc

import (
   "fmt"
   "testing"
)

var tests = []struct {
   content string
   id      int
   url     string
}{
   {
      content: "movie",
      url:     "amcplus.com/movies/nocebo--1061554",
      id:      1061554,
   },
   {
      content: "season",
      url:     "amcplus.com/shows/orphan-black/episodes--1010638",
      id:      1010638,
   },
   {
      content: "series",
      url:     "amcplus.com/shows/orphan-black--1010578",
      id:      1010578,
   },
}

func Test(t *testing.T) {
   fmt.Println(tests)
}
