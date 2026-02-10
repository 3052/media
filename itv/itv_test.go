package itv

import (
   "fmt"
   "testing"
   "time"
)

var tests = []struct {
   content_id string
   key_id     string
   legacy_id  LegacyId
   url        string
}{
   {
      content_id: "MTAtMzkxNS0wMDAyLTAwMV8zNA==",
      key_id:     "zCXIAYrkT9+eG6gbjNG1Qw==",
      legacy_id:  LegacyId{"10", "3915", "0002"},
      url:        "itv.com/watch/community/10a3915/10a3915a0002",
   },
   {
      content_id: "MTAtNTUwMy0wMDAxLTAwMV8yMg==",
      key_id: "FUl4yiBqSRC1imOJbh17og==",
      legacy_id:  LegacyId{"10", "5503", "0001"},
      url:        "itv.com/watch/gone-girl/10a5503a0001",
   },
   {
      content_id: "MTAtMzkxOC0wMDAxLTAwMV8zNA==",
      key_id:     "znjzKgOaRBqJMBDGiUDN8g==",
      legacy_id:  LegacyId{"10", "3918", "0001"},
      url:        "itv.com/watch/joan/10a3918/10a3918a0001",
   },
}

func Test(t *testing.T) {
   for _, test := range tests {
      play, err := test.legacy_id.Playlist()
      if err != nil {
         t.Fatal(err)
      }
      fmt.Printf("%+v\n\n", play)
      time.Sleep(time.Second)
   }
}
