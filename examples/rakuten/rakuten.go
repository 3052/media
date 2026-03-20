package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/rakuten"
   "flag"
   "fmt"
   "log"
)

func (c *client) do() error {
   err := cache.Setup("rosso/rakuten.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   // 1
   widevine := maya.StringVar(&c.Job.Widevine, "w", "Widevine")
   // 2
   address := maya.StringVar(&c.address, "a", "address")
   // 3
   season := maya.StringVar(&c.season, "s", "season ID")
   // 4
   language := maya.StringVar(&c.Language, "A", "audio language")
   episode := maya.StringVar(&c.Episode, "e", "episode ID")
   // 5
   dash_id := maya.StringVar(&c.dash_id, "d", "DASH ID")
   set := maya.Parse()
   switch {
   case set[widevine]:
      return cache.Write(c)
   case set[address]:
      return c.do_address()
   case set[season]:
      return with_cache(c.do_season)
   case set[language]:
      return with_cache(c.do_language)
   case set[dash_id]:
      return with_cache(c.do_dash_id)
   }
   return maya.Usage([][]*flag.Flag{
      {widevine},
      {address},
      {season},
      {language, episode},
      {dash_id},
   })
}
func (c *client) do_address() error {
   var err error
   c.Content, err = rakuten.ParseContent(c.address)
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
   stream, err := c.Content.Stream(
      c.Episode, c.Language, rakuten.Widevine, rakuten.Fhd,
   )
   if err != nil {
      return err
   }
   c.Dash, err = stream.Dash()
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}
func (c *client) do_dash_id() error {
   stream, err := c.Content.Stream(
      c.Episode, c.Language, rakuten.Widevine, rakuten.Hd,
   )
   if err != nil {
      return err
   }
   return c.Job.DownloadDash(c.Dash.Body, c.Dash.Url, c.dash_id, stream.Widevine)
}

func main() {
   log.SetFlags(log.Ltime)
   // server checks location on all requests
   maya.SetProxy("", "*.isma,*.ismv")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

var cache maya.Cache

type client struct {
   Content *rakuten.Content
   Dash    *rakuten.Dash
   // 1
   Job maya.Job
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
