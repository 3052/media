package main

import (
   "41.neocities.org/media/internal"
   "41.neocities.org/media/kanopy"
   "errors"
   "fmt"
   "os"
)

func (f *flags) download() error {
   data, err := os.ReadFile(f.home + "/kanopy.txt")
   if err != nil {
      return err
   }
   var login kanopy.Login
   err = login.Unmarshal(data)
   if err != nil {
      return err
   }
   member, err := login.Membership()
   if err != nil {
      return err
   }
   plays, err := login.Plays(member, f.video_id)
   if err != nil {
      return err
   }
   manifest, ok := plays.Dash()
   if !ok {
      return errors.New("VideoPlays.Dash")
   }
   represents, err := internal.Mpd(manifest)
   if err != nil {
      return err
   }
   for _, represent := range represents {
      switch f.representation {
      case "":
         fmt.Print(&represent, "\n\n")
      case represent.Id:
         f.s.Client = &kanopy.Client{manifest, &login}
         return f.s.Download(&represent)
      }
   }
   return nil
}

func (f *flags) authenticate() error {
   data, err := kanopy.Login{}.Marshal(f.email, f.password)
   if err != nil {
      return err
   }
   return os.WriteFile(f.home+"/kanopy.txt", data, os.ModePerm)
}
