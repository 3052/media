package disney

import (
   "fmt"
   "testing"
)

func TestRegisterDevice(t *testing.T) {
   device, err := fetch_register_device()
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", device)
}
