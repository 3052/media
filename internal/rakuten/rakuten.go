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
   flag.StringVar(&f.content, "c", "", "content ID")
   flag.Var(&f.filters, "f", net.FilterUsage)
   flag.StringVar(&f.movie, "m", "", "movie URL")
   flag.StringVar(&f.show, "s", "", "TV show URL")
   flag.IntVar(&f.config.Threads, "t", 12, "threads")
   flag.Parse()
   return nil
}

var Transport = http.Transport{
   Protocols: &http.Protocols{}, // github.com/golang/go/issues/25793
   Proxy: func(req *http.Request) (*url.URL, error) {
      switch path.Ext(req.URL.Path) {
      case ".isma", ".ismv":
      default:
         log.Println(req.Method, req.URL)
      }
      return http.ProxyFromEnvironment(req)
   },
}

func main() {
   http.DefaultTransport = &rakuten.Transport
   log.SetFlags(log.Ltime)
   var set flag_set
   err := set.New()
   if err != nil {
      log.Fatal(err)
   }
   switch {
   case set.content_language():
      err = set.do_send()
   case set.movie != "":
      err = set.do_movie()
   case set.season != "":
      err = set.do_season()
   case set.show != "":
      err = set.do_show()
   default:
      flag.Usage()
   }
   if err != nil {
      log.Fatal(err)
   }
}

func (f *flag_set) content_language() bool {
   if f.content != "" {
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

func (f *flag_set) do_movie() error {
   var media rakuten.Media
   err := media.Parse(f.movie)
   if err != nil {
      return err
   }
   err = write_file(f.cache+"/rakuten/Media", []byte(f.movie))
   if err != nil {
      return err
   }
   content, err := media.Movie()
   if err != nil {
      return err
   }
   fmt.Println(content)
   return nil
}

type flag_set struct {
   cache    string
   config   net.Config
   filters  net.Filters
   // 1
   movie    string
   
   // 2
   show     string
   // 3
   season   string
   // 4
   content  string
   language string
}

///

// print seasons
func (f *flag_set) do_show() error {
   var media rakuten.Media
   err := media.Parse(f.show)
   if err != nil {
      return err
   }
   err = write_file(f.cache+"/rakuten/Media", []byte(f.show))
   if err != nil {
      return err
   }
   seasons, err := media.Seasons()
   if err != nil {
      return err
   }
   for i, season := range seasons {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(&season)
   }
   return nil
}

// print episodes
func (f *flag_set) do_season() error {
   data, err := os.ReadFile(f.cache + "/rakuten/Media")
   if err != nil {
      return err
   }
   var media rakuten.Media
   err = media.Parse(string(data))
   if err != nil {
      return err
   }
   contents, err := media.Episodes(f.season)
   if err != nil {
      return err
   }
   for i, content := range contents {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(&content)
   }
   return nil
}

func (f *flag_set) do_send() error {
   data, err := os.ReadFile(f.cache + "/rakuten/Media")
   if err != nil {
      return err
   }
   var media rakuten.Media
   err = media.Parse(string(data))
   if err != nil {
      return err
   }
   info, err := media.Wvm(f.content, f.language, rakuten.Fhd)
   if err != nil {
      return err
   }
   resp, err := http.Get(info.Url)
   if err != nil {
      return err
   }
   info, err = media.Wvm(f.content, f.language, rakuten.Hd)
   if err != nil {
      return err
   }
   f.config.Send = func(data []byte) ([]byte, error) {
      return info.Widevine(data)
   }
   return f.filters.Filter(resp, &f.config)
}
