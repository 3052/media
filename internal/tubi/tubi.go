package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/tubi"
   "flag"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

func main() {
   var set flag_set
   err := set.New()
   if err != nil {
      log.Fatal(err)
   }
   http.DefaultTransport = maya.Transport(set.bypass)
   if set.tubi >= 1 {
      err = set.do_tubi()
      if err != nil {
         log.Fatal(err)
      }
   } else {
      flag.Usage()
   }
}

type flag_set struct {
   bypass  string
   cache   string
   config  maya.Config
   filters maya.Filters
   tubi    int
}

func (f *flag_set) do_tubi() error {
   data, err := tubi.NewContent(f.tubi)
   if err != nil {
      return err
   }
   var content tubi.Content
   err = content.Unmarshal(data)
   if err != nil {
      return err
   }
   resp, err := http.Get(content.VideoResources[0].Manifest.Url)
   if err != nil {
      return err
   }
   f.config.Send = func(data []byte) ([]byte, error) {
      return content.VideoResources[0].Widevine(data)
   }
   return f.filters.Filter(resp, &f.config)
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
   flag.Var(&f.filters, "f", maya.FilterUsage)
   flag.StringVar(&f.bypass, "b", ".mp4", "proxy bypass")
   flag.StringVar(&f.config.ClientId, "c", f.config.ClientId, "client ID")
   flag.StringVar(&f.config.PrivateKey, "p", f.config.PrivateKey, "private key")
   flag.IntVar(&f.tubi, "t", 0, "Tubi ID")
   flag.Parse()
   return nil
}
