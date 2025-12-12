package main

import (
   "41.neocities.org/media/pluto"
   "41.neocities.org/net"
   "encoding/xml"
   "flag"
   "errors"
   "fmt"
   "log"
   "net/http"
   "net/url"
   "os"
   "path"
   "path/filepath"
)

func (c *command) do_episode_movie() error {
   clip, err := pluto.NewClip(c.episode_movie)
   if err != nil {
      return err
   }
   file, ok := clip.Dash()
   if !ok {
      return errors.New(".Dash()")
   }
   var cache mpd
   cache.Url, cache.Body, err = file.Mpd()
   if err != nil {
      return err
   }
   data, err := xml.Marshal(cache)
   if err != nil {
      return err
   }
   log.Println("WriteFile", c.name)
   err = os.WriteFile(c.name, data, os.ModePerm)
   if err != nil {
      return err
   }
   return net.Representations(cache.Url, cache.Body)
}

type command struct {
   config net.Config
   name   string
   show string
   episode_movie string
   dash string
}
func (c *command) do_dash() error {
   data, err := os.ReadFile(c.name)
   if err != nil {
      return err
   }
   var cache mpd
   err = xml.Unmarshal(data, &cache)
   if err != nil {
      return err
   }
   c.config.Send = pluto.Widevine
   return c.config.Download(cache.Url, cache.Body, c.dash)
}

type mpd struct {
   Body []byte
   Url     *url.URL
}

func main() {
   log.SetFlags(log.Ltime)
   net.Transport(func(req *http.Request) string {
      if path.Ext(req.URL.Path) == ".m4s" {
         return ""
      }
      return "LP"
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
   c.name = cache + "/pluto/mpd.xml"

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
   series, err := pluto.NewSeries(c.show)
   if err != nil {
      return err
   }
   fmt.Println(series)
   return nil
}
