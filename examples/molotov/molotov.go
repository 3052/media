package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/molotov"
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
   c.name = cache + "/molotov/userCache.xml"
   c.job.ClientId = cache + "/L3/client_id.bin"
   c.job.PrivateKey = cache + "/L3/private_key.pem"
   // 1
   flag.StringVar(&c.email, "e", "", "email")
   flag.StringVar(&c.password, "p", "", "password")
   // 2
   flag.StringVar(&c.address, "a", "", "address")
   // 3
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.job.ClientId, "C", c.job.ClientId, "client ID")
   flag.StringVar(&c.job.PrivateKey, "P", c.job.PrivateKey, "private key")
   flag.Parse()
   if c.email != "" {
      if c.password != "" {
         return c.do_email_password()
      }
   }
   if c.address != "" {
      return c.do_address()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return maya.Usage([][]string{
      {"e", "p"},
      {"a"},
      {"d", "C", "P"},
   })
}

func (c *command) do_email_password() error {
   var login molotov.Login
   err := login.Fetch(c.email, c.password)
   if err != nil {
      return err
   }
   return maya.Write(c.name, &user_cache{Login: &login})
}

func (c *command) do_address() error {
   cache, err := maya.Read[user_cache](c.name)
   if err != nil {
      return err
   }
   err = cache.Login.Refresh()
   if err != nil {
      return err
   }
   var media molotov.MediaId
   err = media.Parse(c.address)
   if err != nil {
      return err
   }
   view, err := cache.Login.ProgramView(&media)
   if err != nil {
      return err
   }
   cache.Asset, err = cache.Login.Asset(view)
   if err != nil {
      return err
   }
   cache.Dash, err = cache.Asset.Dash()
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
   email    string
   password string
   // 2
   address string
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
      return cache.Asset.Widevine(data)
   }
   return c.job.DownloadDash(cache.Dash.Body, cache.Dash.Url, c.dash)
}

type user_cache struct {
   Asset *molotov.Asset
   Dash  *molotov.Dash
   Login *molotov.Login
}

func main() {
   log.SetFlags(log.Ltime)
   maya.Transport(func(req *http.Request) string {
      if path.Ext(req.URL.Path) == ".m4s" {
         return ""
      }
      return "LP"
   })
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}
