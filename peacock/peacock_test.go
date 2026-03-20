package peacock

import (
   "os"
   "os/exec"
   "testing"
)

var bugonia = []string{
   "https://peacocktv.com/watch/asset/movies/bugonia/c84393dc-6aca-3466-b3cd-76f44c79a236",
   "https://peacocktv.com/watch/playback/vod/GMO_00000000261361_02_HDSDR/c84393dc-6aca-3466-b3cd-76f44c79a236",
}

func TestWatch(t *testing.T) {
   t.Log(bugonia)
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
