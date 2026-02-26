package disney

import "testing"

func Test(t *testing.T) {
   t.Log(tests)
}

var tests = []struct {
   entity string
   format string
   url    string
}{
   {
      entity: "movie",
      format: "4K ULTRA HD",
      url:    "https://disneyplus.com/browse/entity-7df81cf5-6be5-4e05-9ff6-da33baf0b94d",
   },
   {
      entity: "series",
      format: "HD",
      url:    "https://disneyplus.com/browse/entity-21e70fbf-6a51-41b3-88e9-f111830b046c",
   },
}
