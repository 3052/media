package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/pluto"
   "flag"
   "fmt"
   "log"
   "net/http"
   "path"
)

func (c *command) do_dash() error {
   var dash pluto.Dash
   err := c.cache.Get("Dash", &dash)
   if err != nil {
      return err
   }
   c.job.Send = pluto.Widevine
   return c.job.DownloadDash(dash.Body, dash.Url, c.dash)
}

func main() {
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

type command struct {
   cache maya.Cache
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

func (c *command) run() error {
   c.cache.Init("L3")
   c.job.ClientId = c.cache.Join("client_id.bin")
   c.job.PrivateKey = c.cache.Join("private_key.pem")
   c.cache.Init("pluto")
   // 1
   flag.StringVar(&c.movie, "m", "", "movie URL")
   flag.StringVar(&c.proxy, "x", "", "proxy")
   // 2
   flag.StringVar(&c.show, "s", "", "show URL")
   // 3
   flag.StringVar(&c.episode, "e", "", "episode ID")
   // 4
   flag.StringVar(&c.dash, "d", "", "DASH ID")
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
   err = c.cache.Set("Dash", dash)
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
   return c.cache.Set("Series", series)
}

func (c *command) do_episode() error {
   var series pluto.Series
   err := c.cache.Get("Series", &series)
   if err != nil {
      return err
   }
   link, err := series.GetEpisodeURL(c.episode)
   if err != nil {
      return err
   }
   var dash pluto.Dash
   err = dash.Fetch(link)
   if err != nil {
      return err
   }
   err = c.cache.Set("Dash", dash)
   if err != nil {
      return err
   }
   return maya.ListDash(dash.Body, dash.Url)
}
