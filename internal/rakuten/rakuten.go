package main

import (
   "41.neocities.org/media/rakuten"
   "41.neocities.org/net"
   "flag"
   "fmt"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func (f *flag_set) do_movie() error {
   var movie rakuten.Movie
   err = movie.ParseURL(f.movie)
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
   return write_file(f.cache+"/rakuten/Cache", data)
}

// print seasons
func (f *flag_set) do_tv_show() error {
   var show rakuten.TvShow
   err = show.ParseURL(f.show)
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
   return write_file(f.cache+"/rakuten/Cache", data)
}

// print episodes
func (f *flag_set) do_season() error {
   data, err := os.ReadFile(f.cache + "/rakuten/Cache")
   if err != nil {
      return err
   }
   var cache rakuten.Cache
   err = json.Unmarshal(data, &cache)
   if err != nil {
      return err
   }
   season, err := cache.TvShow.RequestSeason(f.season)
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

func (f *flag_set) do_movie_dash() error {
   data, err := os.ReadFile(f.cache + "/rakuten/Cache")
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
         f.language, rakuten.Player.Widevine, rakuten.Quality.FHD,
      )
   case cache.TvShow != nil:
      stream, err = cache.TvShow.RequestStream(
         f.item, f.language, rakuten.Player.Widevine, rakuten.Quality.FHD,
      )
   }
   if err != nil {
      return err
   }
   err = stream.Mpd(&cache)
   if err != nil {
      return err
   }
   return net.Representations(cache.MpdBody, cache.Mpd)
}

func (f *flag_set) do_episode_dash() error {
   data, err := os.ReadFile(f.cache + "/rakuten/Cache")
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
         f.language, rakuten.Player.Widevine, rakuten.Quality.FHD,
      )
   case cache.TvShow != nil:
      stream, err = cache.TvShow.RequestStream(
         f.item, f.language, rakuten.Player.Widevine, rakuten.Quality.FHD,
      )
   }
   if err != nil {
      return err
   }
   err = stream.Mpd(&cache)
   if err != nil {
      return err
   }
   return net.Representations(cache.MpdBody, cache.Mpd)
}

func (f *flag_set) do_movie_license() error {
   data, err := os.ReadFile(f.cache + "/rakuten/Cache")
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
         f.language, rakuten.Player.Widevine, rakuten.Quality.HD,
      )
   case cache.TvShow != nil:
      stream, err = cache.TvShow.RequestStream(
         f.item, f.language, rakuten.Player.Widevine, rakuten.Quality.HD,
      )
   }
   if err != nil {
      return err
   }
   f.config.Send = func(data []byte) ([]byte, error) {
      return stream.Widevine(data)
   }
   return f.config.Download(cache.MpdBody, cache.Mpd, f.dash)
}

func (f *flag_set) do_episode_license() error {
   data, err := os.ReadFile(f.cache + "/rakuten/Cache")
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
         f.language, rakuten.Player.Widevine, rakuten.Quality.HD,
      )
   case cache.TvShow != nil:
      stream, err = cache.TvShow.RequestStream(
         f.item, f.language, rakuten.Player.Widevine, rakuten.Quality.HD,
      )
   }
   if err != nil {
      return err
   }
   f.config.Send = func(data []byte) ([]byte, error) {
      return stream.Widevine(data)
   }
   return f.config.Download(cache.MpdBody, cache.Mpd, f.dash)
}

type flag_set struct {
   cache    string
   config   net.Config
   movie    string
   tv_show     string
   season   string
   language string
   episode string
   dash string
}

///

func (f *flag_set) New() error {
   var err error
   f.cache, err = os.UserCacheDir()
   if err != nil {
      return err
   }
   f.cache = filepath.ToSlash(f.cache)
   f.config.ClientId = f.cache + "/L3/client_id.bin"
   f.config.PrivateKey = f.cache + "/L3/private_key.pem"
   flag.StringVar(&f.config.ClientId, "C", f.config.ClientId, "client ID")
   flag.StringVar(&f.config.PrivateKey, "P", f.config.PrivateKey, "private key")
   flag.StringVar(&f.season, "S", "", "season ID")
   flag.StringVar(&f.language, "a", "", "audio language")
   flag.StringVar(&f.item, "c", "", "item ID")
   flag.StringVar(&f.dash, "d", "", "DASH ID")
   flag.StringVar(&f.movie, "m", "", "movie URL")
   flag.StringVar(&f.show, "s", "", "TV show URL")
   flag.IntVar(&f.config.Threads, "t", 12, "threads")
   return nil
}

func main() {
   net.Transport(func(req *http.Request) string {
      switch path.Ext(req.URL.Path) {
      case ".isma", ".ismv":
         return ""
      }
      return "L"
   })
   flag.Parse()
   log.SetFlags(log.Ltime)
   var set flag_set
   err := set.New()
   if err != nil {
      log.Fatal(err)
   }
   
   if set.movie != "" {
      err = set.do_movie()
   }
   
   
   switch {
   case set.movie != "":
   case set.show != "":
      err = set.do_show()
   case set.season != "":
      err = set.do_season()
   case set.item_language():
      err = set.do_item()
   case set.dash != "":
      err = set.do_dash()
   default:
      flag.Usage()
   }
   if err != nil {
      log.Fatal(err)
   }
}
