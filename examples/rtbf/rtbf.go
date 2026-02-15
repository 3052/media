package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/media/rtbf"
   "errors"
   "flag"
   "log"
   "net/http"
   "os"
   "path/filepath"
)

func (c *command) run() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   cache = filepath.ToSlash(cache)
   c.name = cache + "/rtbf/userCache.xml"
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
   var account rtbf.Account
   err := account.Fetch(c.email, c.password)
   if err != nil {
      return err
   }
   return maya.Write(c.name, &user_cache{Account: &account})
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

type user_cache struct {
   Account     *rtbf.Account
   Dash        *rtbf.Dash
   Entitlement *rtbf.Entitlement
}

func (c *command) do_address() error {
   path, err := rtbf.GetPath(c.address)
   if err != nil {
      return err
   }
   asset_id, err := rtbf.FetchAssetId(path)
   if err != nil {
      return err
   }
   cache, err := maya.Read[user_cache](c.name)
   if err != nil {
      return err
   }
   identity, err := cache.Account.Identity()
   if err != nil {
      return err
   }
   session, err := identity.Session()
   if err != nil {
      return err
   }
   cache.Entitlement, err = session.Entitlement(asset_id)
   if err != nil {
      return err
   }
   format, ok := cache.Entitlement.Dash()
   if !ok {
      return errors.New(".Dash()")
   }
   cache.Dash, err = format.Dash()
   if err != nil {
      return err
   }
   err = maya.Write(c.name, cache)
   if err != nil {
      return err
   }
   return maya.ListDash(cache.Dash.Body, cache.Dash.Url)
}

func (c *command) do_dash() error {
   cache, err := maya.Read[user_cache](c.name)
   if err != nil {
      return err
   }
   c.job.Send = func(data []byte) ([]byte, error) {
      return cache.Entitlement.Widevine(data)
   }
   return c.job.DownloadDash(cache.Dash.Body, cache.Dash.Url, c.dash)
}

func main() {
   log.SetFlags(log.Ltime)
   maya.Transport(func(*http.Request) string {
      return "L"
   })
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}
