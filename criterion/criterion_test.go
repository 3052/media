package criterion

import (
   "fmt"
   "testing"
)

var tests = []struct {
   segment string
   url     string
}{
   {
      segment: "SegmentList",
      url:     "https://criterionchannel.com/wildcat",
   },
   {
      segment: "SegmentTemplate",
      url:     "https://criterionchannel.com/my-dinner-with-andre",
   },
}

func Test(t *testing.T) {
   fmt.Printf("%+v\n", tests)
}
