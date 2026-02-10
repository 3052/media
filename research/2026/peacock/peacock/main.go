package main

import (
   "154.pages.dev/log"
   "154.pages.dev/media/internal"
   "flag"
   "os"
   "path/filepath"
)

type flags struct {
   representation string
   email string
   s internal.Stream
   home string
   password string
   peacock string
   v log.Level
}

func (f *flags) New() error {
   var err error
   f.home, err = os.UserHomeDir()
   if err != nil {
      return err
   }
   f.home = filepath.ToSlash(f.home)
   f.s.ClientId = f.home + "/widevine/client_id.bin"
   f.s.PrivateKey = f.home + "/widevine/private_key.pem"
   return nil
}

func main() {
   var f flags
   err := f.New()
   if err != nil {
      panic(err)
   }
   flag.StringVar(&f.peacock, "b", "", "Peacock ID")
   flag.StringVar(&f.email, "e", "", "email")
   flag.StringVar(&f.representation, "i", "", "representation")
   flag.StringVar(&f.password, "p", "", "password")
   flag.TextVar(&f.v.Level, "v", f.v.Level, "level")
   flag.StringVar(&f.s.ClientId, "c", f.s.ClientId, "client ID")
   flag.StringVar(&f.s.PrivateKey, "k", f.s.PrivateKey, "private key")
   flag.Parse()
   f.v.Set()
   log.Transport{}.Set()
   switch {
   case f.password != "":
      err := f.authenticate()
      if err != nil {
         panic(err)
      }
   case f.peacock != "":
      err := f.download()
      if err != nil {
         panic(err)
      }
   default:
      flag.Usage()
   }
}
