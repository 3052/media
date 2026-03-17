package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/peacock"
   "flag"
   "log"
   "net/http"
   "path"
)

var cache maya.Cache

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.m4s")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do() error {
   err := cache.Setup("rosso/peacock.xml")
   if err != nil {
      return err
   }
   err = cache.Read(c)
   // 1
   flag.StringVar(&c.Job.Widevine, "w", c.Job.Widevine, "Widevine")
   // 2
   flag.StringVar(&c.email, "e", "", "email")
   flag.StringVar(&c.password, "p", "", "password")
   // 3
   flag.StringVar(&c.address, "a", "", "address")
   // 4
   flag.StringVar(&c.dash_id, "d", "", "DASH ID")
   set := maya.Parse()
   switch {
   case set["w"]:
      return cache.Write(c)
   case set["e"] && set["p"]:
      return c.do_email_password()
   }
   if err != nil {
      return err
   }
   switch {
   case set["a"]:
      return c.do_address()
   case set["d"]:
      return c.Job.DownloadDash(
         c.Dash.Body, c.Dash.Url, c.dash_id, c.Playout.Widevine,
      )
   }
   return maya.Usage([][]string{
      {"w"},
      {"e", "p"},
      {"a"},
      {"d"},
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
