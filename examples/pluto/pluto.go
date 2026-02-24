package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/pluto"
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
   c.job.ClientId = cache + "/L3/client_id.bin"
   c.job.PrivateKey = cache + "/L3/private_key.pem"
   c.name = cache + "/rosso/pluto.xml"
   // 1
   flag.StringVar(&c.movie, "m", "", "movie URL")
   flag.StringVar(&c.proxy, "x", "", "proxy")
   // 2
   flag.StringVar(&c.show, "s", "", "show URL")
   // 3
   flag.StringVar(&c.episode, "e", "", "episode ID")
   // 4
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.IntVar(&c.job.Threads, "t", 2, "threads")
   flag.StringVar(&c.job.ClientId, "c", c.job.ClientId, "client ID")
   flag.StringVar(&c.job.PrivateKey, "p", c.job.PrivateKey, "private key")
   flag.Parse()
   maya.SetProxy(func(req *http.Request) (string, bool) {
      if path.Ext(req.URL.Path) == ".m4s" {
         return "", false
      }
      return c.proxy, true
   })
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
      {"m", "x"},
      {"s"},
      {"e"},
      {"d", "c", "p"},
   })
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

type user_cache struct {
   Dash   *pluto.Dash
   Series *pluto.Series
}

func main() {
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

type command struct {
   name string
   // 1
   movie string
   proxy string
   // 2
   show string
   // 3
   episode string
   // 4
   dash string
   job  maya.WidevineJob
}
