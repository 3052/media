package main

import (
   "41.neocities.org/media/ctv"
   "41.neocities.org/net"
   "flag"
   "net/http"
   "os"
   "path/filepath"
)

func (f *flag_set) do_address() error {
   path, err := ctv.GetPath(f.address)
   if err != nil {
      return err
   }
   resolve, err := ctv.Resolve(path)
   if err != nil {
      return err
   }
   axis, err := resolve.Axis()
   if err != nil {
      return err
   }
   content, err := axis.Content()
   if err != nil {
      return err
   }
   address, err := axis.Mpd(content)
   if err != nil {
      return err
   }
   resp, err := http.Get(address)
   if err != nil {
      return err
   }
   f.config.Send = ctv.Widevine
   return f.filters.Filter(resp, &f.config)
}

type flag_set struct {
   address string
   config  net.Config
   filters net.Filters
}

func (f *flag_set) New() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   cache = filepath.ToSlash(cache)
   f.config.ClientId = cache + "/L3/client_id.bin"
   f.config.PrivateKey = cache + "/L3/private_key.pem"
   flag.StringVar(&f.address, "a", "", "address")
   flag.StringVar(&f.config.ClientId, "c", f.config.ClientId, "client ID")
   flag.Var(&f.filters, "f", net.FilterUsage)
   flag.StringVar(&f.config.PrivateKey, "p", f.config.PrivateKey, "private key")
   flag.Parse()
   return nil
}

func main() {
   var set flag_set
   err := set.New()
   if err != nil {
      panic(err)
   }
   if set.address != "" {
      err = set.do_address()
      if err != nil {
         panic(err)
      }
   } else {
      flag.Usage()
   }
}
