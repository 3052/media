package amc

import (
   "fmt"
   "testing"
)

func Test(t *testing.T) {
   fmt.Println(series, season, movie)
}

type test struct {
   id  int64
   url string
}

var series = test{
   id:  1010578,
   url: "amcplus.com/shows/orphan-black--1010578",
}

var season = test{
   id:  1010638,
   url: "amcplus.com/shows/orphan-black/episodes--1010638",
}

var movie = test{
   id:  1061554,
   url: "amcplus.com/movies/nocebo--1061554",
}
