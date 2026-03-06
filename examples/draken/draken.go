package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/draken"
   "flag"
   "log"
   "path"
)

func (c *client) do_email_password() error {
   c.Login = &draken.Login{}
   err := c.Login.Fetch(c.email, c.password)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_dash_id() error {
   err := cache.Read(c)
   if err != nil {
      return err
   }
   job.Send = func(data []byte) ([]byte, error) {
      return c.Login.Widevine(c.Playback, data)
   }
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
   err := cache.Setup("rosso/draken.xml")
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
   Dash     *draken.Dash
   Login    *draken.Login
   Playback *draken.Playback
   // 1
   email    string
   password string
   // 2
   address string
   // 3
   dash_id string
}

func (c *client) do_address() error {
   movie, err := draken.FetchMovie(path.Base(c.address))
   if err != nil {
      return err
   }
   err = cache.Update(c, func() error {
      entitlement, err := c.Login.Entitlement(movie)
      if err != nil {
         return err
      }
      c.Playback, err = c.Login.Playback(movie, entitlement)
      if err != nil {
         return err
      }
      c.Dash, err = c.Playback.Dash()
      return err
   })
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}
