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
   net.Transport(func(req *http.Request) string {
      switch path.Ext(req.URL.Path) {
      case ".isma", ".ismv":
         return ""
      }
      return "LP"
   })
   log.SetFlags(log.Ltime)
   var tool runner
   err := tool.run()
   if err != nil {
      log.Fatal(err)
   }
}

func (r *runner) run() error {
   var err error
   r.cache, err = os.UserCacheDir()
   if err != nil {
      return err
   }
   r.cache = filepath.ToSlash(r.cache)
   r.config.ClientId = r.cache + "/L3/client_id.bin"
   r.config.PrivateKey = r.cache + "/L3/private_key.pem"
   flag.StringVar(&r.config.ClientId, "C", r.config.ClientId, "client ID")
   flag.StringVar(&r.config.PrivateKey, "P", r.config.PrivateKey, "private key")
   flag.StringVar(&r.season, "S", "", "season ID")
   flag.StringVar(&r.language, "a", "", "audio language")
   flag.StringVar(&r.dash, "d", "", "DASH ID")
   flag.StringVar(&r.episode, "e", "", "episode ID")
   flag.StringVar(&r.movie, "m", "", "movie URL")
   flag.StringVar(&r.show, "s", "", "TV show URL")
   flag.IntVar(&r.config.Threads, "t", 2, "threads")
   flag.Parse()
   if r.movie != "" {
      return r.do_movie()
   }
   if r.show != "" {
      return r.do_show()
   }
   if r.season != "" {
      return r.do_season()
   }
   if r.language != "" {
      if r.dash != "" {
         return r.do_dash()
      }
      return r.do_language()
   }
   flag.Usage()
   return nil
}

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func (r *runner) do_movie() error {
   var movie rakuten.Movie
   err := movie.ParseURL(r.movie)
   if err != nil {
      return err
   }
   item, err := movie.Request()
   if err != nil {
      return err
   }
   fmt.Println(item)
   data, err := json.Marshal(rakuten.Cache{Movie: &movie})
   if err != nil {
      return err
   }
   return write_file(r.cache+"/rakuten/Cache", data)
}

// print seasons
func (r *runner) do_show() error {
   var show rakuten.TvShow
   err := show.ParseURL(r.show)
   if err != nil {
      return err
   }
   show_data, err := show.Request()
   if err != nil {
      return err
   }
   fmt.Println(show_data)
   data, err := json.Marshal(rakuten.Cache{TvShow: &show})
   if err != nil {
      return err
   }
   return write_file(r.cache+"/rakuten/Cache", data)
}

// print episodes
func (r *runner) do_season() error {
   data, err := os.ReadFile(r.cache + "/rakuten/Cache")
   if err != nil {
      return err
   }
   var cache rakuten.Cache
   err = json.Unmarshal(data, &cache)
   if err != nil {
      return err
   }
   season, err := cache.TvShow.RequestSeason(r.season)
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

func (r *runner) do_language() error {
   data, err := os.ReadFile(r.cache + "/rakuten/Cache")
   if err != nil {
      return err
   }
   var cache rakuten.Cache
   err = json.Unmarshal(data, &cache)
   if err != nil {
      return err
   }
   var stream *rakuten.StreamData
   switch {
   case cache.Movie != nil:
      stream, err = cache.Movie.RequestStream(
         r.language, rakuten.Player.Widevine, rakuten.Quality.FHD,
      )
   case cache.TvShow != nil:
      stream, err = cache.TvShow.RequestStream(
         r.episode, r.language, rakuten.Player.Widevine, rakuten.Quality.FHD,
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
   err = write_file(r.cache + "/rakuten/Cache", data)
   if err != nil {
      return err
   }
   return net.Representations(cache.MpdBody, cache.Mpd)
}

func (r *runner) do_dash() error {
   data, err := os.ReadFile(r.cache + "/rakuten/Cache")
   if err != nil {
      return err
   }
   var cache rakuten.Cache
   err = json.Unmarshal(data, &cache)
   if err != nil {
      return err
   }
   var stream *rakuten.StreamData
   switch {
   case cache.Movie != nil:
      stream, err = cache.Movie.RequestStream(
         r.language, rakuten.Player.Widevine, rakuten.Quality.HD,
      )
   case cache.TvShow != nil:
      stream, err = cache.TvShow.RequestStream(
         r.episode, r.language, rakuten.Player.Widevine, rakuten.Quality.HD,
      )
   }
   if err != nil {
      return err
   }
   r.config.Send = func(data []byte) ([]byte, error) {
      return stream.Widevine(data)
   }
   return r.config.Download(cache.MpdBody, cache.Mpd, r.dash)
}

type runner struct {
   cache    string
   config   net.Config
   dash     string
   episode  string
   language string
   movie    string
   season   string
   show     string
}
