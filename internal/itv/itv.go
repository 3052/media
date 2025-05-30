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

func (f *flag_set) New() error {
   var err error
   f.media, err = os.UserHomeDir()
   if err != nil {
      return err
   }
   f.media = filepath.ToSlash(f.media) + "/media"
   f.e.ClientId = f.media + "/client_id.bin"
   f.e.PrivateKey = f.media + "/private_key.pem"
   flag.StringVar(&f.address, "a", "", "address")
   flag.StringVar(&f.playlist, "b", "", "playlist URL")
   flag.StringVar(&f.e.ClientId, "c", f.e.ClientId, "client ID")
   flag.StringVar(&f.dash, "i", "", "DASH ID")
   flag.StringVar(&f.e.PrivateKey, "p", f.e.PrivateKey, "private key")
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
   case set.dash != "":
      err = set.do_dash()
   default:
      flag.Usage()
   }
   if err != nil {
      panic(err)
   }
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

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
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
   err = write_file(f.media+"/itv/Playlist", data)
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
   return net.Mpd(f.media+"/Mpd", resp)
}

func (f *flag_set) do_dash() error {
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
   return f.e.Download(f.media+"/Mpd", f.dash)
}

type flag_set struct {
   e        net.License
   media    string
   dash     string
   playlist string
   address  string
}

