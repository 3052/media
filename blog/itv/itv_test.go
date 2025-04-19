package itv

import (
   "os"
   "testing"
   "time"
)

var tests = []struct{
   url string
   id string
}{
   {
      url: "itv.com/watch/goldeneye/18910",
      id: "18910",
   },
   {
      url: "itv.com/watch/gone-girl/10a5503a0001B",
      id: "10/5503/0001B",
   },
   {
      url: "itv.com/watch/grace/2a7610",
      id: "2/7610",
   },
   {
      url: "itv.com/watch/joan/10a3918",
      id: "10/3918",
   },
}

func Test(t *testing.T) {
   for _, test1 := range tests {
      resp, err := discovery(test1.id)
      if err != nil {
         t.Fatal(err)
      }
      err = resp.Write(os.Stdout)
      if err != nil {
         t.Fatal(err)
      }
      time.Sleep(time.Second)
   }
}
