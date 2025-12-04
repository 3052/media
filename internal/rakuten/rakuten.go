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

func (o *options) run() (bool, error) {
   var err error
   o.cache, err = os.UserCacheDir()
   if err != nil {
      return false, err
   }
   o.cache = filepath.ToSlash(o.cache)
   o.config.ClientId = o.cache + "/L3/client_id.bin"
   o.config.PrivateKey = o.cache + "/L3/private_key.pem"
   flag.StringVar(&o.config.ClientId, "C", o.config.ClientId, "client ID")
   flag.StringVar(&o.config.PrivateKey, "P", o.config.PrivateKey, "private key")
   flag.StringVar(&o.season, "S", "", "season ID")
   flag.StringVar(&o.language, "a", "", "audio language")
   flag.StringVar(&o.dash, "d", "", "DASH ID")
   flag.StringVar(&o.episode, "e", "", "episode ID")
   flag.StringVar(&o.movie, "m", "", "movie URL")
   flag.StringVar(&o.show, "s", "", "TV show URL")
   flag.IntVar(&o.config.Threads, "t", 2, "threads")
   flag.Parse()
   if o.movie != "" {
      return true, o.do_movie()
   }
   if o.show != "" {
      return true, o.do_show()
   }
   if o.season != "" {
      return true, o.do_season()
   }
   if o.language != "" {
      if o.dash != "" {
         return true, o.do_dash()
      }
      return true, o.do_language()
   }
   return false, nil
}
func main() {
   net.Transport(func(req *http.Request) string {
      switch path.Ext(req.URL.Path) {
      case ".isma", ".ismv":
         return ""
      }
      return "LP"
   })
   log.SetFlags(log.Ltime)
   var opts options
   did_run, err := opts.run()
   if err != nil {
      log.Fatal(err)
   }
   if !did_run {
      flag.Usage()
   }
}

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func (o *options) do_movie() error {
   var movie rakuten.Movie
   err := movie.ParseURL(o.movie)
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
   return write_file(o.cache+"/rakuten/Cache", data)
}

// print seasons
func (o *options) do_show() error {
   var show rakuten.TvShow
   err := show.ParseURL(o.show)
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
   return write_file(o.cache+"/rakuten/Cache", data)
}

// print episodes
func (o *options) do_season() error {
   data, err := os.ReadFile(o.cache + "/rakuten/Cache")
   if err != nil {
      return err
   }
   var cache rakuten.Cache
   err = json.Unmarshal(data, &cache)
   if err != nil {
      return err
   }
   season, err := cache.TvShow.RequestSeason(o.season)
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

func (o *options) do_language() error {
   data, err := os.ReadFile(o.cache + "/rakuten/Cache")
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
         o.language, rakuten.Player.Widevine, rakuten.Quality.FHD,
      )
   case cache.TvShow != nil:
      stream, err = cache.TvShow.RequestStream(
         o.episode, o.language, rakuten.Player.Widevine, rakuten.Quality.FHD,
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
   err = write_file(o.cache + "/rakuten/Cache", data)
   if err != nil {
      return err
   }
   return net.Representations(cache.MpdBody, cache.Mpd)
}

func (o *options) do_dash() error {
   data, err := os.ReadFile(o.cache + "/rakuten/Cache")
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
         o.language, rakuten.Player.Widevine, rakuten.Quality.HD,
      )
   case cache.TvShow != nil:
      stream, err = cache.TvShow.RequestStream(
         o.episode, o.language, rakuten.Player.Widevine, rakuten.Quality.HD,
      )
   }
   if err != nil {
      return err
   }
   o.config.Send = func(data []byte) ([]byte, error) {
      return stream.Widevine(data)
   }
   return o.config.Download(cache.MpdBody, cache.Mpd, o.dash)
}

type options struct {
   cache    string
   config   net.Config
   dash     string
   episode  string
   language string
   movie    string
   season   string
   show     string
}
