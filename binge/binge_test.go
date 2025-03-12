package binge

import (
   "41.neocities.org/widevine"
   "fmt"
   "os"
   "os/exec"
   "strings"
   "testing"
)

func TestLicense(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(home + "/media/binge/auth")
   if err != nil {
      t.Fatal(err)
   }
   var auth1 auth
   err = auth1.unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   token, err := auth1.token()
   if err != nil {
      t.Fatal(err)
   }
   play1, err := token.play()
   if err != nil {
      t.Fatal(err)
   }
   stream1, ok := play1.dash()
   if !ok {
      t.Fatal(".dash()")
   }
   private_key, err := os.ReadFile(home + "/media/private_key.pem")
   if err != nil {
      t.Fatal(err)
   }
   client_id, err := os.ReadFile(home + "/media/client_id.bin")
   if err != nil {
      t.Fatal(err)
   }
   var pssh widevine.Pssh
   pssh.KeyIds = [][]byte{
      []byte(key_id),
   }
   var cdm widevine.Cdm
   err = cdm.New(private_key, client_id, pssh.Marshal())
   if err != nil {
      t.Fatal(err)
   }
   data, err = cdm.RequestBody()
   if err != nil {
      t.Fatal(err)
   }
   _, err = token.widevine(stream1, data)
   if err != nil {
      t.Fatal(err)
   }
}

const key_id = "\x1e\v\xdc\xfd\x069Nٝ\x9f\xe6#\xb0<\xa6\xd2"
func TestToken(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(home + "/media/binge/auth")
   if err != nil {
      t.Fatal(err)
   }
   var auth1 auth
   err = auth1.unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Print(auth1.AccessToken, "\n\n")
   token, err := auth1.token()
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(token.AccessToken)
}

func TestRefresh(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(home + "/media/binge/auth")
   if err != nil {
      t.Fatal(err)
   }
   var auth1 auth
   err = auth1.unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Print(auth1.AccessToken, "\n\n")
   err = auth1.refresh()
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(auth1.AccessToken)
}

func TestWrite(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := exec.Command("password", "binge.com.au").Output()
   if err != nil {
      t.Fatal(err)
   }
   username, password, _ := strings.Cut(string(data), ":")
   data, err = new_auth(username, password)
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile(home + "/media/binge/auth", data, os.ModePerm)
   if err != nil {
      t.Fatal(err)
   }
}
