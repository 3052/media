package disney

import "testing"

func Test(t *testing.T) {
   t.Log(tests)
}

var tests = []struct {
   entity   string
   location string
   url      string
}{
   {
      entity:   "movie",
      location: "US",
      url: "https://disneyplus.com/browse/entity-7df81cf5-6be5-4e05-9ff6-da33baf0b94d",
   },
   {
      entity:   "episode",
      location: "US",
      url: "https://disneyplus.com/browse/entity-21e70fbf-6a51-41b3-88e9-f111830b046c",
   },
   {
      entity:   "movie",
      location: "KR", // MUST DO KR LOGIN FIRST
      url: "https://disneyplus.com/browse/entity-d0d0796c-a144-42fa-a730-4cbd1014ef1f",
   },
}
