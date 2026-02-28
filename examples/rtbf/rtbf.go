package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/rtbf"
   "flag"
   "log"
   "net/http"
)

func (c *command) run() error {
   c.cache.Init("L3")
   c.job.ClientId = c.cache.Join("client_id.bin")
   c.job.PrivateKey = c.cache.Join("private_key.pem")
   c.cache.Init("rtbf")
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
      {"a", "x"},
      {"d", "C", "P"},
   })
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
   var account rtbf.Account
   err = c.cache.Get("Account", &account)
   if err != nil {
      return err
   }
   identity, err := account.Identity()
   if err != nil {
      return err
   }
   session, err := identity.Session()
   if err != nil {
      return err
   }
   entitlement, err := session.Entitlement(asset_id)
   if err != nil {
      return err
   }
   err = c.cache.Set("Entitlement", entitlement)
   if err != nil {
      return err
   }
   format, err := entitlement.Dash()
   if err != nil {
      return err
   }
   dash, err := format.Dash()
   if err != nil {
      return err
   }
   err = c.cache.Set("Dash", dash)
   if err != nil {
      return err
   }
   return maya.ListDash(dash.Body, dash.Url)
}

func (c *command) do_email_password() error {
   var account rtbf.Account
   err := account.Fetch(c.email, c.password)
   if err != nil {
      return err
   }
   return c.cache.Set("Account", account)
}

func (c *command) do_dash() error {
   var dash rtbf.Dash
   err := c.cache.Get("Dash", &dash)
   if err != nil {
      return err
   }
   var entitlement rtbf.Entitlement
   err = c.cache.Get("Entitlement", &entitlement)
   if err != nil {
      return err
   }
   c.job.Send = entitlement.Widevine
   return c.job.DownloadDash(dash.Body, dash.Url, c.dash)
}

func main() {
   maya.SetProxy(func(*http.Request) (string, bool) {
      return "", true
   })
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

type command struct {
   cache maya.Cache
   // 1
   email    string
   password string
   // 2
   address string
   // 3
   dash string
   job  maya.WidevineJob
}
