package main

import (
   "41.neocities.org/media/internal"
   "41.neocities.org/media/rakuten"
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
   flag.Var(&f.address, "a", "address")
   flag.StringVar(&f.language, "b", "", "language")
   flag.StringVar(&f.s.ClientId, "c", f.s.ClientId, "client ID")
   flag.StringVar(&f.representation, "i", "", "representation")
   flag.StringVar(&f.s.PrivateKey, "k", f.s.PrivateKey, "private key")
   flag.Parse()
   if f.address.MarketCode != "" {
      if f.language != "" {
         err := f.download()
         if err != nil {
            panic(err)
         }
      } else {
         err := f.do_language()
         if err != nil {
            panic(err)
         }
      }
   } else {
      flag.Usage()
   }
}

func (f *flags) New() error {
   home, err := os.UserHomeDir()
   if err != nil {
      return err
   }
   home = filepath.ToSlash(home)
   f.s.ClientId = home + "/widevine/client_id.bin"
   f.s.PrivateKey = home + "/widevine/private_key.pem"
   return nil
}

type flags struct {
   address        rakuten.Address
   representation string
   s              internal.Stream
   language       string
}
