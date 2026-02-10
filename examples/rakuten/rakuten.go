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
   flag.StringVar(&c.movie, "m", "", "movie URL")
   // 2
   flag.StringVar(&c.show, "s", "", "TV show URL")
   // 3
   flag.StringVar(&c.season, "S", "", "season ID")
   // 4
   flag.StringVar(&c.episode, "e", "", "episode ID")
   flag.StringVar(&c.language, "a", "", "audio language")
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.job.ClientId, "C", c.job.ClientId, "client ID")
   flag.StringVar(&c.job.PrivateKey, "P", c.job.PrivateKey, "private key")
   flag.Parse()
   if c.movie != "" {
      return c.do_movie()
   }
   if c.show != "" {
      return c.do_show()
   }
   if c.season != "" {
      return c.do_season()
   }
   if c.language != "" {
      if c.dash != "" {
         return c.do_language_dash()
      }
      return c.do_language()
   }
   maya.Usage([][]string{
      {"m"},
      {"s"},
      {"S"},
      {"e", "a", "d", "C", "P"},
   })
   return nil
}

func (c *command) do_movie() error {
   var movie rakuten.Movie
   err := movie.ParseURL(c.movie)
   if err != nil {
      return err
   }
   item, err := movie.Request()
   if err != nil {
      return err
   }
   fmt.Println(item)
   return maya.Write(c.name, &user_cache{Movie: &movie})
}

// print seasons
func (c *command) do_show() error {
   var show rakuten.TvShow
   err := show.ParseURL(c.show)
   if err != nil {
      return err
   }
   show_data, err := show.Request()
   if err != nil {
      return err
   }
   fmt.Println(show_data)
   return maya.Write(c.name, &user_cache{TvShow: &show})
}

// print episodes
func (c *command) do_season() error {
   cache, err := maya.Read[user_cache](c.name)
   if err != nil {
      return err
   }
   season, err := cache.TvShow.RequestSeason(c.season)
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

type command struct {
   name string
   // 1
   movie string
   // 2
   show string
   // 3
   season string
   // 4
   episode  string
   language string
   dash     string
   job      maya.WidevineJob
}

func (c *command) do_language_dash() error {
   cache, err := maya.Read[user_cache](c.name)
   if err != nil {
      return err
   }
   var stream *rakuten.StreamData
   switch {
   case cache.Movie != nil:
      stream, err = cache.Movie.RequestStream(
         c.language, rakuten.Player.Widevine, rakuten.Quality.HD,
      )
   case cache.TvShow != nil:
      stream, err = cache.TvShow.RequestStream(
         c.episode, c.language, rakuten.Player.Widevine, rakuten.Quality.HD,
      )
   }
   if err != nil {
      return err
   }
   c.job.Send = func(data []byte) ([]byte, error) {
      return stream.Widevine(data)
   }
   return c.job.DownloadDash(cache.Dash.Body, cache.Dash.Url, c.dash)
}

func (c *command) do_language() error {
   cache, err := maya.Read[user_cache](c.name)
   if err != nil {
      return err
   }
   var stream *rakuten.StreamData
   switch {
   case cache.Movie != nil:
      stream, err = cache.Movie.RequestStream(
         c.language,
         rakuten.Player.Widevine, rakuten.Quality.FHD,
      )
   case cache.TvShow != nil:
      stream, err = cache.TvShow.RequestStream(
         c.episode, c.language,
         rakuten.Player.Widevine, rakuten.Quality.FHD,
      )
   }
   if err != nil {
      return err
   }
   cache.Dash, err = stream.Dash()
   if err != nil {
      return err
   }
   err = maya.Write(c.name, cache)
   if err != nil {
      return err
   }
   return maya.ListDash(cache.Dash.Body, cache.Dash.Url)
}

func main() {
   log.SetFlags(log.Ltime)
   maya.Transport(func(req *http.Request) string {
      switch path.Ext(req.URL.Path) {
      case ".isma", ".ismv":
         return ""
      }
      return "LP"
   })
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

type user_cache struct {
   Dash   *rakuten.Dash
   Movie  *rakuten.Movie
   TvShow *rakuten.TvShow
}
