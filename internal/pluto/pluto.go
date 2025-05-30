package main

import (
   "41.neocities.org/media/pluto"
   "41.neocities.org/net"
   "errors"
   "flag"
   "fmt"
   "net/http"
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
   f.license.ClientId = f.media + "/client_id.bin"
   f.license.PrivateKey = f.media + "/private_key.pem"
   f.bitrate.Value = [][2]int{
      {100_000, 200_000}, {3_000_000, 5_000_000},
   }
   flag.StringVar(&f.license.ClientId, "c", f.license.ClientId, "client ID")
   flag.StringVar(
      &f.license.PrivateKey, "p", f.license.PrivateKey, "private key",
   )
   flag.IntVar(&net.Threads, "t", 1, "threads")
   flag.StringVar(&pluto.ForwardedFor, "x", "", "x-forwarded-for")
   ///////////////////////////////////////////////////////////////
   flag.StringVar(&f.show, "s", "", "show ID")
   ///////////////////////////////////////////////////////
   flag.StringVar(&f.episode, "e", "", "episode/movie ID")
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
   switch {
   case set.show != "":
      err = set.do_show()
   case set.episode != "":
      err = set.do_episode()
   default:
      flag.Usage()
   }
   if err != nil {
      panic(err)
   }
}

type flag_set struct {
   media   string
   license net.License
   ///////////////////
   show string
   //////////////
   episode string
   bitrate net.Bitrate
}

func (f *flag_set) do_episode() error {
   clips, err := pluto.NewClips(f.episode)
   if err != nil {
      return err
   }
   file, ok := clips.Dash()
   if !ok {
      return errors.New(".Dash()")
   }
   resp, err := file.Mpd()
   if err != nil {
      return err
   }
   f.license.Widevine = pluto.Widevine
   return f.license.Bitrate(resp, &f.bitrate)
}

func (f *flag_set) do_show() error {
   vod, err := pluto.NewVod(f.show)
   if err != nil {
      return err
   }
   fmt.Println(vod)
   return nil
}

