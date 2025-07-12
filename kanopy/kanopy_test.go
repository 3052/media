package kanopy

import (
   "os/exec"
   "strings"
   "testing"
   "time"
)

var tests = []struct {
   url      string
   video_id int
}{
   {
      url:      "kanopy.com/en/product/13808102",
      video_id: 13808102,
   },
   {
      url:      "kanopy.com/en/product/14881167",
      video_id: 14881167,
   },
}

func Test(t *testing.T) {
   data, err := exec.Command("password", "kanopy.com").Output()
   if err != nil {
      t.Fatal(err)
   }
   email, password, _ := strings.Cut(string(data), ":")
   data, err = NewLogin(email, password)
   if err != nil {
      t.Fatal(err)
   }
   var loginVar Login
   err = loginVar.Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   member, err := loginVar.Membership()
   if err != nil {
      t.Fatal(err)
   }
   for _, testVar := range tests {
      _, err = loginVar.Plays(member, testVar.video_id)
      if err != nil {
         t.Fatal(err)
      }
      time.Sleep(time.Second)
   }
}
