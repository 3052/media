package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/paramount"
   "flag"
   "log"
   "net/http"
   "os"
   "path"
   "path/filepath"
)

func (c *command) run() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   cache = filepath.ToSlash(cache)
   c.name = cache + "/paramount/userCache.xml"
   c.job.CertificateChain = cache + "/SL2000/CertificateChain"
   c.job.EncryptSignKey = cache + "/SL2000/EncryptSignKey"
   // 1
   flag.StringVar(&c.username, "U", "", "username")
   flag.StringVar(&c.password, "P", "", "password")
   // 2
   flag.StringVar(&c.paramount, "p", "", "paramount ID")
   // 2, 3
   flag.BoolVar(&c.intl, "i", false, "intl")
   // 3
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.BoolVar(&c.cookie, "c", false, "cookie")
   flag.IntVar(&c.job.Threads, "t", 2, "threads")
   flag.StringVar(&c.job.CertificateChain, "C", c.job.CertificateChain, "certificate chain")
   flag.StringVar(&c.job.EncryptSignKey, "E", c.job.EncryptSignKey, "encrypt sign key")
   flag.Parse()
   if c.username != "" {
      if c.password != "" {
         return c.do_username_password()
      }
   }
   if c.paramount != "" {
      return c.do_paramount()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return maya.Usage([][]string{
      {"U", "P"},
      {"p", "i"},
      {"d", "i", "c", "t", "C", "E"},
   })
}

type user_cache struct {
   Cookie *http.Cookie
   Dash   *paramount.Dash
   Item   *paramount.Item
}

func (c *command) do_username_password() error {
   at, err := paramount.GetAt(paramount.ComCbsApp.AppSecret)
   if err != nil {
      return err
   }
   var cache user_cache
   cache.Cookie, err = paramount.Login(at, c.username, c.password)
   if err != nil {
      return err
   }
   return maya.Write(c.name, &cache)
}

func (c *command) do_paramount() error {
   at, err := paramount.GetAt(c.app_secret())
   if err != nil {
      return err
   }
   cache, err := maya.Read[user_cache](c.name)
   if err != nil {
      cache = &user_cache{}
   }
   cache.Item, err = paramount.FetchItem(at, c.paramount)
   if err != nil {
      return err
   }
   cache.Dash, err = cache.Item.Dash()
   if err != nil {
      return err
   }
   err = maya.Write(c.name, cache)
   if err != nil {
      return err
   }
   return maya.ListDash(cache.Dash.Body, cache.Dash.Url)
}

type command struct {
   name string
   // 1
   username string
   password string
   // 2
   paramount string
   // 2, 3
   intl bool
   // 3
   dash   string
   cookie bool
   job    maya.PlayReadyJob
}

func (c *command) do_dash() error {
   at, err := paramount.GetAt(c.app_secret())
   if err != nil {
      return err
   }
   cache, err := maya.Read[user_cache](c.name)
   if err != nil {
      return err
   }
   if !c.cookie {
      cache.Cookie = nil
   }
   token, err := paramount.PlayReady(at, cache.Item.ContentId, cache.Cookie)
   if err != nil {
      return err
   }
   c.job.Send = func(data []byte) ([]byte, error) {
      return token.Send(data)
   }
   return c.job.DownloadDash(cache.Dash.Body, cache.Dash.Url, c.dash)
}

func (c *command) app_secret() string {
   if c.intl {
      return paramount.ComCbsCa.AppSecret
   }
   return paramount.ComCbsApp.AppSecret
}
func main() {
   log.SetFlags(log.Ltime)
   maya.Transport(func(req *http.Request) string {
      switch path.Ext(req.URL.Path) {
      case ".m4s", ".mp4":
         return ""
      }
      switch path.Base(req.URL.Path) {
      case "anonymous-session-token.json", "getlicense", "rightsmanager.asmx":
         return "L"
      }
      return "LP"
   })
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}
