package tubi

import (
   "fmt"
   "testing"
   "time"
)

var tests = []struct {
   content_id int
   location   string
   url        string
}{
   {
      content_id: 100003573,
      url:        "tubitv.com/movies/100003573",
   },
   {
      content_id: 200042567,
      url:        "tubitv.com/tv-shows/200042567",
   },
   {
      content_id: 643397,
      location:   "Australia",
      url:        "tubitv.com/movies/643397",
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
