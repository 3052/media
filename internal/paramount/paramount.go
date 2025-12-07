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
func (r *runner) do_dash() error {
   data, err := os.ReadFile(r.cache)
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
   token, err := at.Token(r.paramount)
   if err != nil {
      return err
   }
   r.config.Send = func(data []byte) ([]byte, error) {
      return token.Widevine(data)
   }
   return r.config.Download(cache.MpdBody, cache.Mpd, r.dash)
}

func (r *runner) secret() paramount.AppSecret {
   if r.intl {
      return paramount.ComCbsCa
   }
   return paramount.ComCbsApp
}

func main() {
   log.SetFlags(log.Ltime)
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
   var program runner
   err := program.run()
   if err != nil {
      log.Fatal(err)
   }
}

func (r *runner) run() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   cache = filepath.ToSlash(cache)
   r.config.ClientId = cache + "/L3/client_id.bin"
   r.config.PrivateKey = cache + "/L3/private_key.pem"
   r.cache = cache + "/paramount/Cache.json"
   
   flag.StringVar(&r.config.ClientId, "C", r.config.ClientId, "client ID")
   flag.StringVar(&r.config.PrivateKey, "P", r.config.PrivateKey, "private key")
   flag.StringVar(&r.dash, "d", "", "DASH ID")
   flag.BoolVar(&r.intl, "i", false, "intl")
   flag.StringVar(&r.paramount, "p", "", "paramount ID")
   flag.Parse()
   if r.paramount != "" {
      return r.do_paramount()
   }
   if r.dash != "" {
      return r.do_dash()
   }
   flag.Usage()
   return nil
}

type runner struct {
   cache     string
   config    net.Config
   // 1
   paramount string
   intl      bool
   // 2
   dash      string
}

func (r *runner) do_paramount() error {
   at, err := r.secret().At()
   if err != nil {
      return err
   }
   item, err := at.Item(r.paramount)
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
   log.Println("WriteFile", r.cache)
   err = os.WriteFile(r.cache, data, os.ModePerm)
   if err != nil {
      return err
   }
   return net.Representations(cache.MpdBody, cache.Mpd)
}
