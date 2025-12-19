package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/pluto"
   "encoding/xml"
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
   "path"
   "path/filepath"
)

func (c *command) do_dash() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   c.config.Send = pluto.Widevine
   return c.config.Download(cache.Mpd.Url, cache.Mpd.Body, c.dash)
}

type user_cache struct {
   Mpd    *pluto.Mpd
   Series *pluto.Series
}

type command struct {
   config  maya.Config
   dash    string
   episode string
   movie   string
   name    string
   show    string
}

func main() {
   log.SetFlags(log.Ltime)
   maya.Transport(func(req *http.Request) string {
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

func (c *command) run() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   cache = filepath.ToSlash(cache)
   c.config.ClientId = cache + "/L3/client_id.bin"
   c.config.PrivateKey = cache + "/L3/private_key.pem"
   c.name = cache + "/pluto/userCache.xml"

   flag.StringVar(&c.config.ClientId, "c", c.config.ClientId, "client ID")
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.episode, "e", "", "episode ID")
   flag.StringVar(&c.movie, "m", "", "movie ID")
   flag.StringVar(&c.config.PrivateKey, "p", c.config.PrivateKey, "private key")
   flag.StringVar(&c.show, "s", "", "show ID")
   flag.IntVar(&c.config.Threads, "t", 2, "threads")
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
   var mpd pluto.Mpd
   err = mpd.Fetch(series.GetMovieURL())
   if err != nil {
      return err
   }
   err = write(c.name, &user_cache{Mpd: &mpd})
   if err != nil {
      return err
   }
   return maya.Representations(mpd.Url, mpd.Body)
}

func (c *command) do_episode() error {
   cache, err := read(c.name)
   if err != nil {
      return err
   }
   link, err := cache.Series.GetEpisodeURL(c.episode)
   if err != nil {
      return err
   }
   var mpd pluto.Mpd
   err = mpd.Fetch(link)
   if err != nil {
      return err
   }
   cache.Mpd = &mpd
   err = write(c.name, cache)
   if err != nil {
      return err
   }
   return maya.Representations(mpd.Url, mpd.Body)
}
