package nbc

import (
   "fmt"
   "testing"
   "time"
)

var tests = []struct {
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

func Test(t *testing.T) {
   for _, test1 := range tests {
      var meta Metadata
      err := meta.New(test1.id)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(meta)
      time.Sleep(time.Second)
   }
}
