package main

import (
   "41.neocities.org/media/pluto"
   "41.neocities.org/net"
   "errors"
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

var Transport = http.Transport{
   Proxy: func(req *http.Request) (*url.URL, error) {
      if path.Ext(req.URL.Path) != ".m4s" {
         log.Println(req.Method, req.URL)
      }
      return http.ProxyFromEnvironment(req)
   },
}

func main() {
   http.DefaultTransport = &pluto.Transport
   log.SetFlags(log.Ltime)
   var set flag_set
   err := set.New()
   if err != nil {
      log.Fatal(err)
   }
   switch {
   case set.episode != "":
      err = set.do_episode()
   case set.show != "":
      err = set.do_show()
   default:
      flag.Usage()
   }
   if err != nil {
      log.Fatal(err)
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
   f.config.Send = pluto.Widevine
   return f.filters.Filter(resp, &f.config)
}

type flag_set struct {
   cache   string
   config  net.Config
   episode string
   filters net.Filters
   show    string
}

func (f *flag_set) New() error {
   var err error
   f.cache, err = os.UserCacheDir()
   if err != nil {
      return err
   }
   f.cache = filepath.ToSlash(f.cache)
   f.config.ClientId = f.cache + "/L3/client_id.bin"
   f.config.PrivateKey = f.cache + "/L3/private_key.pem"
   flag.StringVar(&f.config.ClientId, "c", f.config.ClientId, "client ID")
   flag.StringVar(&f.episode, "e", "", "episode/movie ID")
   flag.Var(&f.filters, "f", net.FilterUsage)
   flag.StringVar(&f.config.PrivateKey, "p", f.config.PrivateKey, "private key")
   flag.StringVar(&f.show, "s", "", "show ID")
   flag.Parse()
   return nil
}
func (f *flag_set) do_show() error {
   vod, err := pluto.NewVod(f.show)
   if err != nil {
      return err
   }
   fmt.Println(vod)
   return nil
}
