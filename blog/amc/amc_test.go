package amc

import (
   "fmt"
   "testing"
)

func TestSeason(t *testing.T) {
   series, err := series_detail(show.id)
   if err != nil {
      t.Fatal(err)
   }
   for child1 := range series.seasons() {
      season, err := child1.season()
      if err != nil {
         t.Fatal(err)
      }
      var line bool
      for child2 := range season.episodes() {
         if line {
            fmt.Println()
         } else {
            line = true
         }
         fmt.Println(&child2.Properties.Metadata)
      }
      break
   }
}

var show = struct {
   id  int64
   url string
}{
   id:  1010578,
   url: "amcplus.com/shows/orphan-black--1010578",
}

func TestSeries(t *testing.T) {
   series, err := series_detail(show.id)
   if err != nil {
      t.Fatal(err)
   }
   var line bool
   for season := range series.seasons() {
      if line {
         fmt.Println()
      } else {
         line = true
      }
      fmt.Println(&season.Properties.Metadata)
   }
}
