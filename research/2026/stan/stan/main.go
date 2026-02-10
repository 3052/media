package main

import (
   "154.pages.dev/media/internal"
   "154.pages.dev/media/stan"
   "154.pages.dev/text"
   "flag"
   "os"
   "path/filepath"
   "strings"
)

type flags struct {
   code bool
   home string
   host string
   stan int64
   representation string
   s internal.Stream
   token bool
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
   flag.Int64Var(&f.stan, "b", 0, "Stan ID")
   flag.StringVar(&f.s.ClientId, "c", f.s.ClientId, "client ID")
   flag.BoolVar(&f.code, "code", false, "activation code")
   flag.StringVar(
      &f.host, "h", stan.BaseUrl[0], strings.Join(stan.BaseUrl[1:], "\n"),
   )
   flag.StringVar(&f.representation, "i", "", "representation")
   flag.StringVar(&f.s.PrivateKey, "p", f.s.PrivateKey, "private key")
   flag.BoolVar(&f.token, "token", false, "web token")
   flag.Parse()
   text.Transport{}.Set(true)
   switch {
   case f.code:
      err := f.write_code()
      if err != nil {
         panic(err)
      }
   case f.token:
      err := f.write_token()
      if err != nil {
         panic(err)
      }
   case f.stan >= 1:
      err := f.download()
      if err != nil {
         panic(err)
      }
   default:
      flag.Usage()
   }
}
