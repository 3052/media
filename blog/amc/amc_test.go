package amc

import (
   "fmt"
   "os"
   "testing"
)

var show = struct {
   id  int64
   url string
}{
   id:  1010578,
   url: "amcplus.com/shows/orphan-black--1010578",
}

func TestSeason(t *testing.T) {
   series, err := series_detail(show.id)
   if err != nil {
      t.Fatal(err)
   }
   for season := range series.seasons() {
      resp, err := season.Callback.do()
      if err != nil {
         t.Fatal(err)
      }
      defer resp.Body.Close()
      file, err := os.Create("amc.json")
      if err != nil {
         t.Fatal(err)
      }
      defer file.Close()
      _, err = file.ReadFrom(resp.Body)
      if err != nil {
         t.Fatal(err)
      }
      break
   }
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
      fmt.Println(season)
   }
}
