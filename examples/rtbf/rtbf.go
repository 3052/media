package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/rtbf"
   "flag"
   "log"
)

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

var cache maya.Cache

var job  maya.WidevineJob

func (c *client) do() error {
   job.ClientId, _ = maya.ResolveCache("L3/client_id.bin")
   job.PrivateKey, _ = maya.ResolveCache("L3/private_key.pem")
   err := c.cache.Setup("rosso/rtbf.xml")
   if err != nil {
      return err
   }
   // 1
   flag.StringVar(&c.email, "e", "", "email")
   flag.StringVar(&c.password, "p", "", "password")
   // 2
   flag.StringVar(&c.address, "a", "", "address")
   // 3
   flag.StringVar(&c.dash_id, "d", "", "DASH ID")
   flag.StringVar(&job.ClientId, "C", job.ClientId, "client ID")
   flag.StringVar(&job.PrivateKey, "P", job.PrivateKey, "private key")
   flag.Parse()
   if c.email != "" {
      if c.password != "" {
         return c.do_email_password()
      }
   }
   if c.address != "" {
      return c.do_address()
   }
   if c.dash_id != "" {
      return c.do_dash_id()
   }
   return maya.Usage([][]string{
      {"e", "p"},
      {"a"},
      {"d", "C", "P"},
   })
}

type client struct {
   Account     *rtbf.Account
   Dash        *rtbf.Dash
   Entitlement *rtbf.Entitlement
   // 1
   email    string
   password string
   // 2
   address string
   // 3
   dash_id string
}

///

func (c *client) do_email_password() error {
   var account rtbf.Account
   err := account.Fetch(c.email, c.password)
   if err != nil {
      return err
   }
   return c.cache.Write(saved_state{Account: &account})
}

func (c *client) do_dash_id() error {
   err := c.cache.Read(&state)
   if err != nil {
      return err
   }
   job.Send = state.Entitlement.Widevine
   return job.DownloadDash(state.Dash.Body, state.Dash.Url, c.dash_id)
}

func (c *client) do_address() error {
   path, err := rtbf.GetPath(c.address)
   if err != nil {
      return err
   }
   asset_id, err := rtbf.FetchAssetId(path)
   if err != nil {
      return err
   }
   err = c.cache.Read(&state)
   if err != nil {
      return err
   }
   identity, err := state.Account.Identity()
   if err != nil {
      return err
   }
   session, err := identity.Session()
   if err != nil {
      return err
   }
   state.Entitlement, err = session.Entitlement(asset_id)
   if err != nil {
      return err
   }
   format, err := state.Entitlement.Dash()
   if err != nil {
      return err
   }
   state.Dash, err = format.Dash()
   if err != nil {
      return err
   }
   err = c.cache.Write(state)
   if err != nil {
      return err
   }
   return maya.ListDash(state.Dash.Body, state.Dash.Url)
}
