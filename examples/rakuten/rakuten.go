package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/rakuten"
   "flag"
   "fmt"
   "log"
   "net/http"
   "path"
)

func (c *command) do_dash() error {
   var dash rakuten.Dash
   err := c.cache.Get("Dash", &dash)
   if err != nil {
      return err
   }
   var language string
   err = c.cache.Get("Language", &language)
   if err != nil {
      return err
   }
   var media rakuten.Media
   err = c.cache.Get("Media", &media)
   if err != nil {
      return err
   }
   var stream *rakuten.StreamData
   switch media.Type {
   case rakuten.MovieType:
      stream, err = media.MovieStream(
         language, rakuten.Player.Widevine, rakuten.Quality.HD,
      )
   case rakuten.TvShowType:
      var episode string
      err = c.cache.Get("Episode", &episode)
      if err != nil {
         return err
      }
      stream, err = media.EpisodeStream(
         episode, language, rakuten.Player.Widevine, rakuten.Quality.HD,
      )
   }
   if err != nil {
      return err
   }
   c.job.Send = stream.Widevine
   return c.job.DownloadDash(dash.Body, dash.Url, c.dash)
}

func (c *command) do_season() error {
   var media rakuten.Media
   err := c.cache.Get("Media", &media)
   if err != nil {
      return err
   }
   season, err := media.RequestSeason(c.season)
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
func main() {
   maya.SetProxy(func(req *http.Request) (string, bool) {
      // everything needs proxy
      switch path.Ext(req.URL.Path) {
      case ".isma", ".ismv":
         return "", false
      }
      return "", true
   })
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

type command struct {
   cache maya.Cache
   // 1
   address string
   // 2
   season string
   // 3
   episode  string
   language string
   // 4
   dash string
   job  maya.WidevineJob
}

func (c *command) run() error {
   c.cache.Init("L3")
   c.job.ClientId = c.cache.Join("client_id.bin")
   c.job.PrivateKey = c.cache.Join("private_key.pem")
   c.cache.Init("rakuten")
   // 1
   flag.StringVar(&c.address, "a", "", "address")
   // 2
   flag.StringVar(&c.season, "s", "", "season ID")
   // 3
   flag.StringVar(&c.episode, "e", "", "episode ID")
   flag.StringVar(&c.language, "A", "", "audio language")
   // 4
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.job.ClientId, "c", c.job.ClientId, "client ID")
   flag.StringVar(&c.job.PrivateKey, "p", c.job.PrivateKey, "private key")
   flag.Parse()
   if c.address != "" {
      return c.do_address()
   }
   if c.season != "" {
      return c.do_season()
   }
   if c.language != "" {
      return c.do_language()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return maya.Usage([][]string{
      {"a"},
      {"s"},
      {"e", "A"},
      {"d", "c", "p"},
   })
}

func (c *command) do_address() error {
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
   return c.cache.Set("Media", media)
}

func (c *command) do_language() error {
   var media rakuten.Media
   err := c.cache.Get("Media", &media)
   if err != nil {
      return err
   }
   var stream *rakuten.StreamData
   switch media.Type {
   case rakuten.MovieType:
      stream, err = media.MovieStream(
         c.language, rakuten.Player.Widevine, rakuten.Quality.FHD,
      )
   case rakuten.TvShowType:
      stream, err = media.EpisodeStream(
         c.episode, c.language, rakuten.Player.Widevine, rakuten.Quality.FHD,
      )
   }
   if err != nil {
      return err
   }
   dash, err := stream.Dash()
   if err != nil {
      return err
   }
   err = c.cache.Set("Dash", dash)
   if err != nil {
      return err
   }
   err = c.cache.Set("Episode", c.episode)
   if err != nil {
      return err
   }
   err = c.cache.Set("Language", c.language)
   if err != nil {
      return err
   }
   return maya.ListDash(dash.Body, dash.Url)
}
