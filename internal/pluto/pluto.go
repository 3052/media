package main

import (
   "41.neocities.org/media/pluto"
   "41.neocities.org/net"
   "errors"
   "flag"
   "fmt"
   "log"
   "net/http"
   "net/url"
   "os"
   "path/filepath"
)

func main() {
   http.DefaultTransport = &http.Transport{
      Proxy: func(req *http.Request) (*url.URL, error) {
         log.Println(req.Method, req.URL)
         return nil, nil
      },
   }
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
   f.cdm.License = pluto.Widevine
   return f.filters.Filter(resp, &f.cdm)
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
   flag.StringVar(&f.cdm.ClientId, "c", f.cdm.ClientId, "client ID")
   flag.StringVar(
      &f.cdm.PrivateKey, "p", f.cdm.PrivateKey, "private key",
   )
   flag.StringVar(&pluto.ForwardedFor, "x", "", "x-forwarded-for")
   flag.StringVar(&f.show, "s", "", "show ID")
   flag.StringVar(&f.episode, "e", "", "episode/movie ID")
   flag.Var(&f.filters, "f", net.FilterUsage)
   flag.Parse()
   return nil
}

type flag_set struct {
   media   string
   cdm net.Cdm
   show string
   episode string
   filters net.Filters
}

func (f *flag_set) do_show() error {
   vod, err := pluto.NewVod(f.show)
   if err != nil {
      return err
   }
   fmt.Println(vod)
   return nil
}

