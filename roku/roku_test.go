package roku

import (
   "fmt"
   "testing"
   "time"
)

var tests = map[string]struct {
   content_id string
   id         string
   key_id     string
   url        string
}{
   "episode": {
      content_id: "Kg==",
      id:         "105c41ea75775968b670fbb26978ed76",
      key_id:     "vfpNbNs5cC5baB+QYX+afg==",
      url:        "therokuchannel.roku.com/watch/105c41ea75775968b670fbb26978ed76",
   },
   "movie": {
      content_id: "Kg==",
      id:         "597a64a4a25c5bf6af4a8c7053049a6f",
      key_id:     "KDOa149zRSDaJObgVz05Lg==",
      url:        "therokuchannel.roku.com/watch/597a64a4a25c5bf6af4a8c7053049a6f",
   },
}

func Test(t *testing.T) {
   for _, test := range tests {
      var token1 Token
      data, err := token1.Marshal(nil)
      if err != nil {
         t.Fatal(err)
      }
      err = token1.Unmarshal(data)
      if err != nil {
         t.Fatal(err)
      }
      play, err := token1.Playback(test.id)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(play)
      time.Sleep(time.Second)
   }
}
