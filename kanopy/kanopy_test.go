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

func TestKanopy(t *testing.T) {
   data, err := exec.Command("password", "kanopy.com").Output()
   if err != nil {
      t.Fatal(err)
   }
   email, password, _ := strings.Cut(string(data), ":")
   data, err = FetchLogin(email, password)
   if err != nil {
      t.Fatal(err)
   }
   var login_var Login
   err = login_var.Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   member, err := login_var.Membership()
   if err != nil {
      t.Fatal(err)
   }
   for _, test := range tests {
      _, err = login_var.Plays(member, test.video_id)
      if err != nil {
         t.Fatal(err)
      }
      time.Sleep(time.Second)
   }
}
