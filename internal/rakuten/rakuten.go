package main

import (
   "41.neocities.org/media/rakuten"
   "41.neocities.org/net"
   "flag"
   "fmt"
   "log"
   "net/http"
   "net/url"
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
   flag.Parse()
   return nil
}

func (f *flag_set) content_language() bool {
   if f.content != "" {
      if f.language != "" {
         return true
      }
   }
   return false
}

func main() {
   log.SetFlags(log.Ltime)
   http.DefaultTransport = &http.Transport{
      Proxy: func(req *http.Request) (*url.URL, error) {
         log.Println(req.Method, req.URL)
         return http.ProxyFromEnvironment(req)
      },
   }
   var set flag_set
   err := set.New()
   if err != nil {
      panic(err)
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
      panic(err)
   }
}

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}

type flag_set struct {
   cache    string
   config   net.Config
   content  string
   filters  net.Filters
   language string
   movie    string
   season   string
   show     string
}

// print episodes
func (f *flag_set) do_season() error {
   data, err := os.ReadFile(f.cache + "/rakuten/Address")
   if err != nil {
      return err
   }
   var address rakuten.Address
   err = address.Set(string(data))
   if err != nil {
      return err
   }
   contents, err := address.Episodes(f.season)
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

func (f *flag_set) do_show() error {
   var address rakuten.Address
   err := address.Set(f.show)
   if err != nil {
      return err
   }
   err = write_file(f.cache+"/rakuten/Address", []byte(f.show))
   if err != nil {
      return err
   }
   seasons, err := address.Seasons()
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

func (f *flag_set) do_movie() error {
   var address rakuten.Address
   err := address.Set(f.movie)
   if err != nil {
      return err
   }
   err = write_file(f.cache+"/rakuten/Address", []byte(f.movie))
   if err != nil {
      return err
   }
   content, err := address.Movie()
   if err != nil {
      return err
   }
   fmt.Println(content)
   return nil
}

func (f *flag_set) do_send() error {
   data, err := os.ReadFile(f.cache + "/rakuten/Address")
   if err != nil {
      return err
   }
   var address rakuten.Address
   err = address.Set(string(data))
   if err != nil {
      return err
   }
   info, err := address.Wvm(f.content, f.language, rakuten.Fhd)
   if err != nil {
      return err
   }
   resp, err := http.Get(info.Url)
   if err != nil {
      return err
   }
   info, err = address.Wvm(f.content, f.language, rakuten.Hd)
   if err != nil {
      return err
   }
   f.config.Send = func(data []byte) ([]byte, error) {
      return info.Widevine(data)
   }
   return f.filters.Filter(resp, &f.config)
}
