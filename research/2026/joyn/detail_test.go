package joyn

import (
   "154.pages.dev/encoding"
   "fmt"
   "testing"
   "time"
)

var tests = []struct {
   key_id string
   path   string
}{
   {
      // joyn.de/filme/barry-seal-only-in-america
      key_id: "9mY0MZrt58qhF/FvD837QA==",
      path:   "/filme/barry-seal-only-in-america",
   },
   {
      // joyn.de/serien/one-tree-hill/1-2-quaelende-angst
      path: "/serien/one-tree-hill/1-2-quaelende-angst",
   },
}

func TestDetail(t *testing.T) {
   for _, test := range tests {
      detail, err := Path(test.path).Detail()
      if err != nil {
         t.Fatal(err)
      }
      name, err := encoding.Name(Namer{detail})
      if err != nil {
         t.Fatal(err)
      }
      fmt.Printf("%q\n", name)
      fmt.Printf("%+v\n", detail)
      time.Sleep(time.Second)
   }
}
