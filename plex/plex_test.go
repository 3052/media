package plex

import "testing"

var watch_tests = []struct {
   drm bool
   url string
}{
   {
      drm: true,
      url: "https://watch.plex.tv/watch/movie/ghost-in-the-shell",
   },
   {
      url: "https://watch.plex.tv/movie/limitless",
   },
   {
      url: "https://watch.plex.tv/show/broadchurch/season/3/episode/5",
   },
}

func Test(t *testing.T) {
   t.Log(watch_tests)
}
