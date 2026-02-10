package max

import (
   "fmt"
   "os"
   "testing"
   "time"
)

var tests = []struct {
   url        string
   video_type string
   key_id     []string
}{
   {
      url:        "play.max.com/video/watch/5c762883-279e-40ed-ab84-43fdda9d88a0/560abdc4-ee5e-4f86-807e-38bb9feabe0e",
      video_type: "MOVIE",
      key_id: []string{
         "AQC1NR9S5CJX8MEgkYbXpg==",
         "AQFz3ZsVFjESfMh2rISgjw==",
         "AQLexzSxi5gJMbgkQogYJQ==",
         "AQW5UW421beBH+jIn3XASw==",
      },
   },
   {
      video_type: "EPISODE",
      url:        "play.max.com/video/watch/28ae9450-8192-4277-b661-e76eaad9b2e6/e19442fb-c7ac-4879-8d50-a301f613cb96",
      key_id:     nil,
   },
}

func TestLicense(t *testing.T) {
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
