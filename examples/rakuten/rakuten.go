package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/rakuten"
   "flag"
   "fmt"
   "log"
)

var cache maya.Cache

var job maya.WidevineJob

func main() {
   log.SetFlags(log.Ltime)
   // server checks location on all requests
   maya.SetProxy("", "*.isma,*.ismv")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do() error {
   job.ClientId, _ = maya.ResolveCache("L3/client_id.bin")
   job.PrivateKey, _ = maya.ResolveCache("L3/private_key.pem")
   err := cache.Setup("rosso/rakuten.xml")
   if err != nil {
      return err
   }
   // 1
   flag.StringVar(&c.address, "a", "", "address")
   // 2
   flag.StringVar(&c.season, "s", "", "season ID")
   // 3
   flag.StringVar(&c.Language, "A", "", "audio language")
   flag.StringVar(&c.Episode, "e", "", "episode ID")
   // 4
   flag.StringVar(&c.dash_id, "d", "", "DASH ID")
   flag.StringVar(&job.ClientId, "c", job.ClientId, "client ID")
   flag.StringVar(&job.PrivateKey, "p", job.PrivateKey, "private key")
   flag.Parse()
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
      {"a"},
      {"s"},
      {"A", "e"},
      {"d", "c", "p"},
   })
}

type client struct {
   Content *rakuten.Content
   Dash  *rakuten.Dash
   // 1
   address string
   // 2
   season string
   // 3
   Language string
   Episode  string
   // 4
   dash_id string
}

///

func (c *client) do_address() error {
   var err error
   c.Content, err = rakuten.ParseMedia(c.address)
   if err != nil {
      return err
   }
   switch c.Content.Type {
   case rakuten.MovieType:
      item, err := c.Content.RequestMovie()
      if err != nil {
         return err
      }
      fmt.Println(item)
   case rakuten.TvShowType:
      item, err := c.Content.RequestTvShow()
      if err != nil {
         return err
      }
      fmt.Println(item)
   }
   return cache.Write(c)
}

func (c *client) do_season() error {
   err := cache.Read(c)
   if err != nil {
      return err
   }
   season, err := c.Content.RequestSeason(c.season)
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

func (c *client) do_language() error {
   err := cache.Update(c, func() error {
      var (
         stream *rakuten.StreamData
         err    error
      )
      switch c.Content.Type {
      case rakuten.MovieType:
         stream, err = c.Content.MovieStream(
            c.Language, rakuten.Player.Widevine, rakuten.Quality.FHD,
         )
      case rakuten.TvShowType:
         stream, err = c.Content.EpisodeStream(
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
   switch c.Content.Type {
   case rakuten.MovieType:
      stream, err = c.Content.MovieStream(
         c.Language,
         rakuten.Player.Widevine,
         rakuten.Quality.HD,
      )
   case rakuten.TvShowType:
      stream, err = c.Content.EpisodeStream(
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
