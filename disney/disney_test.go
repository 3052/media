package disney

import (
   "fmt"
   "testing"
)

func TestRegisterDevice(t *testing.T) {
   var token_value Token
   err := token_value.RegisterDevice()
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(&token_value)
}

func TestAuthenticateWithOtp(t *testing.T) {
   var token_value Token
   token_value.AccessToken = otp_test.access_token
   authenticate, err := token_value.AuthenticateWithOtp(
      otp_test.email, otp_test.passcode,
   )
   if err != nil {
      t.Fatal(err)
   }
   inactive, err := token_value.LoginWithActionGrant(authenticate.ActionGrant)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", inactive)
}

func TestEntity(t *testing.T) {
   t.Log(entity_tests)
}

var entity_tests = []struct {
   entity string
   format string
   url    string
}{
   {
      entity: "movie",
      format: "4K ULTRA HD",
      url:    "https://disneyplus.com/browse/entity-7df81cf5-6be5-4e05-9ff6-da33baf0b94d",
   },
   {
      entity: "movie",
      format: "4K ULTRA HD",
      url:    "https://disneyplus.com/browse/entity-917f1bf3-3db4-4df0-afe2-60b2c5e67618",
   },
   {
      entity: "series",
      format: "HD",
      url:    "https://disneyplus.com/browse/entity-21e70fbf-6a51-41b3-88e9-f111830b046c",
   },
}

func TestRequestOtp(t *testing.T) {
   var token_value Token
   token_value.AccessToken = otp_test.access_token
   otp, err := token_value.RequestOtp(otp_test.email)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(otp)
}

var otp_test struct {
   email        string
   passcode     string
   access_token string
}
