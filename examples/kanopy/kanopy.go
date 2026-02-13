package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/kanopy"
   "errors"
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
   c.name = cache + "/kanopy/userCache.xml"
   c.job.ClientId = cache + "/L3/client_id.bin"
   c.job.PrivateKey = cache + "/L3/private_key.pem"
   // 1
   flag.StringVar(&c.email, "e", "", "email")
   flag.StringVar(&c.password, "p", "", "password")
   // 2
   flag.IntVar(&c.kanopy, "k", 0, "Kanopy ID")
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
   if c.kanopy >= 1 {
      return c.do_kanopy()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   maya.Usage([][]string{
      {"e", "p"},
      {"k"},
      {"d", "t", "C", "P"},
   })
   return nil
}

func (c *command) do_dash() error {
   cache, err := maya.Read[user_cache](c.name)
   if err != nil {
      return err
   }
   c.job.Send = func(data []byte) ([]byte, error) {
      return cache.Login.Widevine(cache.StreamInfo, data)
   }
   return c.job.DownloadDash(cache.Dash.Body, cache.Dash.Url, c.dash)
}

type user_cache struct {
   Dash       *kanopy.Dash
   Login      *kanopy.Login
   StreamInfo *kanopy.StreamInfo
}

func main() {
   maya.Transport(func(req *http.Request) string {
      if path.Ext(req.URL.Path) == ".m4s" {
         return ""
      }
      return "LP"
   })
   log.SetFlags(log.Ltime)
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

type command struct {
   name string
   // 1
   email    string
   password string
   // 2
   kanopy int
   // 3
   dash string
   job  maya.WidevineJob
}

func (c *command) do_email_password() error {
   var login kanopy.Login
   err := login.Fetch(c.email, c.password)
   if err != nil {
      return err
   }
   return maya.Write(c.name, &user_cache{Login: &login})
}

func (c *command) do_kanopy() error {
   cache, err := maya.Read[user_cache](c.name)
   if err != nil {
      return err
   }
   member, err := cache.Login.Membership()
   if err != nil {
      return err
   }
   play, err := cache.Login.PlayResponse(member, c.kanopy)
   if err != nil {
      return err
   }
   var ok bool
   cache.StreamInfo, ok = play.Dash()
   if !ok {
      return errors.New(".Dash()")
   }
   cache.Dash, err = cache.StreamInfo.Dash()
   if err != nil {
      return err
   }
   err = maya.Write(c.name, cache)
   if err != nil {
      return err
   }
   return maya.ListDash(cache.Dash.Body, cache.Dash.Url)
}
