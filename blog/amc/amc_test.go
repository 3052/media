package amc

import (
   "fmt"
   "testing"
)

var show = struct{
   id int64
   url string
}{
   id: 1010578,
   url: "amcplus.com/shows/orphan-black--1010578",
}

func Test(t *testing.T) {
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
      fmt.Println(season)
   }
}
