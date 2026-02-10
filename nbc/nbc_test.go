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
   key_id  string
}{
   {
      id:      3494500,
      lock:    true,
      program: "movie",
      url:     "nbc.com/the-matrix/video/the-matrix/3494500",
   },
   {
      id:      9000283422,
      key_id:  "0552e44842654a4e81b326004be47be0",
      program: "episode",
      url:     "nbc.com/saturday-night-live/video/may-18-jake-gyllenhaal/9000283438",
   },
   {
      id:      9000283435,
      key_id:  "a48d84f23ec74aa1ba8b1d4c863ac02b",
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
