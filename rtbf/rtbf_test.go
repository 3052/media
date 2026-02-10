package rtbf

import (
   "fmt"
   "testing"
   "time"
)

var tests = []struct {
   key_id string
   path   string
   url    string
}{
   {
      key_id: "Ma5jT/1dR8K/ljWx/1Pb4A==",
      path:   "/media/titanic-3286058",
      url:    "auvio.rtbf.be/media/titanic-3286058",
   },
   {
      key_id: "xESyRLihQMacu++BvoakfA==",
      path:   "/media/agatha-christie-pourquoi-pas-evans-agatha-christie-pourquoi-pas-evans-3280380",
      url:    "auvio.rtbf.be/media/agatha-christie-pourquoi-pas-evans-agatha-christie-pourquoi-pas-evans-3280380",
   },
}

func TestPage(t *testing.T) {
   for _, test := range tests {
      content1, err := Address{test.path}.Content()
      if err != nil {
         t.Fatal(err)
      }
      fmt.Printf("%+v\n", content1)
      time.Sleep(time.Second)
   }
}
