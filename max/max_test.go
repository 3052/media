package max

import (
   "fmt"
   "os"
   "testing"
   "time"
)

var tests = []struct {
   location   []string
   url        string
   video_type string
}{
   {
      url:        "play.max.com/video/watch/5c762883-279e-40ed-ab84-43fdda9d88a0/560abdc4-ee5e-4f86-807e-38bb9feabe0e",
      video_type: "MOVIE",
      location: []string{
         "Brazil",
         "Chile",
         "Colombia",
         "Denmark",
         "Finland",
         "France",
         "Mexico",
         "Norway",
         "Peru",
         "Sweden",
         "United States",
      },
   },
   {
      url:        "play.max.com/video/watch/857fc45b-5652-42ca-9192-ac1e5e456300/c6258f3b-1c15-4cef-8f1c-1848b22e3f11",
      video_type: "MOVIE",
      location: []string{
         "Chile",
         "Colombia",
         "Indonesia",
         "Malaysia",
         "Mexico",
         "Peru",
         "Philippines",
         "Singapore",
         "Thailand",
      },
   },
   {
      video_type: "EPISODE",
      url:        "play.max.com/video/watch/28ae9450-8192-4277-b661-e76eaad9b2e6/e19442fb-c7ac-4879-8d50-a301f613cb96",
      location: []string{
         "Belgium",
         "Brazil",
         "Bulgaria",
         "Chile",
         "Colombia",
         "Croatia",
         "Czech Republic",
         "Denmark",
         "Finland",
         "France",
         "Hungary",
         "Indonesia",
         "Malaysia",
         "Mexico",
         "Netherlands",
         "Norway",
         "Peru",
         "Philippines",
         "Poland",
         "Portugal",
         "Romania",
         "Singapore",
         "Slovakia",
         "Spain",
         "Sweden",
         "Thailand",
         "United States",
      },
   },
}

func Test(t *testing.T) {
   data, err := os.ReadFile("login.txt")
   if err != nil {
      t.Fatal(err)
   }
   var login1 Login
   err = login1.Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   for _, test := range tests {
      var watch WatchUrl
      err = watch.Set(test.url)
      if err != nil {
         t.Fatal(err)
      }
      play, err := login1.Playback(&watch)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Printf("%+v\n", play.Fallback)
      time.Sleep(time.Second)
   }
}
