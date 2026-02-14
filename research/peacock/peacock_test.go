package peacock

import (
   "fmt"
   "net/http"
   "os"
   "os/exec"
   "testing"
)

// peacocktv.com/watch/playback/vod/GMO_00000000091566_02_HDSDR/6668f89a-b581-36ac-9895-7783aa16b471
const content_id = "GMO_00000000091566_02_HDSDR"

func TestPlayout(t *testing.T) {
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(cache + "/peacock/peacock.txt")
   if err != nil {
      t.Fatal(err)
   }
   id, err := http.ParseSetCookie(string(data))
   if err != nil {
      t.Fatal(err)
   }
   var token AuthToken
   err = token.Fetch(id)
   if err != nil {
      t.Fatal(err)
   }
   play, err := token.Playout(content_id)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", play)
}

func TestSignRead(t *testing.T) {
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   data, err := os.ReadFile(cache + "/peacock/peacock.txt")
   if err != nil {
      t.Fatal(err)
   }
   id, err := http.ParseSetCookie(string(data))
   if err != nil {
      t.Fatal(err)
   }
   var token AuthToken
   err = token.Fetch(id)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", token)
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
   id, err := FetchIdSession(user, password)
   if err != nil {
      t.Fatal(err)
   }
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   err = os.WriteFile(
      cache+"/peacock/peacock.txt", []byte(id.String()), os.ModePerm,
   )
   if err != nil {
      t.Fatal(err)
   }
}
