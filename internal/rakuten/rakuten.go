package main

import (
   "41.neocities.org/media/rakuten"
   "41.neocities.org/net"
   "encoding/json"
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
   "path"
   "path/filepath"
)

func main() {
   log.SetFlags(log.Ltime)
   net.Transport(func(req *http.Request) string {
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

func (c *command) run() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   cache = filepath.ToSlash(cache)
   c.config.ClientId = cache + "/L3/client_id.bin"
   c.config.PrivateKey = cache + "/L3/private_key.pem"
   c.name = cache + "/rakuten/user_cache.json"

   flag.StringVar(&c.config.ClientId, "C", c.config.ClientId, "client ID")
   flag.StringVar(&c.config.PrivateKey, "P", c.config.PrivateKey, "private key")
   flag.StringVar(&c.season, "S", "", "season ID")
   flag.StringVar(&c.language, "a", "", "audio language")
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.episode, "e", "", "episode ID")
   flag.StringVar(&c.movie, "m", "", "movie URL")
   flag.StringVar(&c.show, "s", "", "TV show URL")
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
         return c.do_dash()
      }
      return c.do_language()
   }
   flag.Usage()
   return nil
}

func write(name string, cache *user_cache) error {
   data, err := json.Marshal(cache)
   if err != nil {
      return err
   }
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
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

func read(name string) (*user_cache, error) {
   data, err := os.ReadFile(name)
   if err != nil {
      return nil, err
   }
   cache := &user_cache{}
   err = json.Unmarshal(data, cache)
   if err != nil {
      return nil, err
   }
   return cache, nil
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

type user_cache struct {
   Movie   *rakuten.Movie
   Mpd     *url.URL
   MpdBody []byte
   TvShow  *rakuten.TvShow
}

type command struct {
   config   net.Config
   name    string
   // 1
   movie    string
   // 2
   show     string
   // 3
   season   string
   
   // 4
   episode  string
   // 5
   language string
   dash     string
}

///

func (c *command) do_language() error {
   data, err := os.ReadFile(c.name + "/rakuten/user_cache")
   if err != nil {
      return err
   }
   var cache rakuten.user_cache
   err = json.Unmarshal(data, &cache)
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
   data, err = json.Marshal(cache)
   if err != nil {
      return err
   }
   err = write_file(c.name + "/rakuten/user_cache", data)
   if err != nil {
      return err
   }
   return net.Representations(cache.MpdBody, cache.Mpd)
}

func (c *command) do_dash() error {
   data, err := os.ReadFile(c.name + "/rakuten/user_cache")
   if err != nil {
      return err
   }
   var cache rakuten.user_cache
   err = json.Unmarshal(data, &cache)
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
   return c.config.Download(cache.MpdBody, cache.Mpd, c.dash)
}
