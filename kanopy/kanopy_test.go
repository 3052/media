package kanopy

import (
   "os/exec"
   "strings"
   "testing"
   "time"
)

var tests = []struct {
   key_id   string
   url      string
   video_id int
}{
   {
      key_id:   "DUCS1DH4TB6Po1oEkG9xUA==",
      url:      "kanopy.com/en/product/13808102",
      video_id: 13808102,
   },
   {
      key_id:   "sYcEuBtnTH6Bqn65yIE0Ww==",
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
   var login1 Login
   data, err = login1.Marshal(email, password)
   if err != nil {
      t.Fatal(err)
   }
   err = login1.Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   member, err := login1.Membership()
   if err != nil {
      t.Fatal(err)
   }
   for _, test := range tests {
      _, err = login1.Plays(member, test.video_id)
      if err != nil {
         t.Fatal(err)
      }
      time.Sleep(time.Second)
   }
}
