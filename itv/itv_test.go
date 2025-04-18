package itv

import (
   "os"
   "testing"
   "time"
)

var tests = []struct {
   id  string
   url string
}{
   {
      id:  "18910",
      url: "itv.com/watch/goldeneye/18910",
   },
   {
      id:  "10a5503a0001B",
      url: "itv.com/watch/gone-girl/10a5503a0001B",
   },
   {
      id:  "10_5503_0001B",
      url: "itv.com/watch/gone-girl/10_5503_0001B",
   },
   {
      url:        "itv.com/watch/community/10a3915/10a3915a0002",
      id: "10a3915a0002",
   },
   {
      url:        "itv.com/watch/joan/10a3918/10a3918a0001",
      id: "10a3918a0001",
   },
}

func Test(t *testing.T) {
   for _, test1 := range tests {
      var id legacy_id
      id.Set(test1.id)
      resp, err := id.programme_page()
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
