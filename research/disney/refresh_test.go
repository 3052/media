package disney

import (
   "os"
   "testing"
)

func Test(t *testing.T) {
   resp, err := refresh()
   if err != nil {
      t.Fatal(err)
   }
   err = resp.Write(os.Stdout)
   if err != nil {
      t.Fatal(err)
   }
}
