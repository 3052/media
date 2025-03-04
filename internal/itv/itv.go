package main

import (
   "41.neocities.org/media/internal"
   "41.neocities.org/media/itv"
   "errors"
   "flag"
   "log"
   "os"
   "path"
   "path/filepath"
)

type flags struct {
   address string
   dash    string
   e       internal.License
   media   string
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
   flag.StringVar(&f.address, "a", "", "address")
   flag.StringVar(&f.e.ClientId, "c", f.e.ClientId, "client ID")
   flag.StringVar(&f.dash, "i", "", "DASH ID")
   flag.StringVar(&f.e.PrivateKey, "p", f.e.PrivateKey, "private key")
   flag.Parse()
   if f.address != "" {
      err := f.download()
      if err != nil {
         panic(err)
      }
   } else {
      flag.Usage()
   }
}

func (f *flags) download() error {
   if f.dash != "" {
      data, err := os.ReadFile(f.media + "/itv/Playlist")
      if err != nil {
         return err
      }
      var play itv.Playlist
      err = play.Unmarshal(data)
      if err != nil {
         return err
      }
      file, _ := play.FullHd()
      f.e.Widevine = func(data []byte) ([]byte, error) {
         return file.Widevine(data)
      }
      return f.e.Download(f.media + "/Mpd", f.dash)
   }
   var id itv.EpisodeId
   err := id.Set(path.Base(f.address))
   if err != nil {
      return err
   }
   data, err := id.Playlist()
   if err != nil {
      return err
   }
   var play itv.Playlist
   err = play.Unmarshal(data)
   if err != nil {
      return err
   }
   log.Println("WriteFile", f.media + "/itv/Playlist")
   err = os.WriteFile(f.media + "/itv/Playlist", data, os.ModePerm)
   if err != nil {
      return err
   }
   file, ok := play.FullHd()
   if !ok {
      return errors.New(".FullHd()")
   }
   resp, err := file.Mpd()
   if err != nil {
      return err
   }
   return internal.Mpd(f.media + "/Mpd", resp)
}
