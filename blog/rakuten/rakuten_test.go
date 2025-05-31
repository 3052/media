package rakuten

import (
   "fmt"
   "os"
   "testing"
)

var test = struct {
   season string
   url    string
}{
   season: "clink-1",
   url:    "//rakuten.tv/uk?content_type=tv_shows&tv_show_id=clink",
}

func TestSeason(t *testing.T) {
   var show tv_show
   err := show.Set(test.url)
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile("tv_show", []byte(test.url), os.ModePerm)
   if err != nil {
      t.Fatal(err)
   }
   seasons, err := show.seasons()
   if err != nil {
      t.Fatal(err)
   }
   for _, season1 := range seasons {
      fmt.Println(&season1)
   }
}

func TestEpisode(t *testing.T) {
   data, err := os.ReadFile("tv_show")
   if err != nil {
      t.Fatal(err)
   }
   var show tv_show
   err = show.Set(string(data))
   if err != nil {
      t.Fatal(err)
   }
   episodes, err := show.episodes(test.season)
   if err != nil {
      t.Fatal(err)
   }
   for i, episode1 := range episodes {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(&episode1)
   }
}
