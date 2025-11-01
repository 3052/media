package tubi

import (
   "fmt"
   "testing"
   "time"
)

var tests = []struct {
   content_id int
   drm        bool
   location   string
   url        string
}{
   {
      content_id: 667315,
      url:        "tubitv.com/movies/667315",
      drm:        false,
   },
   {
      content_id: 200042567,
      url:        "tubitv.com/tv-shows/200042567",
      drm:        true,
   },
   {
      content_id: 643397,
      location:   "Australia",
      url:        "tubitv.com/movies/643397",
      drm:        false,
   },
}

func Test(t *testing.T) {
   for _, testVar := range tests {
      data, err := NewContent(testVar.content_id)
      if err != nil {
         t.Fatal(err)
      }
      contentVar := &Content{}
      err = contentVar.Unmarshal(data)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(contentVar)
      time.Sleep(time.Second)
   }
}
