package main

import (
   "41.neocities.org/media/plex"
   "41.neocities.org/net"
   "errors"
   "flag"
   "os"
   "path/filepath"
)

func (f *flag_set) New() error {
   var err error
   f.cache, err = os.UserCacheDir()
   if err != nil {
      return err
   }
   f.cache = filepath.ToSlash(f.cache)
   f.license.ClientId = f.cache + "/L3/client_id.bin"
   f.license.PrivateKey = f.cache + "/L3/private_key.pem"
   flag.StringVar(&f.license.ClientId, "c", f.license.ClientId, "client ID")
   flag.StringVar(&f.license.PrivateKey, "p", f.license.PrivateKey, "private key")
   flag.StringVar(&plex.ForwardedFor, "x", "", "x-forwarded-for")
   flag.StringVar(&f.address, "a", "", "address")
   flag.Var(&f.bitrate, "b", "bitrate")
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

type flag_set struct {
   cache   string
   license net.License
   address string
   bitrate net.Bitrate
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
   f.license.Widevine = func(data []byte) ([]byte, error) {
      return user.Widevine(part, data)
   }
   return f.license.Bitrate(resp, &f.bitrate)
}
