package main

import (
   "41.neocities.org/media/paramount"
   "41.neocities.org/net"
   "encoding/json"
   "flag"
   "log"
   "net/http"
   "os"
   "path"
   "path/filepath"
)

// INTL does NOT allow anonymous key request, so if you are INTL you
// will need to use US VPN until someone codes the INTL login
func (f *flag_set) do_dash() error {
   data, err := os.ReadFile(f.cache + "/paramount/Cache")
   if err != nil {
      return err
   }
   var cache paramount.Cache
   err = json.Unmarshal(data, &cache)
   if err != nil {
      return err
   }
   at, err := paramount.ComCbsApp.At()
   if err != nil {
      return err
   }
   token, err := at.Token(f.paramount)
   if err != nil {
      return err
   }
   f.config.Send = func(data []byte) ([]byte, error) {
      return token.Widevine(data)
   }
   return f.config.Download(cache.MpdBody, cache.Mpd, f.dash)
}

func (f *flag_set) secret() paramount.AppSecret {
   if f.intl {
      return paramount.ComCbsCa
   }
   return paramount.ComCbsApp
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
   flag.BoolVar(&f.intl, "i", false, "intl")
   flag.StringVar(&f.paramount, "p", "", "paramount ID")
   flag.StringVar(&f.dash, "d", "", "DASH ID")
   flag.Parse()
   return nil
}

func main() {
   net.Transport(func(req *http.Request) string {
      switch path.Ext(req.URL.Path) {
      case ".m4s", ".mp4":
         return ""
      }
      switch path.Base(req.URL.Path) {
      case "anonymous-session-token.json", "getlicense":
         return "L"
      }
      return "LP"
   })
   log.SetFlags(log.Ltime)
   var set flag_set
   err := set.New()
   if err != nil {
      log.Fatal(err)
   }
   switch {
   case set.paramount != "":
      err = set.do_paramount()
   case set.dash != "":
      err = set.do_dash()
   default:
      flag.Usage()
   }
   if err != nil {
      log.Fatal(err)
   }
}

type flag_set struct {
   cache     string
   config    net.Config
   intl      bool
   // 1
   paramount string
   // 2
   dash      string
}

func write_file(name string, data []byte) error {
   log.Println("WriteFile", name)
   return os.WriteFile(name, data, os.ModePerm)
}
func (f *flag_set) do_paramount() error {
   at, err := f.secret().At()
   if err != nil {
      return err
   }
   item, err := at.Item(f.paramount)
   if err != nil {
      return err
   }
   var cache paramount.Cache
   err = item.Mpd(&cache)
   if err != nil {
      return err
   }
   data, err := json.Marshal(cache)
   if err != nil {
      return err
   }
   err = write_file(f.cache + "/paramount/Cache", data)
   if err != nil {
      return err
   }
   return net.Representations(cache.MpdBody, cache.Mpd)
}
