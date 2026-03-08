// everything needs proxy
package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/rakuten"
   "flag"
   "fmt"
   "log"
)

func (c *client) do_season() error {
   err := cache.Read(c)
   if err != nil {
      return err
   }
   season, err := c.Media.RequestSeason(c.season)
   if err != nil {
      return err
   }
   for i, item := range season.Episodes {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(&item)
   }
   return nil
}

type client struct {
   Dash  *rakuten.Dash
   Media *rakuten.Media
   // 1
   proxy string
   // 2
   address string
   // 3
   season string
   // 4
   Language string
   Episode  string
   // 5
   dash_id string
}

func (c *client) do() error {
   job.ClientId, _ = maya.ResolveCache("L3/client_id.bin")
   job.PrivateKey, _ = maya.ResolveCache("L3/private_key.pem")
   err := cache.Setup("rosso/rakuten.xml")
   if err != nil {
      return err
   }
   // 1
   flag.StringVar(
      &c.proxy, "x", "", "proxy (server checks location on all requests)",
   )
   // 2
   flag.StringVar(&c.address, "a", "", "address")
   // 3
   flag.StringVar(&c.season, "s", "", "season ID")
   // 4
   flag.StringVar(&c.Language, "A", "", "audio language")
   flag.StringVar(&c.Episode, "e", "", "episode ID")
   // 5
   flag.StringVar(&c.dash_id, "d", "", "DASH ID")
   flag.StringVar(&job.ClientId, "c", job.ClientId, "client ID")
   flag.StringVar(&job.PrivateKey, "p", job.PrivateKey, "private key")
   flag.Parse()
   err = maya.SetProxy(c.proxy, "*.isma,*.ismv")
   if err != nil {
      return err
   }
   if c.address != "" {
      return c.do_address()
   }
   if c.season != "" {
      return c.do_season()
   }
   if c.Language != "" {
      return c.do_language()
   }
   if c.dash_id != "" {
      return c.do_dash_id()
   }
   return maya.Usage([][]string{
      {"x"},
      {"a"},
      {"s"},
      {"A", "e"},
      {"d", "c", "p"},
   })
}
var cache maya.Cache

var job maya.WidevineJob

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

///

func (c *client) do_address() error {
   var err error
   c.Media, err = rakuten.ParseMedia(c.address)
   if err != nil {
      return err
   }
   switch c.Media.Type {
   case rakuten.MovieType:
      item, err := c.Media.RequestMovie()
      if err != nil {
         return err
      }
      fmt.Println(item)
   case rakuten.TvShowType:
      item, err := c.Media.RequestTvShow()
      if err != nil {
         return err
      }
      fmt.Println(item)
   }
   return cache.Write(c)
}

func (c *client) do_language() error {
   err := cache.Update(c, func() error {
      var (
         stream *rakuten.StreamData
         err    error
      )
      switch c.Media.Type {
      case rakuten.MovieType:
         stream, err = c.Media.MovieStream(
            c.Language, rakuten.Player.Widevine, rakuten.Quality.FHD,
         )
      case rakuten.TvShowType:
         stream, err = c.Media.EpisodeStream(
            c.Episode, c.Language, rakuten.Player.Widevine, rakuten.Quality.FHD,
         )
      }
      if err != nil {
         return err
      }
      c.Dash, err = stream.Dash()
      return err
   })
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}

func (c *client) do_dash_id() error {
   err := cache.Read(c)
   if err != nil {
      return err
   }
   var stream *rakuten.StreamData
   switch c.Media.Type {
   case rakuten.MovieType:
      stream, err = c.Media.MovieStream(
         c.Language,
         rakuten.Player.Widevine,
         rakuten.Quality.HD,
      )
   case rakuten.TvShowType:
      stream, err = c.Media.EpisodeStream(
         c.Episode,
         c.Language,
         rakuten.Player.Widevine,
         rakuten.Quality.HD,
      )
   }
   if err != nil {
      return err
   }
   job.Send = stream.Widevine
   return job.DownloadDash(c.Dash.Body, c.Dash.Url, c.dash_id)
}
