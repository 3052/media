package disney

import (
   "os"
   "testing"
)

func TestRegisterDevice(t *testing.T) {
   resp, err := register_device()
   if err != nil {
      t.Fatal(err)
   }
   err = resp.Write(os.Stdout)
   if err != nil {
      t.Fatal(err)
   }
}
