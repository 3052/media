package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/peacock"
   "flag"
   "log"
   "net/http"
   "path"
)

func (c *client) do_dash_id() error {
   _, err := cache.Read(c)
   if err != nil {
      return err
   }
   job.Send = c.Playout.Widevine
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

var job  maya.WidevineJob

func (c *client) do() error {
   job.ClientId, _ = maya.ResolveCache("L3/client_id.bin")
   job.PrivateKey, _ = maya.ResolveCache("L3/private_key.pem")
   err := cache.Setup("rosso/peacock.xml")
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
   c.Cookie, err = peacock.FetchIdSession(c.email, c.password)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_address() error {
   _, err := cache.Read(c)
   if err != nil {
      return err
   }
   var token peacock.Token
   err = token.Fetch(c.Cookie)
   if err != nil {
      return err
   }
   c.Playout, err = token.Playout(path.Base(c.address))
   if err != nil {
      return err
   }
   endpoint, err := c.Playout.Fastly()
   if err != nil {
      return err
   }
   c.Dash, err = endpoint.Dash()
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
   Cookie  *http.Cookie
   Dash    *peacock.Dash
   Playout *peacock.Playout
   // 1
   email    string
   password string
   // 2
   address string
   // 3
   dash_id string
}
