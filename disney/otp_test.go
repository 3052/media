package disney

import (
   "fmt"
   "os"
   "testing"
)

const email = "27@riseup.net"

func TestAuthenticateWithOtp(t *testing.T) {
   var device_item Device
   device_item.Token.AccessToken = ""
   resp, err := device_item.authenticate_with_otp(email, 123456)
   if err != nil {
      t.Fatal(err)
   }
   err = resp.Write(os.Stdout)
   if err != nil {
      t.Fatal(err)
   }
}

func TestRequestOtp(t *testing.T) {
   device_item, err := RegisterDevice()
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(device_item.Token.AccessToken)
   otp, err := device_item.RequestOtp(email)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(otp)
}
