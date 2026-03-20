package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/pluto"
   "flag"
   "fmt"
   "log"
   "path"
)

func (c *client) do() error {
   err := cache.Setup("rosso/pluto.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   widevine := maya.StringVar(&c.Job.Widevine, "w", "Widevine")
   //----------------------------------------------------------
   movie := maya.StringVar(&c.movie, "m", "movie URL")
   //----------------------------------------------------------
   show := maya.StringVar(&c.show, "s", "show URL")
   //----------------------------------------------------------
   episode := maya.StringVar(&c.episode, "e", "episode ID")
   //----------------------------------------------------------
   dash_id := maya.StringVar(&c.dash_id, "d", "DASH ID")
   set := maya.Parse()
   switch {
   case set[widevine]:
      return cache.Write(c)
   case set[movie]:
      return c.do_movie()
   case set[show]:
      return c.do_show()
   case set[episode]:
      return with_cache(c.do_episode)
   case set[dash_id]:
      return with_cache(c.do_dash_id)
   }
   return maya.Usage([][]*flag.Flag{
      {widevine},
      {movie},
      {show},
      {episode},
      {dash_id},
   })
}

func (c *client) do_dash_id() error {
   return c.Job.DownloadDash(c.Dash.Body, c.Dash.Url, c.dash_id, pluto.Widevine)
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

func (c *client) do_movie() error {
   series, err := pluto.FetchSeries(path.Base(c.movie))
   if err != nil {
      return err
   }
   c.Dash, err = pluto.FetchDash(series.GetMovieUrl())
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
   var err error
   c.Series, err = pluto.FetchSeries(path.Base(c.show))
   if err != nil {
      return err
   }
   fmt.Println(&c.Series.Vod[0])
   return cache.Write(c)
}

type client struct {
   Dash   *pluto.Dash
   Series *pluto.Series
   //------------------
   Job maya.Job
   //------------------
   movie string
   //------------------
   show string
   //------------------
   episode string
   //------------------
   dash_id string
}

func (c *client) do_episode() error {
   url, err := c.Series.GetEpisodeUrl(c.episode)
   if err != nil {
      return err
   }
   c.Dash, err = pluto.FetchDash(url)
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}
