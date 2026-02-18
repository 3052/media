package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/rakuten"
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
   "path"
   "path/filepath"
)

func (c *command) do_dash() error {
   cache, err := maya.Read[user_cache](c.name)
   if err != nil {
      return err
   }
   var stream *rakuten.StreamData
   switch cache.Media.Type {
   case rakuten.MovieType:
      stream, err = cache.Media.MovieStream(
         cache.Language,
         rakuten.Player.Widevine,
         rakuten.Quality.HD,
      )
   case rakuten.TvShowType:
      stream, err = cache.Media.EpisodeStream(
         cache.Episode,
         cache.Language,
         rakuten.Player.Widevine,
         rakuten.Quality.HD,
      )
   }
   if err != nil {
      return err
   }
   c.job.Send = stream.Widevine
   return c.job.DownloadDash(cache.Dash.Body, cache.Dash.Url, c.dash)
}

func (c *command) run() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   cache = filepath.ToSlash(cache)
   c.name = cache + "/rakuten/userCache.xml"
   c.job.ClientId = cache + "/L3/client_id.bin"
   c.job.PrivateKey = cache + "/L3/private_key.pem"
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
   return maya.Write(c.name, &user_cache{Media: &media})
}

func (c *command) do_season() error {
   cache, err := maya.Read[user_cache](c.name)
   if err != nil {
      return err
   }
   season, err := cache.Media.RequestSeason(c.season)
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
   log.SetFlags(log.Ltime)
   maya.Transport(func(req *http.Request) string {
      // everything needs proxy
      switch path.Ext(req.URL.Path) {
      case ".isma", ".ismv":
         return "P"
      }
      return "LP"
   })
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

type command struct {
   name string
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

type user_cache struct {
   Dash     *rakuten.Dash
   Episode  string
   Language string
   Media    *rakuten.Media
}

func (c *command) do_language() error {
   cache, err := maya.Read[user_cache](c.name)
   if err != nil {
      return err
   }
   var stream *rakuten.StreamData
   switch cache.Media.Type {
   case rakuten.MovieType:
      stream, err = cache.Media.MovieStream(
         c.language, rakuten.Player.Widevine, rakuten.Quality.FHD,
      )
   case rakuten.TvShowType:
      stream, err = cache.Media.EpisodeStream(
         c.episode, c.language, rakuten.Player.Widevine, rakuten.Quality.FHD,
      )
   }
   if err != nil {
      return err
   }
   cache.Dash, err = stream.Dash()
   if err != nil {
      return err
   }
   cache.Episode = c.episode
   cache.Language = c.language
   err = maya.Write(c.name, cache)
   if err != nil {
      return err
   }
   return maya.ListDash(cache.Dash.Body, cache.Dash.Url)
}
