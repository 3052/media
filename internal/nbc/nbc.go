package main

import (
   "41.neocities.org/media/nbc"
   "41.neocities.org/net"
   "flag"
   "net/http"
   "os"
   "path/filepath"
)

type flag_set struct {
   nbc   int
   cdm     net.License
   cache string
   dash  string
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
   flag.StringVar(&f.cdm.ClientId, "c", f.cdm.ClientId, "client ID")
   flag.StringVar(&f.dash, "d", "", "dash ID")
   flag.IntVar(&f.nbc, "n", 0, "NBC ID")
   flag.StringVar(&f.cdm.PrivateKey, "p", f.cdm.PrivateKey, "private key")
   flag.IntVar(&net.Threads, "t", 2, "threads")
   flag.Parse()
   return nil
}

func main() {
   var set flag_set
   err := set.New()
   if err != nil {
      panic(err)
   }
   switch {
   case set.nbc >= 1:
      err = set.do_nbc()
   case set.dash != "":
      err = set.do_dash()
   default:
      flag.Usage()
   }
   if err != nil {
      panic(err)
   }
}

func (f *flag_set) do_nbc() error {
   var metadata nbc.Metadata
   err := metadata.New(f.nbc)
   if err != nil {
      return err
   }
   vod, err := metadata.Vod()
   if err != nil {
      return err
   }
   resp, err := http.Get(vod.PlaybackUrl)
   if err != nil {
      return err
   }
   return net.Mpd(f.cache+"/Mpd", resp)
}

func (f *flag_set) do_dash() error {
   f.cdm.Widevine = nbc.Widevine
   return f.cdm.Download(f.cache+"/Mpd", f.dash)
}
