package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/pluto"
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
   "path"
   "path/filepath"
)

func (c *command) run() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   cache = filepath.ToSlash(cache)
   c.name = cache + "/pluto/userCache.xml"
   c.job.ClientId = cache + "/L3/client_id.bin"
   c.job.PrivateKey = cache + "/L3/private_key.pem"
   // 1
   flag.StringVar(&c.movie, "m", "", "movie URL")
   // 2
   flag.StringVar(&c.show, "s", "", "show URL")
   // 3
   flag.StringVar(&c.episode, "e", "", "episode ID")
   // 4
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.job.ClientId, "c", c.job.ClientId, "client ID")
   flag.StringVar(&c.job.PrivateKey, "p", c.job.PrivateKey, "private key")
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
   return maya.Usage([][]string{
      {"m"},
      {"s"},
      {"e"},
      {"d", "c", "p"},
   })
}

type command struct {
   name string
   // 1
   movie string
   // 2
   show string
   // 3
   episode string
   // 4
   dash string
   job  maya.WidevineJob
}

func (c *command) do_episode() error {
   cache, err := maya.Read[user_cache](c.name)
   if err != nil {
      return err
   }
   link, err := cache.Series.GetEpisodeURL(c.episode)
   if err != nil {
      return err
   }
   var dash pluto.Dash
   err = dash.Fetch(link)
   if err != nil {
      return err
   }
   cache.Dash = &dash
   err = maya.Write(c.name, cache)
   if err != nil {
      return err
   }
   return maya.ListDash(dash.Body, dash.Url)
}

func (c *command) do_dash() error {
   cache, err := maya.Read[user_cache](c.name)
   if err != nil {
      return err
   }
   c.job.Send = pluto.Widevine
   return c.job.DownloadDash(cache.Dash.Body, cache.Dash.Url, c.dash)
}

func (c *command) do_movie() error {
   var series pluto.Series
   err := series.Fetch(path.Base(c.movie))
   if err != nil {
      return err
   }
   var dash pluto.Dash
   err = dash.Fetch(series.GetMovieURL())
   if err != nil {
      return err
   }
   err = maya.Write(c.name, &user_cache{Dash: &dash})
   if err != nil {
      return err
   }
   return maya.ListDash(dash.Body, dash.Url)
}

func (c *command) do_show() error {
   var series pluto.Series
   err := series.Fetch(path.Base(c.show))
   if err != nil {
      return err
   }
   fmt.Println(&series.Vod[0])
   return maya.Write(c.name, &user_cache{Series: &series})
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

type user_cache struct {
   Dash   *pluto.Dash
   Series *pluto.Series
}
