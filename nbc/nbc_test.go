package nbc

import (
   "fmt"
   "os"
   "testing"
)

func TestPlayReady(t *testing.T) {
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile(
      cache + "/nbc/PlayReady", []byte(playReady().String()), os.ModePerm,
   )
   if err != nil {
      t.Fatal(err)
   }
}

func TestVideo(t *testing.T) {
   fmt.Println(video_tests)
}

var video_tests = []struct {
   url     string
   program string
   id      int
   lock    bool
}{
   {
      id:      3494500,
      lock:    true,
      program: "movie",
      url:     "nbc.com/the-matrix/video/the-matrix/3494500",
   },
   {
      id:      9000283422,
      program: "episode",
      url:     "nbc.com/saturday-night-live/video/may-18-jake-gyllenhaal/9000283438",
   },
   {
      id:      9000283435,
      lock:    true,
      program: "episode",
      url:     "nbc.com/saturday-night-live/video/march-30-ramy-youssef/9000283435",
   },
}
