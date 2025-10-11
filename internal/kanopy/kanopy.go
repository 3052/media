package main

import (
   "41.neocities.org/media/kanopy"
   "41.neocities.org/net"
   "errors"
   "flag"
   "log"
   "os"
   "path/filepath"
)

type flag_set struct {
   cdm      net.Cdm
   email    string
   filters  net.Filters
   kanopy   int
   cache    string
   password string
}

func (f *flag_set) New() error {
   var err error
   f.cache, err = os.UserCacheDir()
   if err != nil {
      return err
   }
   f.cache = filepath.ToSlash(f.cache)
   f.cdm.ClientId = f.cache + "/L3/client_id.bin"
   f.cdm.PrivateKey = f.cache + "/L3/private_key.pem"
   flag.StringVar(&f.cdm.ClientId, "C", f.cdm.ClientId, "client ID")
   flag.StringVar(&f.cdm.PrivateKey, "P", f.cdm.PrivateKey, "private key")
   flag.StringVar(&f.email, "e", "", "email")
   flag.Var(&f.filters, "f", net.FilterUsage)
   flag.IntVar(&f.kanopy, "k", 0, "Kanopy ID")
   flag.StringVar(&f.password, "p", "", "password")
   flag.Parse()
   return nil
}

func (f *flag_set) do_email() error {
   data, err := kanopy.NewLogin(f.email, f.password)
   if err != nil {
      return err
   }
   return write_file(f.cache+"/kanopy/Login", data)
}

///

func main() {
   var set flag_set
   err := set.New()
   if err != nil {
      panic(err)
   }
   func() {
      if set.email != "" {
         if set.password != "" {
            err = set.do_email()
            return
         }
      }
      if set.kanopy >= 1 {
         err = set.do_kanopy()
      } else {
         flag.Usage()
      }
   }()
   if err != nil {
      panic(err)
   }
}

func (f *flag_set) do_kanopy() error {
   data, err := os.ReadFile(f.cache + "/kanopy/Login")
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
   data, err = login.Plays(member, f.kanopy)
   if err != nil {
      return err
   }
   var plays kanopy.Plays
   err = plays.Unmarshal(data)
   if err != nil {
      return err
   }
   err = write_file(f.cache+"/kanopy/Plays", data)
   if err != nil {
      return err
   }
   manifest, ok := plays.Dash()
   if !ok {
      return errors.New(".Dash()")
   }
   resp, err := manifest.Mpd()
   if err != nil {
      return err
   }
   return net.Mpd(f.cache+"/Mpd", resp)
}

func (f *flag_set) do_dash() error {
   data, err := os.ReadFile(f.cache + "/kanopy/Login")
   if err != nil {
      return err
   }
   var login kanopy.Login
   err = login.Unmarshal(data)
   if err != nil {
      return err
   }
   data, err = os.ReadFile(f.cache + "/kanopy/Plays")
   if err != nil {
      return err
   }
   var plays kanopy.Plays
   err = plays.Unmarshal(data)
   if err != nil {
      return err
   }
   manifest, _ := plays.Dash()
   f.cdm.Widevine = func(data []byte) ([]byte, error) {
      return login.Widevine(manifest, data)
   }
   return f.cdm.Download(f.cache+"/Mpd", f.dash)
}

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}
