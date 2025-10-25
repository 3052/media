package criterion

import (
   "fmt"
   "testing"
)

var tests = []struct {
   segment string
   slug    string
   url     string
}{
   {
      segment: "SegmentList",
      slug:    "wildcat",
      url:     "criterionchannel.com/wildcat",
   },
   {
      segment: "SegmentTemplate",
      slug:    "my-dinner-with-andre",
      url:     "criterionchannel.com/videos/my-dinner-with-andre",
   },
}

func Test(t *testing.T) {
   fmt.Printf("%+v\n", tests)
}
