package criterion

import (
   "fmt"
   "os/exec"
   "testing"
)

func TestToken(t *testing.T) {
   user, err := output("credential", "-h", "criterionchannel.com", "-k", "user")
   if err != nil {
      t.Fatal(err)
   }
   password, err := output("credential", "-h", "criterionchannel.com")
   if err != nil {
      t.Fatal(err)
   }
   data, err := FetchToken(user, password)
   if err != nil {
      t.Fatal(err)
   }
   var token_value Token
   err = token_value.Unmarshal(data)
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

var videos = []struct {
   segment string
   url     string
}{
   {
      segment: "SegmentList",
      url:     "https://criterionchannel.com/wildcat",
   },
   {
      segment: "SegmentTemplate",
      url:     "https://criterionchannel.com/my-dinner-with-andre",
   },
}

func TestVideo(t *testing.T) {
   fmt.Printf("%+v\n", videos)
}
