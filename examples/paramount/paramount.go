package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/paramount"
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
   c.job.CertificateChain = cache + "/SL2000/CertificateChain"
   c.job.EncryptSignKey = cache + "/SL2000/EncryptSignKey"
   c.name = cache + "/rosso/paramount.xml"
   // 1
   flag.StringVar(&c.username, "U", "", "username")
   flag.StringVar(&c.password, "P", "", "password")
   // 1, 2, 3
   flag.StringVar(&c.proxy, "x", "", "proxy")
   // 2
   flag.StringVar(&c.paramount, "p", "", "paramount ID")
   // 3
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.BoolVar(&c.cookie, "c", false, "cookie")
   flag.IntVar(&c.job.Threads, "t", 2, "threads")
   flag.StringVar(&c.job.CertificateChain, "C", c.job.CertificateChain, "certificate chain")
   flag.StringVar(&c.job.EncryptSignKey, "E", c.job.EncryptSignKey, "encrypt sign key")
   flag.Parse()
   maya.SetProxy(func(req *http.Request) (string, bool) {
      switch path.Ext(req.URL.Path) {
      case ".m4s", ".mp4":
         return "", false
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
      {"U", "P", "x"},
      {"p", "x"},
      {"d", "x", "c", "t", "C", "E"},
   })
}

func (c *command) do_dash() error {
   app_secret, err := paramount.FetchAppSecret()
   if err != nil {
      return err
   }
   at, err := paramount.GetAt(app_secret)
   if err != nil {
      return err
   }
   var cache user_cache
   err = maya.Read(c.name, &cache)
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

func (c *command) do_username_password() error {
   app_secret, err := paramount.FetchAppSecret()
   if err != nil {
      return err
   }
   at, err := paramount.GetAt(app_secret)
   if err != nil {
      return err
   }
   var cache user_cache
   cache.Cookie, err = paramount.Login(at, c.username, c.password)
   if err != nil {
      return err
   }
   return maya.Write(c.name, cache)
}

func (c *command) do_paramount() error {
   app_secret, err := paramount.FetchAppSecret()
   if err != nil {
      return err
   }
   at, err := paramount.GetAt(app_secret)
   if err != nil {
      return err
   }
   item, err := paramount.FetchItem(at, c.paramount)
   if err != nil {
      return err
   }
   var cache user_cache
   err = maya.Read(c.name, &cache)
   if err != nil {
      log.Print(err)
   }
   cache.Dash, err = item.Dash()
   if err != nil {
      return err
   }
   cache.ContentId = c.paramount
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
   // 1, 2, 3
   proxy string
   // 2
   paramount string
   // 3
   dash   string
   cookie bool
   job    maya.PlayReadyJob
}

type user_cache struct {
   ContentId string
   Cookie    *http.Cookie
   Dash      *paramount.Dash
}

func main() {
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}
