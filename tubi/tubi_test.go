package tubi

import (
   "fmt"
   "testing"
   "time"
)

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

var tests = []struct {
   content_id int
   key_id     string
   location   string
   url        string
}{
   {
      content_id: 643397,
      location:   "Australia",
      url:        "tubitv.com/movies/643397",
   },
   {
      content_id: 100001047,
      url:        "tubitv.com/movies/100001047",
   },
   {
      content_id: 200042567,
      key_id:     "Ndopo1ozQ8iSL75MAfbL6A==",
      url:        "tubitv.com/tv-shows/200042567",
   },
}
