package tubi

import (
   "fmt"
   "testing"
   "time"
)

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

func TestLicense(t *testing.T) {
   for _, test := range tests {
      content1 := &Content{}
      data, err := content1.Marshal(test.content_id)
      if err != nil {
         t.Fatal(err)
      }
      err = content1.Unmarshal(data)
      if err != nil {
         t.Fatal(err)
      }
      if content1.Episode() {
         data, err = content1.Marshal(content1.SeriesId)
         if err != nil {
            t.Fatal(err)
         }
         err = content1.Unmarshal(data)
         if err != nil {
            t.Fatal(err)
         }
         var ok bool
         content1, ok = content1.Get(test.content_id)
         if !ok {
            t.Fatal("Content.Get")
         }
      }
      fmt.Println(content1)
      time.Sleep(time.Second)
   }
}
