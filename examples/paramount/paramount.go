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

func (c *command) app_secret() string {
   if c.intl {
      return paramount.AppSecrets[0].ComCbsCa
   }
   return paramount.AppSecrets[0].ComCbsApp
}

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
   flag.StringVar(&c.username, "username", "", "username")
   flag.StringVar(&c.password, "password", "", "password")
   // 1, 2
   flag.BoolVar(&c.intl, "i", false, "intl")
   // 2
   flag.StringVar(&c.paramount, "p", "", "paramount ID")
   flag.StringVar(&c.proxy, "P", "", "proxy")
   // 3
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.BoolVar(&c.cookie, "c", false, "cookie")
   flag.StringVar(&c.job.CertificateChain, "C", c.job.CertificateChain, "certificate chain")
   flag.StringVar(&c.job.EncryptSignKey, "E", c.job.EncryptSignKey, "encrypt sign key")
   flag.Parse()
   maya.SetProxy(func(req *http.Request) (string, bool) {
      switch path.Ext(req.URL.Path) {
      case ".m4s", ".mp4":
         return "", false
      }
      switch path.Base(req.URL.Path) {
      case "anonymous-session-token.json", "getlicense", "rightsmanager.asmx":
         return "", true
      }
      return c.proxy, true
   })
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
      {"U", "P", "i"},
      {"p", "i"},
      {"d", "c", "C", "E"},
   })
}

func (c *command) do_username_password() error {
   at, err := paramount.GetAt(c.app_secret())
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

type user_cache struct {
   Dash      *paramount.Dash
   Cookie    *http.Cookie
   ContentId string
   Intl      bool
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
   item, err := paramount.FetchItem(at, c.paramount)
   if err != nil {
      return err
   }
   cache.Dash, err = item.Dash()
   if err != nil {
      return err
   }
   cache.ContentId = c.paramount
   cache.Intl = c.intl
   err = maya.Write(c.name, cache)
   if err != nil {
      return err
   }
   return maya.ListDash(cache.Dash.Body, cache.Dash.Url)
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
   token, err := paramount.PlayReady(at, cache.ContentId, cache.Cookie)
   if err != nil {
      return err
   }
   c.job.Send = token.Send
   return c.job.DownloadDash(cache.Dash.Body, cache.Dash.Url, c.dash)
}

func main() {
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

type command struct {
   name string
   // 1
   username string
   password string
   // 1, 2
   intl bool
   // 2
   paramount string
   proxy     string
   // 3
   dash   string
   cookie bool
   job    maya.PlayReadyJob
}
