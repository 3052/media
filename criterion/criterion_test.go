package criterion

import (
   "os"
   "testing"
   "time"
)

var tests = []struct {
   slug    string
   url     string
   segment string
}{
   {
      slug:    "wildcat",
      url:     "criterionchannel.com/wildcat",
      segment: "SegmentList",
   },
   {
      slug:    "my-dinner-with-andre",
      url:     "criterionchannel.com/videos/my-dinner-with-andre",
      segment: "SegmentTemplate",
   },
}

func Test(t *testing.T) {
   data, err := os.ReadFile("token.txt")
   if err != nil {
      t.Fatal(err)
   }
   var tokenVar Token
   err = tokenVar.Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   for _, testVar := range tests {
      _, err = tokenVar.Video(testVar.slug)
      if err != nil {
         t.Fatal(err)
      }
      time.Sleep(time.Second)
   }
}
