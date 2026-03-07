package disney

import (
   "os"
   "testing"
)

func Test(t *testing.T) {
   resp, err := request_otp()
   if err != nil {
      t.Fatal(err)
   }
   err = resp.Write(os.Stdout)
   if err != nil {
      t.Fatal(err)
   }
}
