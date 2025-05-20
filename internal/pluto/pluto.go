package main

import (
   "41.neocities.org/media/pluto"
   "41.neocities.org/net"
   "errors"
   "flag"
   "fmt"
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
   f.e.ClientId = f.media + "/client_id.bin"
   f.e.PrivateKey = f.media + "/private_key.pem"
   f.bitrate.Value = [][2]int{
      {100_000, 200_000}, {3_000_000, 5_000_000},
   }
   return nil
}

func (f *flags) do_show() error {
   vod, err := pluto.NewVod(f.show)
   if err != nil {
      return err
   }
   fmt.Println(vod)
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
   flag.IntVar(&net.Threads, "t", 1, "threads")
   flag.StringVar(&pluto.ForwardedFor, "x", "", "x-forwarded-for")
   /////////////////////////////////////////////////////////////////////
   flag.StringVar(&f.show, "s", "", "show ID")
   /////////////////////////////////////////////////////////////////////
   flag.StringVar(&f.episode, "e", "", "episode/movie ID")
   flag.Var(&f.bitrate, "b", "bitrate")
   flag.Parse()
   switch {
   case f.show != "":
      err = f.do_show()
   case f.episode != "":
      err = f.do_episode()
   default:
      flag.Usage()
   }
   if err != nil {
      panic(err)
   }
}

type flags struct {
   media string
   e     net.License
   ////////////////////////////
   show string
   //////////////////////
   episode string
   bitrate net.Bitrate
}

func (f *flags) do_episode() error {
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
   f.e.Widevine = pluto.Widevine
   return f.e.Bitrate(resp, &f.bitrate)
}
