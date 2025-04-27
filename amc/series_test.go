package amc

import (
   "fmt"
   "testing"
)

func TestEpisodes(t *testing.T) {
   season, err := season_episodes(season.id)
   if err != nil {
      t.Fatal(err)
   }
   var line bool
   for episode := range season.episodes() {
      if line {
         fmt.Println()
      } else {
         line = true
      }
      fmt.Println(&episode.Properties.Metadata)
   }
}

var show = struct {
   id  int64
   url string
}{
   id:  1010578,
   url: "amcplus.com/shows/orphan-black--1010578",
}

var season = struct {
   id  int64
   url string
}{
   url: "amcplus.com/shows/orphan-black/episodes--1010638",
   id:  1010638,
}

func TestSeasons(t *testing.T) {
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
