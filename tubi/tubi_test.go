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
      content_id: 100047438,
      drm:        true,
      url:        "tubitv.com/movies/100047438",
   },
   {
      content_id: 200042567,
      url:        "tubitv.com/tv-shows/200042567",
      drm:        true,
   },
   {
      content_id: 667315,
      url:        "tubitv.com/movies/667315",
      drm:        false,
   },
   {
      content_id: 643397,
      location:   "Australia",
      url:        "tubitv.com/movies/643397",
      drm:        false,
   },
}

func Test(t *testing.T) {
   for _, test_var := range tests {
      data, err := NewContent(test_var.content_id)
      if err != nil {
         t.Fatal(err)
      }
      content_var := &Content{}
      err = content_var.Unmarshal(data)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(content_var)
      time.Sleep(time.Second)
   }
}
