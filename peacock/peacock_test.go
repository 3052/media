package peacock

import (
   "net/http"
   "os"
   "os/exec"
   "testing"
)

var watch = struct {
   content_id string
   url        string
}{
   content_id: "GMO_00000000091566_02_HDSDR",
   url:        "https://peacocktv.com/watch/playback/vod/GMO_00000000091566_02_HDSDR/6668f89a-b581-36ac-9895-7783aa16b471",
}

func TestWatch(t *testing.T) {
   t.Log(watch)
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
   var auth Token
   err = auth.Fetch(id)
   if err != nil {
      t.Fatal(err)
   }
   t.Log(auth)
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

func output(name string, arg ...string) (string, error) {
   data, err := exec.Command(name, arg...).Output()
   if err != nil {
      return "", err
   }
   return string(data), nil
}
