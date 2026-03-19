package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/peacock"
   "flag"
   "log"
   "net/http"
   "path"
)

func (c *client) do() error {
   err := cache.Setup("rosso/peacock.xml")
   if err != nil {
      return err
   }
   read_err := cache.Read(c)
   // 1
   widevine := maya.StringVar(&c.Job.Widevine, "w", "Widevine")
   // 2
   email := maya.StringVar(&c.email, "e", "email")
   password := maya.StringVar(&c.password, "p", "password")
   // 3
   address := maya.StringVar(&c.address, "a", "address")
   // 4
   dash_id := maya.StringVar(&c.dash_id, "d", "DASH ID")
   set := maya.Parse()
   switch {
   case len(set) == 0:
      return maya.Usage([][]*flag.Flag{
         {widevine},
         {email, password},
         {address},
         {dash_id},
      })
   case set[widevine.Name]:
      return cache.Write(c)
   case set[email.Name] && set[password.Name]:
      return c.do_email_password()
   case read_err != nil:
      return read_err
   case set[address.Name]:
      return c.do_address()
   case set[dash_id.Name]:
      return c.do_dash_id()
   }
   return nil
}

func (c *client) do_dash_id() error {
   return c.Job.DownloadDash(
      c.Dash.Body, c.Dash.Url, c.dash_id, c.Playout.Widevine,
   )
}

func (c *client) do_address() error {
   token, err := peacock.FetchToken(c.Cookie)
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

func (c *client) do_email_password() error {
   var err error
   c.Cookie, err = peacock.FetchIdSession(c.email, c.password)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

var cache maya.Cache

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.m4s")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   Cookie  *http.Cookie
   Dash    *peacock.Dash
   Playout *peacock.Playout
   // 1
   Job maya.Job
   // 2
   email    string
   password string
   // 3
   address string
   // 4
   dash_id string
}
