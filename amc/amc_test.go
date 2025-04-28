package amc

import (
   "fmt"
   "testing"
)

/*
flag.StringVar(&series, "series", "", "series ID")
flag.StringVar(&season, "s", "", "season ID")
flag.StringVar(&episode, "e", "", "episode or movie ID")
flag.StringVar(&dash, "d", "", "DASH ID")
*/

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

type test struct {
   id int64
   url string
}

var show = test{
   id:  1010578,
   url: "amcplus.com/shows/orphan-black--1010578",
}

var season = test{
   id:  1010638,
   url: "amcplus.com/shows/orphan-black/episodes--1010638",
}

var movie = test{
   id: 1061554,
   url: "amcplus.com/movies/nocebo--1061554",
}
