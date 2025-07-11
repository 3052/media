package roku

import (
   "fmt"
   "testing"
   "time"
)

var tests = []struct{
   id    string
   typeVar string
   url   string
}{
   {
      id:    "597a64a4a25c5bf6af4a8c7053049a6f",
      typeVar: "movie",
      url:   "therokuchannel.roku.com/watch/597a64a4a25c5bf6af4a8c7053049a6f",
   },
   {
      id:    "105c41ea75775968b670fbb26978ed76",
      typeVar: "episode",
      url:   "therokuchannel.roku.com/watch/105c41ea75775968b670fbb26978ed76",
   },
}

func Test(t *testing.T) {
   for _, testVar := range tests {
      data, err := (*Code).AccountToken(nil)
      if err != nil {
         t.Fatal(err)
      }
      var tokenVar AccountToken
      err = tokenVar.Unmarshal(data)
      if err != nil {
         t.Fatal(err)
      }
      play, err := tokenVar.Playback(testVar.id)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(play)
      time.Sleep(time.Second)
   }
}
