package main

import (
   "41.neocities.org/media/plex"
   "41.neocities.org/net"
   "errors"
   "flag"
   "os"
   "path/filepath"
)

func (f *flags) New() error {
   var err error
   f.media, err = os.UserHomeDir()
   if err != nil {
      return err
   }
   f.media = filepath.ToSlash(f.media) + "/media"
   f.license.ClientId = f.media + "/client_id.bin"
   f.license.PrivateKey = f.media + "/private_key.pem"
   f.bitrate.Value = [][2]int{
      {128_000, 256_000}, {3_000_000, 4_000_000},
   }
   return nil
}

func main() {
   var f flags
   err := f.New()
   if err != nil {
      panic(err)
   }
   flag.StringVar(&f.license.ClientId, "c", f.license.ClientId, "client ID")
   flag.StringVar(
      &f.license.PrivateKey, "p", f.license.PrivateKey, "private key",
   )
   flag.StringVar(&plex.ForwardedFor, "x", "", "x-forwarded-for")
   /////////////////////////////////////////////////////////////////////////
   flag.StringVar(&f.address, "a", "", "address")
   flag.Var(&f.bitrate, "b", "bitrate")
   flag.Parse()
   if f.address != "" {
      err = f.do_address()
      if err != nil {
         panic(err)
      }
   } else {
      flag.Usage()
   }
}

type flags struct {
   media   string
   license net.License
   ////////////////////////////////
   address string
   bitrate net.Bitrate
}

func (f *flags) do_address() error {
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
