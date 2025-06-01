package main

import (
   "41.neocities.org/net"
   "flag"
   "log"
   "net/http"
   "net/url"
   "os"
   "path/filepath"
)

func main() {
   var set flag_set
   err := set.New()
   if err != nil {
      panic(err)
   }
   http.DefaultTransport = net.Transport(set.proxy)
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
   f.media, err = os.UserHomeDir()
   if err != nil {
      return err
   }
   f.media = filepath.ToSlash(f.media) + "/media"
   f.cdm.ClientId = f.media + "/client_id.bin"
   f.cdm.PrivateKey = f.media + "/private_key.pem"
   flag.StringVar(&f.cdm.ClientId, "C", f.cdm.ClientId, "client ID")
   flag.StringVar(&f.cdm.PrivateKey, "P", f.cdm.PrivateKey, "private key")
   flag.StringVar(&f.season, "S", "", "season ID")
   flag.StringVar(&f.language, "a", "", "audio language")
   flag.StringVar(&f.content, "c", "", "content ID")
   flag.StringVar(&f.movie, "m", "", "movie URL")
   flag.StringVar(&f.show, "s", "", "TV show URL")
   flag.Func("x", "proxy", func(data string) error {
      var err error
      f.proxy, err = url.Parse(data)
      return err
   })
   flag.Parse()
   return nil
}

type flag_set struct {
   cdm      net.Cdm
   content  string
   filters  net.Filters
   language string
   media    string
   movie    string
   proxy    *url.URL
   season   string
   show     string
}
