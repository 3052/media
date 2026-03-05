package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/paramount"
   "flag"
   "log"
   "net/http"
)

func (c *client) do_dash_id() error {
   app_secret, err := paramount.FetchAppSecret()
   if err != nil {
      return err
   }
   at, err := paramount.GetAt(app_secret)
   if err != nil {
      return err
   }
   err = cache.Read(c)
   if err != nil {
      return err
   }
   if !c.cookie {
      c.Cookie = nil
   }
   token, err := paramount.PlayReady(at, c.paramount, c.Cookie)
   if err != nil {
      return err
   }
   job.Send = token.Send
   return job.DownloadDash(c.Dash.Body, c.Dash.Url, c.dash_id)
}

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.m4s,*.mp4")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

var cache maya.Cache

var job maya.PlayReadyJob

func (c *client) do() error {
   job.CertificateChain, _ = maya.ResolveCache("SL2000/CertificateChain")
   job.EncryptSignKey, _ = maya.ResolveCache("SL2000/EncryptSignKey")
   err := cache.Setup("rosso/paramount.xml")
   if err != nil {
      return err
   }
   // 1
   flag.StringVar(&c.username, "U", "", "username")
   flag.StringVar(&c.password, "P", "", "password")
   // 2
   flag.StringVar(&c.paramount, "p", "", "paramount ID")
   // 3
   flag.StringVar(&c.dash_id, "d", "", "DASH ID")
   flag.BoolVar(&c.cookie, "c", false, "cookie")
   flag.StringVar(&job.CertificateChain, "C", job.CertificateChain, "certificate chain")
   flag.StringVar(&job.EncryptSignKey, "E", job.EncryptSignKey, "encrypt sign key")
   flag.Parse()
   if c.username != "" {
      if c.password != "" {
         return c.do_username_password()
      }
   }
   if c.paramount != "" {
      return c.do_paramount()
   }
   if c.dash_id != "" {
      return c.do_dash_id()
   }
   return maya.Usage([][]string{
      {"U", "P"},
      {"p"},
      {"d", "c", "C", "E"},
   })
}

func (c *client) do_username_password() error {
   app_secret, err := paramount.FetchAppSecret()
   if err != nil {
      return err
   }
   at, err := paramount.GetAt(app_secret)
   if err != nil {
      return err
   }
   c.Cookie, err = paramount.Login(at, c.username, c.password)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_paramount() error {
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
   c.Dash, err = item.Dash()
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}

type client struct {
   Cookie *http.Cookie
   Dash   *paramount.Dash
   // 1
   username string
   password string
   // 2
   paramount string
   // 3
   dash_id string
   cookie  bool
}
