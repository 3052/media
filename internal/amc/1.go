package main

import (
   "41.neocities.org/media/amc"
   "41.neocities.org/media/internal"
   "errors"
   "flag"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

func main() {
   var f flags
   err := f.New()
   if err != nil {
      panic(err)
   }
   flag.StringVar(&f.e.ClientId, "client", f.e.ClientId, "client ID")
   flag.StringVar(&f.email, "email", "", "email")
   flag.StringVar(&f.e.PrivateKey, "key", f.e.PrivateKey, "private key")
   flag.StringVar(&f.password, "password", "", "password")
   
   flag.Int64Var(&f.episode, "e", 0, "episode or movie ID")
   flag.Int64Var(&f.season, "s", 0, "season ID")
   flag.Int64Var(&f.series, "series", 0, "series ID")
   
   
   flag.Parse()
   if f.email != "" {
      if f.password != "" {
         err := f.do_email_password()
         if err != nil {
            panic(err)
         }
      }
   } else if f.amc >= 1 {
      err := f.do_amc()
      if err != nil {
         panic(err)
      }
   } else if f.dash != "" {
      err := f.do_dash()
      if err != nil {
         panic(err)
      }
   } else {
      flag.Usage()
   }
}

func (f *flags) do_amc() error {
   data, err := os.ReadFile(f.media + "/amc/Auth")
   if err != nil {
      return err
   }
   var auth amc.Auth
   err = auth.Unmarshal(data)
   if err != nil {
      return err
   }
   data, err = auth.Refresh()
   if err != nil {
      return err
   }
   err = write_file(f.media+"/amc/Auth", data)
   if err != nil {
      return err
   }
   err = auth.Unmarshal(data)
   if err != nil {
      return err
   }
   data, err = auth.Playback(f.amc)
   if err != nil {
      return err
   }
   err = write_file(f.media+"/amc/Playback", data)
   if err != nil {
      return err
   }
   var play amc.Playback
   err = play.Unmarshal(data)
   if err != nil {
      return err
   }
   source, ok := play.Dash()
   if !ok {
      return errors.New(".Dash()")
   }
   resp, err := http.Get(source.Src)
   if err != nil {
      return err
   }
   return internal.Mpd(f.media+"/Mpd", resp)
}

func (f *flags) do_dash() error {
   data, err := os.ReadFile(f.media + "/amc/Playback")
   if err != nil {
      return err
   }
   var play amc.Playback
   err = play.Unmarshal(data)
   if err != nil {
      return err
   }
   source, _ := play.Dash()
   f.e.Widevine = func(data []byte) ([]byte, error) {
      return play.Widevine(source, data)
   }
   return f.e.Download(f.media+"/Mpd", f.dash)
}
