package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/rakuten"
   "encoding/xml"
   "flag"
   "fmt"
   "os"
   "path/filepath"
)

func (c *command) do_language() error {
   data, err := os.ReadFile(c.name + "/rakuten/userCache.xml")
   if err != nil {
      return err
   }
   var cache user_cache
   err = xml.Unmarshal(data, &cache)
   if err != nil {
      return err
   }
   var stream *rakuten.StreamData
   switch {
   case cache.Movie != nil:
      stream, err = cache.Movie.RequestStream(
         c.language, rakuten.Player.Widevine, rakuten.Quality.FHD,
      )
   case cache.TvShow != nil:
      stream, err = cache.TvShow.RequestStream(
         c.episode, c.language, rakuten.Player.Widevine, rakuten.Quality.FHD,
      )
   }
   if err != nil {
      return err
   }
   err = stream.Mpd(&cache)
   if err != nil {
      return err
   }
   err = write(c.name, &cache)
   if err != nil {
      return err
   }
   return maya.Representations(cache.Mpd.Url, cache.Mpd.Body)
}

func (c *command) do_dash() error {
   data, err := os.ReadFile(c.name + "/rakuten/userCache.xml")
   if err != nil {
      return err
   }
   var cache user_cache
   err = xml.Unmarshal(data, &cache)
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
   c.config.Send = func(data []byte) ([]byte, error) {
      return stream.Widevine(data)
   }
   return c.config.Download(cache.Mpd.Url, cache.Mpd.Body, c.dash)
}

func (c *command) run() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   cache = filepath.ToSlash(cache)
   c.config.ClientId = cache + "/L3/client_id.bin"
   c.config.PrivateKey = cache + "/L3/private_key.pem"
   c.name = cache + "/rakuten/userCache.xml"

   flag.StringVar(&c.config.ClientId, "C", c.config.ClientId, "client ID")
   flag.StringVar(&c.config.PrivateKey, "P", c.config.PrivateKey, "private key")
   flag.StringVar(&c.season, "S", "", "season ID")
   flag.StringVar(&c.language, "a", "", "audio language")
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.episode, "e", "", "episode ID")
   flag.StringVar(&c.movie, "m", "", "movie URL")
   flag.StringVar(&c.show, "s", "", "TV show URL")
   flag.Parse()

   if c.show != "" {
      return c.do_show()
   }
   if c.season != "" {
      return c.do_season()
   }
   
   if c.movie != "" {
      return c.do_movie()
   }
   if c.language != "" {
      if c.dash != "" {
         return c.do_dash()
      }
      return c.do_language()
   }
   flag.Usage()
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
   return write(c.name, &user_cache{Movie: &movie})
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
   return write(c.name, &user_cache{TvShow: &show})
}

// print episodes
func (c *command) do_season() error {
   cache, err := read(c.name)
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
