package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/paramount"
   "flag"
   "log"
   "net/http"
   "path"
)

func (c *command) run() error {
   c.cache.Init("SL2000")
   c.job.CertificateChain = c.cache.Join("CertificateChain")
   c.job.EncryptSignKey = c.cache.Join("EncryptSignKey")
   c.cache.Init("paramount")
   // 1
   flag.StringVar(&c.username, "U", "", "username")
   flag.StringVar(&c.password, "P", "", "password")
   // 2
   flag.StringVar(&c.paramount, "p", "", "paramount ID")
   // 3
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.BoolVar(&c.cookie, "c", false, "cookie")
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
      {"p"},
      {"d", "c", "C", "E"},
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
   var content_id string
   err = c.cache.Get("ContentId", content_id)
   if err != nil {
      return err
   }
   var cookie *http.Cookie
   if c.cookie {
      cookie = &http.Cookie{}
      err = c.cache.Get("Cookie", cookie)
      if err != nil {
         return err
      }
   }
   token, err := paramount.PlayReady(at, content_id, cookie)
   if err != nil {
      return err
   }
   c.job.Send = token.Send
   var dash paramount.Dash
   err = c.cache.Get("Dash", &dash)
   if err != nil {
      return err
   }
   return c.job.DownloadDash(dash.Body, dash.Url, c.dash)
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
   cookie, err := paramount.Login(at, c.username, c.password)
   if err != nil {
      return err
   }
   return c.cache.Set("Cookie", cookie)
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
   dash, err := item.Dash()
   if err != nil {
      return err
   }
   err = c.cache.Set("Dash", dash)
   if err != nil {
      return err
   }
   err = c.cache.Set("ContentId", c.paramount)
   if err != nil {
      return err
   }
   return maya.ListDash(dash.Body, dash.Url)
}

type command struct {
   cache maya.Cache
   // 1
   username string
   password string
   // 2
   paramount string
   // 3
   dash   string
   cookie bool
   job    maya.PlayReadyJob
}

func main() {
   maya.SetProxy(func(req *http.Request) (string, bool) {
      switch path.Ext(req.URL.Path) {
      case ".m4s", ".mp4":
         return "", false
      }
      return "", true
   })
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}
