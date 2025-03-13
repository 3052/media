package binge

import (
   "41.neocities.org/widevine"
   "fmt"
   "os"
   "os/exec"
   "strings"
   "testing"
)

const (
   asset_id = 7738
   key_id = "\x1e\v\xdc\xfd\x069Nٝ\x9f\xe6#\xb0<\xa6\xd2"
)

func TestRefresh(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(home + "/media/binge/Auth")
   if err != nil {
      t.Fatal(err)
   }
   var auth1 Auth
   err = auth1.Unmarshal(data)
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
   data, err = NewAuth(username, password)
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile(home+"/media/binge/Auth", data, os.ModePerm)
   if err != nil {
      t.Fatal(err)
   }
}

func TestLicense(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(home + "/media/binge/Auth")
   if err != nil {
      t.Fatal(err)
   }
   var auth1 Auth
   err = auth1.Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   token, err := auth1.Token()
   if err != nil {
      t.Fatal(err)
   }
   play1, err := token.Play(asset_id)
   if err != nil {
      t.Fatal(err)
   }
   stream1, ok := play1.Dash()
   if !ok {
      t.Fatal(".Dash()")
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
   _, err = token.Widevine(stream1, data)
   if err != nil {
      t.Fatal(err)
   }
}

func TestPlay(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(home + "/media/binge/Auth")
   if err != nil {
      t.Fatal(err)
   }
   var auth1 Auth
   err = auth1.Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   token, err := auth1.Token()
   if err != nil {
      t.Fatal(err)
   }
   play1, err := token.Play(asset_id)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", play1)
}
