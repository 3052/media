package main

import (
   "41.neocities.org/media/plex"
   "41.neocities.org/net"
   "errors"
   "flag"
   "os"
   "path/filepath"
)

type flag_set struct {
   address string
   cache   string
   config  net.Config
   filters net.Filters
}

func (f *flag_set) New() error {
   var err error
   f.cache, err = os.UserCacheDir()
   if err != nil {
      return err
   }
   f.cache = filepath.ToSlash(f.cache)
   f.config.ClientId = f.cache + "/L3/client_id.bin"
   f.config.PrivateKey = f.cache + "/L3/private_key.pem"
   flag.StringVar(&f.address, "a", "", "address")
   flag.StringVar(&f.config.ClientId, "c", f.config.ClientId, "client ID")
   flag.Var(&f.filters, "f", net.FilterUsage)
   flag.StringVar(&f.config.PrivateKey, "p", f.config.PrivateKey, "private key")
   flag.StringVar(&plex.ForwardedFor, "x", "", "x-forwarded-for")
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

func (f *flag_set) do_address() error {
   data, err := plex.NewUser()
   if err != nil {
      return err
   }
   var user plex.User
   err = user.Unmarshal(data)
   if err != nil {
      return err
   }
   match, err := user.Match(plex.Path(f.address))
   if err != nil {
      return err
   }
   data1, err := user.Metadata(match)
   if err != nil {
      return err
   }
   var metadata plex.Metadata
   err = metadata.Unmarshal(data1)
   if err != nil {
      return err
   }
   part, ok := metadata.Dash()
   if !ok {
      return errors.New(".Dash()")
   }
   resp, err := user.Mpd(part)
   if err != nil {
      return err
   }
   f.config.Send = func(data []byte) ([]byte, error) {
      return user.Widevine(part, data)
   }
   return f.filters.Filter(resp, &f.config)
}
