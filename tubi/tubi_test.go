package tubi

import "testing"

var tests = []struct {
   resolution string
   url        string
}{
   {
      url:        "https://tubitv.com/movies/714654",
      resolution: "1080p",
   },
   {
      url:        "https://tubitv.com/movies/617502",
      resolution: "720p",
   },
   {
      url:        "https://tubitv.com/tv-shows/200203258",
      resolution: "720p",
   },
}

func Test(t *testing.T) {
   t.Log(tests)
}
