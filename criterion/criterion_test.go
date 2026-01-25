package criterion

import "testing"

var videos = []struct {
   segment string
   url     string
}{
   {
      url: "https://criterionchannel.com/121280-ritual",
   },
   {
      segment: "SegmentList",
      url:     "https://criterionchannel.com/wildcat",
   },
   {
      segment: "SegmentTemplate",
      url:     "https://criterionchannel.com/my-dinner-with-andre",
   },
}

func TestVideo(t *testing.T) {
   t.Log(videos)
}
