package disney

import (
   "io"
   "os"
   "testing"
)

func TestRegisterDevice(t *testing.T) {
   resp, err := register_device()
   if err != nil {
      t.Fatal(err)
   }
   defer resp.Body.Close()
   data, err := io.ReadAll(resp.Body)
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile("register_device.json", data, os.ModePerm)
   if err != nil {
      t.Fatal(err)
   }
}
