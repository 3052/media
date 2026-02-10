package main

import (
   "41.neocities.org/media/internal"
   "41.neocities.org/media/nbc"
   "flag"
   "net/http"
   "os"
   "path/filepath"
)

func (f *flags) download() error {
   if f.representation != "" {
      f.e.Widevine = nbc.Widevine
      return f.e.Download(f.media + "/Mpd", f.representation)
   }
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
   return internal.Mpd(f.media + "/Mpd", resp)
}

type flags struct {
   e              internal.License
   media           string
   nbc            int
   representation string
}

func (f *flags) New() error {
   var err error
   f.media, err = os.UserHomeDir()
   if err != nil {
      return err
   }
   f.media = filepath.ToSlash(f.media) + "/media"
   f.e.ClientId = f.media + "/client_id.bin"
   f.e.PrivateKey = f.media + "/private_key.pem"
   return nil
}

func main() {
   var f flags
   err := f.New()
   if err != nil {
      panic(err)
   }
   flag.StringVar(&f.e.ClientId, "c", f.e.ClientId, "client ID")
   flag.StringVar(&f.e.PrivateKey, "p", f.e.PrivateKey, "private key")
   flag.StringVar(&f.representation, "i", "", "representation")
   flag.IntVar(&f.nbc, "b", 0, "NBC ID")
   flag.Parse()
   if f.nbc >= 1 {
      err := f.download()
      if err != nil {
         panic(err)
      }
   } else {
      flag.Usage()
   }
}
