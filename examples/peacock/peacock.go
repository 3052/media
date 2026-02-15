package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/peacock"
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
   c.job.ClientId = cache + "/L3/client_id.bin"
   c.job.PrivateKey = cache + "/L3/private_key.pem"
   c.name = cache + "/peacock/userCache.xml"
   // 1
   flag.StringVar(&c.email, "email", "", "email")
   flag.StringVar(&c.password, "password", "", "password")
   // 2
   flag.StringVar(&c.peacock, "p", "", "Peacock ID")
   // 3
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.IntVar(&c.job.Threads, "t", 2, "threads")
   flag.StringVar(&c.job.ClientId, "C", c.job.ClientId, "client ID")
   flag.StringVar(&c.job.PrivateKey, "P", c.job.PrivateKey, "private key")
   flag.Parse()
   if c.email != "" {
      if c.password != "" {
         return c.do_email_password()
      }
   }
   if c.peacock != "" {
      return c.do_peacock()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return maya.Usage([][]string{
      {"email", "password"},
      {"p"},
      {"d", "t", "C", "P"},
   })
}

func (c *command) do_email_password() error {
   var (
      cache user_cache
      err   error
   )
   cache.Cookie, err = peacock.FetchIdSession(c.email, c.password)
   if err != nil {
      return err
   }
   return maya.Write(c.name, &cache)
}

type command struct {
   name string
   // 1
   email    string
   password string
   // 2
   peacock string
   // 3
   dash string
   job  maya.WidevineJob
}

func (c *command) do_dash() error {
   cache, err := maya.Read[user_cache](c.name)
   if err != nil {
      return err
   }
   c.job.Send = func(data []byte) ([]byte, error) {
      return cache.Playout.Widevine(data)
   }
   return c.job.DownloadDash(cache.Dash.Body, cache.Dash.Url, c.dash)
}

func (c *command) do_peacock() error {
   cache, err := maya.Read[user_cache](c.name)
   if err != nil {
      return err
   }
   var token peacock.Token
   err = token.Fetch(cache.Cookie)
   if err != nil {
      return err
   }
   cache.Playout, err = token.Playout(c.peacock)
   if err != nil {
      return err
   }
   endpoint, err := cache.Playout.Fastly()
   if err != nil {
      return err
   }
   cache.Dash, err = endpoint.Dash()
   if err != nil {
      return err
   }
   err = maya.Write(c.name, cache)
   if err != nil {
      return err
   }
   return maya.ListDash(cache.Dash.Body, cache.Dash.Url)
}

type user_cache struct {
   Cookie  *http.Cookie
   Dash    *peacock.Dash
   Playout *peacock.Playout
}
func main() {
   maya.Transport(func(req *http.Request) string {
      if path.Ext(req.URL.Path) == ".m4s" {
         return ""
      }
      return "L"
   })
   log.SetFlags(log.Ltime)
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}
