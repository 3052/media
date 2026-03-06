package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/molotov"
   "flag"
   "log"
)

func (c *client) do_dash_id() error {
   err := cache.Read(c)
   if err != nil {
      return err
   }
   job.Send = c.Asset.Widevine
   return job.DownloadDash(c.Dash.Body, c.Dash.Url, c.dash_id)
}

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.m4s")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

var cache maya.Cache

var job maya.WidevineJob

func (c *client) do() error {
   job.ClientId, _ = maya.ResolveCache("L3/client_id.bin")
   job.PrivateKey, _ = maya.ResolveCache("L3/private_key.pem")
   err := cache.Setup("rosso/molotov.xml")
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

func (c *client) do_email_password() error {
   var err error
   c.Login, err = molotov.FetchLogin(c.email, c.password)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

type client struct {
   Asset *molotov.Asset
   Dash  *molotov.Dash
   Login *molotov.Login
   // 1
   email    string
   password string
   // 2
   address string
   // 3
   dash_id string
}

func (c *client) do_address() error {
   var media molotov.MediaId
   err := media.Parse(c.address)
   if err != nil {
      return err
   }
   err = cache.Update(c, func() error {
      err := c.Login.Refresh()
      if err != nil {
         return err
      }
      view, err := c.Login.ProgramView(&media)
      if err != nil {
         return err
      }
      c.Asset, err = c.Login.Asset(view)
      if err != nil {
         return err
      }
      c.Dash, err = c.Asset.Dash()
      return err
   })
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}
