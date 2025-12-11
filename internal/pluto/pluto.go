package main

import (
   "41.neocities.org/media/pluto"
   "41.neocities.org/net"
   "encoding/json"
   "flag"
   "fmt"
   "log"
   "net/http"
   "net/url"
   "os"
   "path"
   "path/filepath"
)

func (c *command) do_dash() error {
   data, err := os.ReadFile(c.name)
   if err != nil {
      return err
   }
   var cache user_cache
   err = json.Unmarshal(data, &cache)
   if err != nil {
      return err
   }
   c.config.Send = pluto.Widevine
   return c.config.Download(cache.Mpd, cache.MpdBody, c.dash)
}

type user_cache struct {
   Mpd     *url.URL
   MpdBody []byte
}

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

func (c *command) run() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   cache = filepath.ToSlash(cache)
   c.config.ClientId = cache + "/L3/client_id.bin"
   c.config.PrivateKey = cache + "/L3/private_key.pem"
   c.name = cache + "/pluto/mpd.json"

   flag.StringVar(&c.config.ClientId, "c", c.config.ClientId, "client ID")
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.episode_movie, "e", "", "episode/movie ID")
   flag.StringVar(&c.config.PrivateKey, "p", c.config.PrivateKey, "private key")
   flag.StringVar(&c.show, "s", "", "show ID")
   flag.Parse()

   if c.show != "" {
      return c.do_show()
   }
   if c.episode_movie != "" {
      return c.do_episode_movie()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   flag.Usage()
   return nil
}

func (c *command) do_show() error {
   var series pluto.Series
   err := series.Fetch(c.show)
   if err != nil {
      return err
   }
   fmt.Println(&series) // FIXME
   return nil
}

func (c *command) do_episode_movie() error {
   var series pluto.Series
   err := series.Fetch(c.episode_movie)
   if err != nil {
      return err
   }
   var cache user_cache
   cache.Mpd, cache.MpdBody, err = series.Mpd()
   if err != nil {
      return err
   }
   data, err := json.Marshal(cache)
   if err != nil {
      return err
   }
   log.Println("WriteFile", c.name)
   err = os.WriteFile(c.name, data, os.ModePerm)
   if err != nil {
      return err
   }
   return net.Representations(cache.Mpd, cache.MpdBody)
}

type command struct {
   config net.Config
   name   string
   // 1
   show string
   // 2
   episode_movie string
   // 3
   dash string
}
