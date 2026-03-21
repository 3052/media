package paramount

import "testing"

var videos = []struct {
   resolution string
   paramount        string
   justWatch string
}{
   {
      justWatch: "https://justwatch.com/us/movie/zodiac",
      paramount:        "https://paramountplus.com/movies/video/wjQ4RChi6BHHu4MVTncppVuCwu44uq2Q",
      resolution: "2160p",
   },
   {
      justWatch: "https://justwatch.com/us/tv-show/criminal-minds",
      paramount:        "https://paramountplus.com/shows/video/esJvFlqdrcS_kFHnpxSuYp449E7tTexD",
      resolution: "2160p",
   },
   {
      justWatch: "https://justwatch.com/us/movie/scream-1",
      paramount:        "https://paramountplus.com/gb/movies/video/yh8qG9949D8fvJa0dFmY1C_SilOYt2hS",
      resolution: "2160p",
   },
   {
      justWatch: "https://justwatch.com/us/movie/paw-patrol-ready-race-rescue",
   },
}

func Test(t *testing.T) {
   t.Log(videos)
}
