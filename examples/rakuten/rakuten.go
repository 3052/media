package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/rakuten"
   "flag"
   "fmt"
   "log"
)

func (c *client) do_dash_id() error {
   err := cache.Read(c)
   if err != nil {
      return err
   }
   stream, err := c.Content.Stream(
      c.Episode, c.Language, rakuten.Widevine, rakuten.Hd,
   )
   if err != nil {
      return err
   }
   job.Send = stream.Widevine
   return job.DownloadDash(c.Dash.Body, c.Dash.Url, c.dash_id)
}

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

func (c *client) do_address() error {
   c.Content = &rakuten.Content{}
   err := c.Content.Parse(c.address)
   if err != nil {
      return err
   }
   switch {
   case c.Content.IsMovie():
      movie, err := c.Content.Movie()
      if err != nil {
         return err
      }
      fmt.Println(movie)
   case c.Content.IsTvShow():
      show, err := c.Content.TvShow()
      if err != nil {
         return err
      }
      fmt.Println(show)
   }
   return cache.Write(c)
}

func (c *client) do_season() error {
   err := cache.Read(c)
   if err != nil {
      return err
   }
   season, err := c.Content.Season(c.season)
   if err != nil {
      return err
   }
   for i, episode := range season.Episodes {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(&episode)
   }
   return nil
}

func (c *client) do_language() error {
   err := cache.Update(c, func() error {
      stream, err := c.Content.Stream(
         c.Episode, c.Language, rakuten.Widevine, rakuten.Fhd,
      )
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
