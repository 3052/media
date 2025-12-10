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

func main() {
   log.SetFlags(log.Ltime)
   net.Transport(func(req *http.Request) string {
      if path.Ext(req.URL.Path) == ".m4s" {
         return ""
      }
      return "L"
   })
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

type command struct {
   config  net.Config
   name   string
   // 1
   show    string
   // 2
   episode string
   // 3
   dash string
}

///

func (c *command) run() error {
   c.cache, err = os.UserCacheDir()
   if err != nil {
      return err
   }
   c.cache = filepath.ToSlash(c.cache)
   c.config.ClientId = c.cache + "/L3/client_id.bin"
   c.config.PrivateKey = c.cache + "/L3/private_key.pem"
   flag.StringVar(&c.config.ClientId, "c", c.config.ClientId, "client ID")
   flag.StringVar(&c.episode, "e", "", "episode/movie ID")
   flag.Var(&c.filters, "f", net.FilterUsage)
   flag.StringVar(&c.config.PrivateKey, "p", c.config.PrivateKey, "private key")
   flag.StringVar(&c.show, "s", "", "show ID")
   flag.Parse()
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
   var err error
   return nil
}

func (c *command) do_episode() error {
   clips, err := pluto.NewClips(c.episode)
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
   c.config.Send = pluto.Widevine
   return c.filters.Filter(resp, &c.config)
}

func (c *command) do_show() error {
   vod, err := pluto.NewVod(c.show)
   if err != nil {
      return err
   }
   fmt.Println(vod)
   return nil
}
