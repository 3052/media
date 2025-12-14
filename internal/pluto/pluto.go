package main

import (
   "41.neocities.org/media/pluto"
   "41.neocities.org/net"
   "encoding/xml"
   "flag"
   "fmt"
   "log"
   "net/http"
   "net/url"
   "os"
   "path"
   "path/filepath"
)

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
   c.name = cache + "/pluto/user_cache.xml"

   flag.StringVar(&c.config.ClientId, "c", c.config.ClientId, "client ID")
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.episode, "e", "", "episode ID")
   flag.StringVar(&c.movie, "m", "", "movie ID")
   flag.StringVar(&c.config.PrivateKey, "p", c.config.PrivateKey, "private key")
   flag.StringVar(&c.show, "s", "", "show ID")
   flag.Parse()

   if c.movie != "" {
      return c.do_movie()
   }
   if c.show != "" {
      return c.do_show()
   }
   if c.episode != "" {
      return c.do_episode()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   flag.Usage()
   return nil
}

type user_cache struct {
   Mpd *url.URL
   MpdBody []byte
}

type command struct {
   config        net.Config
   name          string
   // 1
   movie string
   // 2
   show          string
   // 3
   episode string
   // 4
   dash          string
}

///

func (c *command) do_movie() error {
   var series pluto.Series
   err := series.Fetch(c.episode_movie)
   if err != nil {
      return err
   }
   var cache user_cache
   
   cache.Mpd, cache.MpdBody, err = series.Mpd()
   
   cache.Url, cache.Body, err = series.Mpd()
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

func (c *command) do_episode() error {
   var series pluto.Series
   err := series.Fetch(c.episode_movie)
   if err != nil {
      return err
   }
   var cache mpd
   cache.Url, cache.Body, err = series.Mpd()
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

func (c *command) do_show() error {
   var series pluto.Series
   err := series.Fetch(c.show)
   if err != nil {
      return err
   }
   fmt.Println(&series.Vod[0])
   return nil
}
