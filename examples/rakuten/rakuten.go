package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/rakuten"
   "flag"
   "fmt"
   "log"
)

func main() {
   log.SetFlags(log.Ltime)
   // everything needs proxy
   maya.SetProxy("", "*.isma,*.ismv")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

var cache maya.Cache

var job  maya.WidevineJob

func (c *client) do() error {
   job.ClientId, _ = maya.ResolveCache("L3/client_id.bin")
   job.PrivateKey, _ = maya.ResolveCache("L3/private_key.pem")
   err := c.cache.Setup("rosso/rakuten.xml")
   if err != nil {
      return err
   }
   // 1
   flag.StringVar(&c.address, "a", "", "address")
   // 2
   flag.StringVar(&c.season, "s", "", "season ID")
   // 3
   flag.StringVar(&c.Episode, "e", "", "episode ID")
   flag.StringVar(&c.Language, "A", "", "audio language")
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
      {"e", "A"},
      {"d", "c", "p"},
   })
}

type client struct {
   Dash     *rakuten.Dash
   Media    *rakuten.Media
   // 1
   address string
   // 2
   season string
   // 3
   Episode  string
   Language string
   // 4
   dash_id string
}

///

func (c *client) do_address() error {
   var media rakuten.Media
   err := media.ParseURL(c.address)
   if err != nil {
      return err
   }
   switch media.Type {
   case rakuten.MovieType:
      item, err := media.RequestMovie()
      if err != nil {
         return err
      }
      fmt.Println(item)
   case rakuten.TvShowType:
      item, err := media.RequestTvShow()
      if err != nil {
         return err
      }
      fmt.Println(item)
   }
   return c.cache.Write(saved_state{Media: &media})
}

func (c *client) do_language() error {
   err := c.cache.Read(&state)
   if err != nil {
      return err
   }
   var stream *rakuten.StreamData
   switch state.Media.Type {
   case rakuten.MovieType:
      stream, err = state.Media.MovieStream(
         c.Language, rakuten.Player.Widevine, rakuten.Quality.FHD,
      )
   case rakuten.TvShowType:
      stream, err = state.Media.EpisodeStream(
         c.Episode, c.Language, rakuten.Player.Widevine, rakuten.Quality.FHD,
      )
   }
   if err != nil {
      return err
   }
   state.Dash, err = stream.Dash()
   if err != nil {
      return err
   }
   err = c.cache.Write(state)
   if err != nil {
      return err
   }
   return maya.ListDash(state.Dash.Body, state.Dash.Url)
}

func (c *client) do_dash_id() error {
   err := c.cache.Read(&state)
   if err != nil {
      return err
   }
   var stream *rakuten.StreamData
   switch state.Media.Type {
   case rakuten.MovieType:
      stream, err = state.Media.MovieStream(
         state.Language,
         rakuten.Player.Widevine,
         rakuten.Quality.HD,
      )
   case rakuten.TvShowType:
      stream, err = state.Media.EpisodeStream(
         state.Episode,
         state.Language,
         rakuten.Player.Widevine,
         rakuten.Quality.HD,
      )
   }
   if err != nil {
      return err
   }
   job.Send = stream.Widevine
   return job.DownloadDash(state.Dash.Body, state.Dash.Url, c.dash_id)
}

func (c *client) do_season() error {
   err := c.cache.Read(&state)
   if err != nil {
      return err
   }
   season, err := state.Media.RequestSeason(c.season)
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
