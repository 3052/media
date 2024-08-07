package roku

import (
   "154.pages.dev/widevine"
   "encoding/hex"
   "fmt"
   "os"
   "path"
   "testing"
   "time"
)

func TestLicense(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   private_key, err := os.ReadFile(home + "/widevine/private_key.pem")
   if err != nil {
      t.Fatal(err)
   }
   client_id, err := os.ReadFile(home + "/widevine/client_id.bin")
   if err != nil {
      t.Fatal(err)
   }
   for _, test := range tests{
      var pssh widevine.Pssh
      pssh.KeyId, err = hex.DecodeString(test.key_id)
      if err != nil {
         t.Fatal(err)
      }
      var module widevine.Cdm
      err = module.New(private_key, client_id, pssh.Encode())
      if err != nil {
         t.Fatal(err)
      }
      var auth AccountAuth
      auth.New(nil)
      play, err := auth.Playback(path.Base(test.url))
      if err != nil {
         t.Fatal(err)
      }
      key, err := module.Key(play, pssh.KeyId)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Printf("%x\n", key)
      time.Sleep(time.Second)
   }
}

func TestCode(t *testing.T) {
   // AccountAuth
   var auth AccountAuth
   err := auth.New(nil)
   if err != nil {
      t.Fatal(err)
   }
   os.WriteFile("auth.txt", auth.Data, 0666)
   auth.Unmarshal()
   // AccountCode
   code, err := auth.Code()
   if err != nil {
      t.Fatal(err)
   }
   os.WriteFile("code.txt", code.Data, 0666)
   code.Unmarshal()
   fmt.Println(code)
}

func TestTokenWrite(t *testing.T) {
   var err error
   // AccountAuth
   var auth AccountAuth
   auth.Data, err = os.ReadFile("auth.txt")
   if err != nil {
      t.Fatal(err)
   }
   auth.Unmarshal()
   // AccountCode
   var code AccountCode
   code.Data, err = os.ReadFile("code.txt")
   if err != nil {
      t.Fatal(err)
   }
   code.Unmarshal()
   // AccountToken
   token, err := auth.Token(code)
   if err != nil {
      t.Fatal(err)
   }
   os.WriteFile("token.txt", token.Data, 0666)
}

func TestTokenRead(t *testing.T) {
   var err      error
   // AccountToken
   var token AccountToken
   token.Data, err = os.ReadFile("token.txt")
   if err != nil {
      t.Fatal(err)
   }
   token.Unmarshal()
   // AccountAuth
   var auth AccountAuth
   err = auth.New(&token)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", auth)
}
