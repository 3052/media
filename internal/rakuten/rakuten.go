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

func (f *flag_set) item_language() bool {
   if f.item != "" {
      if f.language != "" {
         return true
      }
   }
   return false
}

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

func main() {
   net.Transport(func(req *http.Request) string {
      switch path.Ext(req.URL.Path) {
      case ".isma", ".ismv":
         return ""
      }
      return "L"
   })
   log.SetFlags(log.Ltime)
   var set flag_set
   err := set.New()
   if err != nil {
      log.Fatal(err)
   }
   switch {
   case set.movie != "":
      err = set.do_movie()
   case set.show != "":
      err = set.do_show()
   case set.season != "":
      err = set.do_season()
   case set.item_language():
      err = set.do_send()
   case set.dash != "":
      err = set.do_dash()
   default:
      flag.Usage()
   }
   if err != nil {
      log.Fatal(err)
   }
}

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
   flag.Parse()
   return nil
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
func (f *flag_set) do_show() error {
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

type flag_set struct {
   cache    string
   config   net.Config
   // 1
   movie    string
   // 2
   show     string
   
   // 3
   season   string
   // 4
   item  string
   language string
   // 5
   dash string
}

///

// print episodes
func (f *flag_set) do_season() error {
   data, err := os.ReadFile(f.cache + "/rakuten/Cache")
   if err != nil {
      return err
   }
   var media rakuten.Media
   err = media.Parse(string(data))
   if err != nil {
      return err
   }
   items, err := media.Episodes(f.season)
   if err != nil {
      return err
   }
   for i, item := range items {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(&item)
   }
   return nil
}

func (f *flag_set) do_send() error {
   data, err := os.ReadFile(f.cache + "/rakuten/Cache")
   if err != nil {
      return err
   }
   var media rakuten.Media
   err = media.Parse(string(data))
   if err != nil {
      return err
   }
   info, err := media.Wvm(f.item, f.language, rakuten.Fhd)
   if err != nil {
      return err
   }
   resp, err := http.Get(info.Url)
   if err != nil {
      return err
   }
   info, err = media.Wvm(f.item, f.language, rakuten.Hd)
   if err != nil {
      return err
   }
   f.config.Send = func(data []byte) ([]byte, error) {
      return info.Widevine(data)
   }
   return f.filters.Filter(resp, &f.config)
}
