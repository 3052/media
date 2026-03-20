package tubi

import "testing"

var tests = []struct {
   drm        bool
   resolution string
   url        string
}{
   {
      drm:        true,
      resolution: "720p",
      url:        "https://tubitv.com/movies/617502",
   },
   {
      drm:        false,
      resolution: "1080p",
      url:        "https://tubitv.com/movies/714654",
   },
   {
      drm:        false,
      resolution: "720p",
      url:        "https://tubitv.com/tv-shows/200203258",
   },
}

func Test(t *testing.T) {
   t.Log(tests)
}
