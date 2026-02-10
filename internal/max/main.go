package main

import (
   "41.neocities.org/media/internal"
   "41.neocities.org/media/max"
   "flag"
   "log"
   "os"
   "path/filepath"
)

func main() {
   var f flags
   err := f.New()
   if err != nil {
      panic(err)
   }
   flag.StringVar(&f.s.ClientId, "c", f.s.ClientId, "client ID")
   flag.StringVar(&f.representation, "i", "", "representation")
   flag.BoolVar(
      &f.initiate, "initiate", false, "/authentication/linkDevice/initiate",
   )
   flag.StringVar(&f.s.PrivateKey, "k", f.s.PrivateKey, "private key")
   flag.BoolVar(
      &f.login, "login", false, "/authentication/linkDevice/login",
   )
   flag.Var(&f.url, "a", "address")
   flag.Parse()
   switch {
   case f.initiate:
      err := f.do_initiate()
      if err != nil {
         panic(err)
      }
   case f.login:
      err := f.do_login()
      if err != nil {
         panic(err)
      }
   case f.url.VideoId != "":
      err := f.download()
      if err != nil {
         panic(err)
      }
   default:
      flag.Usage()
   }
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

type flags struct {
   home           string
   initiate       bool
   login          bool
   representation string
   s              internal.Stream
   url            max.WatchUrl
}
