package peacock

import (
   "fmt"
   "os"
   "os/exec"
   "testing"
)

func TestSignRead(t *testing.T) {
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(cache + "/peacock/peacock.json")
   if err != nil {
      t.Fatal(err)
   }
   var sign SignIn
   sign.Unmarshal(data)
   auth, err := sign.Auth()
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", auth)
}

func output(name string, arg ...string) (string, error) {
   data, err := exec.Command(name, arg...).Output()
   if err != nil {
      return "", err
   }
   return string(data), nil
}

func TestSignWrite(t *testing.T) {
   user, err := output("credential", "-h=peacocktv.com", "-k=user")
   if err != nil {
      t.Fatal(err)
   }
   password, err := output("credential", "-h=peacocktv.com")
   if err != nil {
      t.Fatal(err)
   }
   var sign SignIn
   err = sign.New(user, password)
   if err != nil {
      t.Fatal(err)
   }
   data, err := sign.Marshal()
   if err != nil {
      t.Fatal(err)
   }
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile(cache + "/peacock/peacock.json", data, os.ModePerm)
   if err != nil {
      t.Fatal(err)
   }
}
func TestVideo(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   text, err := os.ReadFile(home + "/peacock.json")
   if err != nil {
      t.Fatal(err)
   }
   var sign SignIn
   sign.Unmarshal(text)
   auth, err := sign.Auth()
   if err != nil {
      t.Fatal(err)
   }
   video, err := auth.Video(content_id)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", video)
}

// peacocktv.com/watch/playback/vod/GMO_00000000224510_02_HDSDR
const (
   content_id = "GMO_00000000224510_02_HDSDR"
   raw_key_id = "0016e23473ebe77d93d8d1a72dc690d7"
)
