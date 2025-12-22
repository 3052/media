package itv

import (
   "os"
   "testing"
)

func Test(t *testing.T) {
   data, err := os.ReadFile("10a3918.txt")
   if err != nil {
      t.Fatal(err)
   }
   var next next_data
   err = next.ExtractFromHTML(data)
   if err != nil {
      t.Fatal(err)
   }
   t.Logf("%+v", next)
}
