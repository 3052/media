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

func (c *command) do_episode() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   address, err := cache.Series.GetEpisodeURL(c.episode)
   if err != nil {
      return err
   }
   cache.Mpd, cache.MpdBody, err = pluto.Mpd(address)
   if err != nil {
      return err
   }
   err = write(c.name, cache)
   if err != nil {
      return err
   }
   return net.Representations(cache.Mpd, cache.MpdBody)
}

type command struct {
   config  net.Config
   dash    string
   episode string
   movie   string
   name    string
   show    string
}

func (c *command) do_dash() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   c.config.Send = pluto.Widevine
   return c.config.Download(cache.Mpd, cache.MpdBody, c.dash)
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

func (c *command) do_movie() error {
   var series pluto.Series
   err := series.Fetch(c.movie)
   if err != nil {
      return err
   }
   var cache user_cache
   cache.Mpd, cache.MpdBody, err = pluto.Mpd(series.GetMovieURL())
   if err != nil {
      return err
   }
   err = write(c.name, &cache)
   if err != nil {
      return err
   }
   return net.Representations(cache.Mpd, cache.MpdBody)
}

func (c *command) do_show() error {
   var series pluto.Series
   err := series.Fetch(c.show)
   if err != nil {
      return err
   }
   fmt.Println(&series.Vod[0])
   return write(c.name, &user_cache{Series: &series})
}

func read(name string) (*user_cache, error) {
   data, err := os.ReadFile(name)
   if err != nil {
      return nil, err
   }
   cache := &user_cache{}
   err = xml.Unmarshal(data, cache)
   if err != nil {
      return nil, err
   }
   return cache, nil
}

func write(name string, cache *user_cache) error {
   data, err := xml.Marshal(cache)
   if err != nil {
      return err
   }
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

type user_cache struct {
   Mpd     *url.URL
   MpdBody []byte
   Series  *pluto.Series
}
