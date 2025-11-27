package roku

import (
   "fmt"
   "testing"
   "time"
)

var tests = []struct {
   id         string
   type_value string
   url        string
}{
   {
      id:         "597a64a4a25c5bf6af4a8c7053049a6f",
      type_value: "movie",
      url:        "https://therokuchannel.roku.com/watch/597a64a4a25c5bf6af4a8c7053049a6f",
   },
   {
      id:         "105c41ea75775968b670fbb26978ed76",
      type_value: "episode",
      url:        "https://therokuchannel.roku.com/watch/105c41ea75775968b670fbb26978ed76",
   },
}

func TestRoku(t *testing.T) {
   for _, test := range tests {
      data, err := (*Code).AccountToken(nil)
      if err != nil {
         t.Fatal(err)
      }
      var token_var AccountToken
      err = token_var.Unmarshal(data)
      if err != nil {
         t.Fatal(err)
      }
      play, err := token_var.Playback(test.id)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(play)
      time.Sleep(time.Second)
   }
}
