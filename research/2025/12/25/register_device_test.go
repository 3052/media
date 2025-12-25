package disney

import (
   "fmt"
   "testing"
)

func TestRegisterDevice(t *testing.T) {
   var device register_device
   err := device.fetch()
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", device)
}
