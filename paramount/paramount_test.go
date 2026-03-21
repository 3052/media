package paramount

import "testing"

var videos = []struct {
   justWatch  string
   paramount  string
   resolution string
   status     string
}{
   {
      justWatch: "https://justwatch.com/us/tv-show/cia",
      paramount: "https://paramountplus.com/shows/video/8PO2sBBr6lFb7J4nklXuzNZRhUR_V9dd",
      status:    "AVAILABLE",
   },
   {
      justWatch:  "https://justwatch.com/us/movie/zodiac",
      resolution: "2160p",
      paramount:  "https://paramountplus.com/movies/video/wjQ4RChi6BHHu4MVTncppVuCwu44uq2Q",
      status:     "PREMIUM",
   },
   {
      justWatch: "https://justwatch.com/us/tv-show/the-price-is-right",
      paramount: "https://paramountplus.com/shows/video/ALVE01KKH4B7WREZF804N1RV4TSY4S",
      status:    "PREMIUM",
   },
}

func Test(t *testing.T) {
   t.Log(videos)
}
