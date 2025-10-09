package main

import (
   "41.neocities.org/media/itv"
   "41.neocities.org/net"
   "errors"
   "flag"
   "fmt"
   "log"
   "os"
   "path/filepath"
)

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

type flag_set struct {
   address  string
   cdm      net.Cdm
   filters  net.Filters
   media    string
   playlist string
}

func (f *flag_set) do_address() error {
   var id itv.LegacyId
   err := id.Set(f.address)
   if err != nil {
      return err
   }
   titles, err := id.Titles()
   if err != nil {
      return err
   }
   for i, title := range titles {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(&title)
   }
   return nil
}

func (f *flag_set) New() error {
   var err error
   f.media, err = os.UserHomeDir()
   if err != nil {
      return err
   }
   f.media = filepath.ToSlash(f.media) + "/media"
   f.cdm.ClientId = f.media + "/client_id.bin"
   f.cdm.PrivateKey = f.media + "/private_key.pem"
   flag.StringVar(&f.cdm.ClientId, "C", f.cdm.ClientId, "client ID")
   flag.StringVar(&f.cdm.PrivateKey, "P", f.cdm.PrivateKey, "private key")
   flag.StringVar(&f.address, "a", "", "address")
   flag.Var(&f.filters, "f", net.FilterUsage)
   flag.StringVar(&f.playlist, "p", "", "playlist URL")
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
   case set.address != "":
      err = set.do_address()
   case set.playlist != "":
      err = set.do_playlist()
   default:
      flag.Usage()
   }
   if err != nil {
      panic(err)
   }
}

func (f *flag_set) do_playlist() error {
   var title itv.Title
   title.LatestAvailableVersion.PlaylistUrl = f.playlist
   data, err := title.Playlist()
   if err != nil {
      return err
   }
   var play itv.Playlist
   err = play.Unmarshal(data)
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
   f.cdm.License = func(data []byte) ([]byte, error) {
      return file.Widevine(data)
   }
   return f.filters.Filter(resp, &f.cdm)
}
