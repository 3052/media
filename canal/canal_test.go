package canal

import (
   "fmt"
   "testing"
)

var tests = []struct {
   id  string
   url string
}{
   {
      id:  "1EBvrU5Q2IFTIWSC2_4cAlD98U0OR0ejZm_dgGJi",
      url: "https://canalplus.cz/stream/film/argylle-tajny-agent",
   },
   {
      id:  "XT0kyelnPAOl3f-Bx7etkj_yX3nDHom_ymdCRK5A",
      url: "https://canalplus.cz/stream/series/fbi",
   },
   {
      id:  "cnygdzw_ntkhIekB6ruh9M2U-k6UQFjQ__DYJALw",
      url: "https://canalplus.cz/stream/series/silo",
   },
}

func TestAssets(t *testing.T) {
   for _, test1 := range tests {
      fmt.Println(test1.url)
   }
}
