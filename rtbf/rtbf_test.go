package rtbf

import (
   "154.pages.dev/text"
   "154.pages.dev/widevine"
   "encoding/base64"
   "fmt"
   "os"
   "testing"
   "time"
)

func TestPage(t *testing.T) {
   for _, medium := range media {
      page, err := NewPage(medium.path)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Printf("%+v\n", page)
      name, err := text.Name(page)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Printf("%q\n", name)
      time.Sleep(time.Second)
   }
}

func TestWidevine(t *testing.T) {
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
   medium := media[0]
   var pssh widevine.Pssh
   pssh.KeyId, err = base64.StdEncoding.DecodeString(medium.key_id)
   if err != nil {
      t.Fatal(err)
   }
   var module widevine.Cdm
   err = module.New(private_key, client_id, pssh.Encode())
   if err != nil {
      t.Fatal(err)
   }
   text, err := os.ReadFile("account.txt")
   if err != nil {
      t.Fatal(err)
   }
   var account AccountLogin
   err = account.Unmarshal(text)
   if err != nil {
      t.Fatal(err)
   }
   token, err := account.Token()
   if err != nil {
      t.Fatal(err)
   }
   gigya, err := token.Login()
   if err != nil {
      t.Fatal(err)
   }
   page, err := NewPage(medium.path)
   if err != nil {
      t.Fatal(err)
   }
   title, err := gigya.Entitlement(page)
   if err != nil {
      t.Fatal(err)
   }
   key, err := module.Key(title, pssh.KeyId)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%x\n", key)
}

func TestWebToken(t *testing.T) {
   text, err := os.ReadFile("account.txt")
   if err != nil {
      t.Fatal(err)
   }
   var account AccountLogin
   err = account.Unmarshal(text)
   if err != nil {
      t.Fatal(err)
   }
   token, err := account.Token()
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", token)
}
func TestEntitlement(t *testing.T) {
   text, err := os.ReadFile("account.txt")
   if err != nil {
      t.Fatal(err)
   }
   var account AccountLogin
   err = account.Unmarshal(text)
   if err != nil {
      t.Fatal(err)
   }
   token, err := account.Token()
   if err != nil {
      t.Fatal(err)
   }
   gigya, err := token.Login()
   if err != nil {
      t.Fatal(err)
   }
   page, err := NewPage(media[0].path)
   if err != nil {
      t.Fatal(err)
   }
   title, err := gigya.Entitlement(page)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", title)
   fmt.Println(title.Dash())
}
