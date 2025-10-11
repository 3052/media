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
   if set.movie != "" {
      err = set.do_movie()
   } else if set.show != "" {
      err = set.do_show()
   } else if set.season != "" {
      err = set.do_season()
   } else if set.content != "" {
      if set.language != "" {
         err = set.do_content()
      }
   } else {
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

func (f *flag_set) New() error {
   var err error
   f.cache, err = os.UserCacheDir()
   if err != nil {
      return err
   }
   f.cache = filepath.ToSlash(f.cache)
   f.cdm.ClientId = f.cache + "/L3/client_id.bin"
   f.cdm.PrivateKey = f.cache + "/L3/private_key.pem"
   flag.StringVar(&f.cdm.ClientId, "C", f.cdm.ClientId, "client ID")
   flag.StringVar(&f.cdm.PrivateKey, "P", f.cdm.PrivateKey, "private key")
   flag.StringVar(&f.season, "S", "", "season ID")
   flag.StringVar(&f.language, "a", "", "audio language")
   flag.StringVar(&f.content, "c", "", "content ID")
   flag.StringVar(&f.movie, "m", "", "movie URL")
   flag.StringVar(&f.show, "s", "", "TV show URL")
   flag.Parse()
   return nil
}

type flag_set struct {
   cdm      net.Cdm
   content  string
   filters  net.Filters
   language string
   cache    string
   movie    string
   season   string
   show     string
}
func (f *flag_set) do_content() error {
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
   f.cdm.License = func(data []byte) ([]byte, error) {
      return info.License(data)
   }
   return f.filters.Filter(resp, &f.cdm)
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
