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
      content_id: 312926,
      url:        "tubitv.com/movies/312926",
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
   for _, test1 := range tests {
      data, err := NewContent(test1.content_id)
      if err != nil {
         t.Fatal(err)
      }
      content1 := &Content{}
      err = content1.Unmarshal(data)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(content1)
      time.Sleep(time.Second)
   }
}
