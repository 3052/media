package criterion

import (
   "os"
   "testing"
   "time"
)

var tests = []struct {
   key_id string
   slug   string
   url    string
}{
   {
      key_id: "",
      slug:   "wildcat",
      url:    "criterionchannel.com/wildcat",
   },
   {
      key_id: "e4576465a745213f336c1ef1bf5d513e",
      slug:   "my-dinner-with-andre",
      url:    "criterionchannel.com/videos/my-dinner-with-andre",
   },
}

func Test(t *testing.T) {
   data, err := os.ReadFile("token.txt")
   if err != nil {
      t.Fatal(err)
   }
   var token1 Token
   err = token1.Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   for _, test1 := range tests {
      _, err = token1.Video(test1.slug)
      if err != nil {
         t.Fatal(err)
      }
      time.Sleep(time.Second)
   }
}
