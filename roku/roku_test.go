package roku

import (
   "fmt"
   "testing"
   "time"
)

func Test(t *testing.T) {
   for _, test1 := range tests {
      data, err := (*Code).AccountToken(nil)
      if err != nil {
         t.Fatal(err)
      }
      var token1 AccountToken
      err = token1.Unmarshal(data)
      if err != nil {
         t.Fatal(err)
      }
      play, err := token1.Playback(test1.id)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(play)
      time.Sleep(time.Second)
   }
}

var tests = []struct {
   id    string
   type1 string
   url   string
}{
   {
      id:    "597a64a4a25c5bf6af4a8c7053049a6f",
      type1: "movie",
      url:   "therokuchannel.roku.com/watch/597a64a4a25c5bf6af4a8c7053049a6f",
   },
   {
      id:    "105c41ea75775968b670fbb26978ed76",
      type1: "episode",
      url:   "therokuchannel.roku.com/watch/105c41ea75775968b670fbb26978ed76",
   },
}
