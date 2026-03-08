package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/pluto"
   "flag"
   "fmt"
   "log"
   "path"
)

func (c *client) do_dash_id() error {
   err := cache.Read(c)
   if err != nil {
      return err
   }
   job.Send = pluto.Widevine
   return job.DownloadDash(c.Dash.Body, c.Dash.Url, c.dash_id)
}

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.m4s")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

var cache maya.Cache

var job maya.WidevineJob

func (c *client) do() error {
   job.ClientId, _ = maya.ResolveCache("L3/client_id.bin")
   job.PrivateKey, _ = maya.ResolveCache("L3/private_key.pem")
   err := cache.Setup("rosso/pluto.xml")
   if err != nil {
      return err
   }
   // 1
   flag.StringVar(&c.movie, "m", "", "movie URL")
   // 2
   flag.StringVar(&c.show, "s", "", "show URL")
   // 3
   flag.StringVar(&c.episode, "e", "", "episode ID")
   // 4
   flag.StringVar(&c.dash_id, "d", "", "DASH ID")
   flag.StringVar(&job.ClientId, "c", job.ClientId, "client ID")
   flag.StringVar(&job.PrivateKey, "p", job.PrivateKey, "private key")
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
   if c.dash_id != "" {
      return c.do_dash_id()
   }
   return maya.Usage([][]string{
      {"m"},
      {"s"},
      {"e"},
      {"d", "c", "p"},
   })
}

func (c *client) do_movie() error {
   var series pluto.Series
   err := series.Fetch(path.Base(c.movie))
   if err != nil {
      return err
   }
   c.Dash = &pluto.Dash{}
   err = c.Dash.Fetch(series.GetMovieUrl())
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}

func (c *client) do_show() error {
   c.Series = &pluto.Series{}
   err := c.Series.Fetch(path.Base(c.show))
   if err != nil {
      return err
   }
   fmt.Println(&c.Series.Vod[0])
   return cache.Write(c)
}

type client struct {
   Dash   *pluto.Dash
   Series *pluto.Series
   // 1
   movie string
   // 2
   show string
   // 3
   episode string
   // 4
   dash_id string
}

func (c *client) do_episode() error {
   err := cache.Read(c)
   if err != nil {
      return err
   }
   url, err := c.Series.GetEpisodeUrl(c.episode)
   if err != nil {
      return err
   }
   c.Dash = &pluto.Dash{}
   err = c.Dash.Fetch(url)
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}
